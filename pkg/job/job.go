package job

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/josedelrio85/test_cache/pkg"
)

type Job interface {
	Run()
	SleepTime() time.Duration
}

func RegisterJob(j Job) {
	go func(j Job) {
		for {
			j.Run()
			time.Sleep(j.SleepTime())
		}
	}(j)
}

type myJob struct {
	cacheClient *memcache.Client
	safeCounter pkg.SafeCounter
}

func NewJob(cacheClient *memcache.Client, counter pkg.SafeCounter) *myJob {
	return &myJob{
		cacheClient: cacheClient,
		safeCounter: counter,
	}
}

func (m myJob) Run() {
	m.safeCounter.Inc("counter")
	counterStr := fmt.Sprintf("%d", m.safeCounter.Value("counter"))
	fmt.Println("Running...", counterStr)
	if err := m.cacheClient.Set(&memcache.Item{Key: "counter", Value: []byte(counterStr)}); err != nil {
		fmt.Println(err)
	}
}

func (myJob) SleepTime() time.Duration {
	return time.Second * 5
}
