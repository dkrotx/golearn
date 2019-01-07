package main

import (
	"context"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/shirou/gopsutil/load"
	"github.com/uber-go/tally"
	tallystatsd "github.com/uber-go/tally/statsd"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Context struct {
	logger *zap.SugaredLogger
	scope  tally.Scope
	histogram tally.Histogram
}

func randomDuration(fromMs, toMs int) time.Duration {
	val := fromMs + rand.Int()%(toMs-fromMs)
	return time.Millisecond * time.Duration(val)
}

func (ctx *Context) handler(w http.ResponseWriter, r *http.Request) {
	ctx.scope.Counter("req_cnt").Inc(1)
	ctx.logger.Infof("New Request: %s", r.URL.Path)

	duration := randomDuration(20, 100)
	ctx.scope.Timer("handler_time").Record(duration)
	ctx.histogram.RecordDuration(duration)
	fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path[1:])
}

func (ctx *Context) errorsHandler(w http.ResponseWriter, r *http.Request) {
	who := strings.ToLower(r.URL.Query().Get("who"))
	available := map[string]bool {
		"redis": true,
		"mysql": true,
		"file_write": true,
	}

	// we don't want to produce absolutely random error and flood graphs
	if available[who] {
		ctx.logger.Infof("New (synthetic) error: %s", who)
		ctx.scope.SubScope("errors").Counter(who).Inc(1)
	}
}

func (ctx *Context) sysinfoReport() error {
	stat, err := load.Avg()
	if err != nil {
		return err
	}

	// scope send gauges as int64 (why?!), so we should care about precision itself
	ctx.scope.Gauge("la.1min").Update(stat.Load1*100)
	ctx.scope.Gauge("la.5min").Update(stat.Load5*100)
	ctx.scope.Gauge("la.15min").Update(stat.Load15*100)

	ctx.logger.Infof("LA: %.2f %.2f %.2f", stat.Load1, stat.Load5, stat.Load15)
	return nil
}

func (ctx *Context) sysinfoReportLoop(c context.Context) {
	timer := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-c.Done():
			return
		case <-timer.C:
			if err := ctx.sysinfoReport(); err != nil {
				ctx.logger.Errorf("failed to get load avg: %v", err)
			}
		}
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %s", err)
		os.Exit(1)
	}
}

func main() {
	statter, err := statsd.NewBufferedClient("test.dkrot.pro:8125", "daemons",
		100*time.Millisecond, 1440)
	checkErr(err, "NewBufferedClient failed")

	reporter := tallystatsd.NewReporter(statter, tallystatsd.Options{SampleRate: 1.0})

	logger, err := zap.NewDevelopment()
	checkErr(err, "zap-logger failed")

	scope, _ := tally.NewRootScope(tally.ScopeOptions{Reporter: reporter, Prefix: "test-service"},
		time.Second)

	backgroundCtx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ctx := &Context{
		logger: logger.Sugar(),
		scope: scope,
		histogram: scope.Histogram("handler_hist",
			tally.DurationBuckets{time.Millisecond * 30, time.Millisecond * 60, time.Millisecond * 90, time.Millisecond * 200},
		),
	}

	go ctx.sysinfoReportLoop(backgroundCtx)


	http.HandleFunc("/", ctx.handler)
	http.HandleFunc("/error", ctx.errorsHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
