package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func choice(arr...string) string {
	return arr[rand.Int() % len(arr)]
}

func generateTaskInfo() (inf TaskInfo) {
	inf.Application = choice(
		"aurora-update-garbage-collector",
		"cag-kafka-ingester",
		"clusto-sync",
		"endpoint-exerciser",
		"instance-locator",
		"statistics-collector",
		"udeploy-aggregator",
	)

	inf.Zone = choice("dca1", "dca4", "dca8", "phx2", "phx3", "sjc1", "sjc4")
	inf.Host = fmt.Sprintf("%s%d-%s", choice("agent", "compute"), rand.Int() % 10000, inf.Zone)
	inf.Revision = choice("21.0.3", "21.1.1")
	inf.Instance = rand.Int() % 1000

	return
}

func BenchmarkRegexpSearch(t *testing.B) {
	rep := NewServiceRepository()

	for i := 0; i < 60000; i++ {
		rep.AddTask("udeploy-aggregator", generateTaskInfo())
	}
	for i := 0; i < 200000; i++ {
		svc := ServiceName(fmt.Sprintf("test-service-%d", rand.Int() % 50))
		rep.AddTask(svc, generateTaskInfo())
	}

		t.Run("full match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := rep.MatchAny("udeploy-aggregator", "dca")
			assert.NotEmpty(t, res)
		}
	})
}