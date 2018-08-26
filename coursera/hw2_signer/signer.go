package main

/*
 * Strategy for speedup is very simple:
 *   - wrap MD5-calls with mutex
 *   - with CRC32 - do all calls in-parallel
 */

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Execute pipeline of jobs joined by channels (like pipes)
func ExecutePipeline(jobs ...job) {
	var prevChan chan interface{}
	wg := &sync.WaitGroup{}

	for _, jb := range jobs {
		ch := make(chan interface{}, MaxInputDataLen)

		wg.Add(1)
		go func(j job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer close(out)
			defer wg.Done()
			j(in, out)
		}(jb, prevChan, ch, wg)

		prevChan = ch
	}

	wg.Wait()
}

var md5CallMutex sync.Mutex

func CallMd5(s string) string {
	md5CallMutex.Lock()
	defer md5CallMutex.Unlock()
	return DataSignerMd5(s)
}

func SingleHashImpl(x int, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	val := strconv.Itoa(x)

	cc := []chan string{make(chan string), make(chan string)}
	go func() {
		cc[0] <- DataSignerCrc32(val)
	}()

	go func() {
		cc[1] <- DataSignerCrc32(CallMd5(val))
	}()

	res := <-cc[0] + "~" + <-cc[1]
	out <- res
}

// crc32(data) + "~" + crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for x := range in {
		wg.Add(1)
		go SingleHashImpl(x.(int), out, wg)
	}

	wg.Wait()
}

func MultiHashImpl(s string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	const crcAmount = 6
	var cc [crcAmount]chan string

	for i := 0; i < crcAmount; i++ {
		cc[i] = make(chan string)
		go func(th int) {
			cc[th] <- DataSignerCrc32(strconv.Itoa(th) + s)
		}(i)
	}

	var parts []string
	for _, c := range cc {
		parts = append(parts, <-c)
	}

	out <- strings.Join(parts, "")
}

// concat: crc32(str(th) + data), where th is 0..5
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for x := range in {
		wg.Add(1)
		go MultiHashImpl(x.(string), out, wg)
	}

	wg.Wait()
}

// sort and concat input strings via '_'
func CombineResults(in, out chan interface{}) {
	var arr []string

	for v := range in {
		arr = append(arr, v.(string))
	}
	sort.Strings(arr)
	out <- strings.Join(arr, "_")
}
