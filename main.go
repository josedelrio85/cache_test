package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/josedelrio85/test_cache/pkg"
	"github.com/josedelrio85/test_cache/pkg/job"
	"github.com/josedelrio85/test_cache/pkg/service"
)

func setupDependencies() (service.Actionable, service.Actionable) {
	shoutService := service.NewActionShout("FIGHT FOR YOUR RIGHTS!")
	whistleService := service.NewActionWhistle("e_du_ardo...", 1)
	return shoutService, whistleService
}

type MyJob struct {
	cacheClient *memcache.Client
}

var counter = 0

func (m MyJob) Run() {
	counter++
	counterStr := fmt.Sprintf("%d", counter)
	fmt.Println("Running...", counterStr)
	if err := m.cacheClient.Set(&memcache.Item{Key: "counter", Value: []byte(counterStr)}); err != nil {
		fmt.Println(err)
	}
}

func (MyJob) SleepTime() time.Duration {
	return time.Second * 5
}

//https://michaelheap.com/golang-using-memcached/

func main() {
	cache := initMemcache()
	defer cache.Close()
	if err := cache.Ping(); err != nil {
		log.Fatalf("can't ping memcached %v", err)
	}
	myJob := MyJob{
		cacheClient: cache,
	}

	job.RegisterJob(&myJob)

	shoutService, whistleService := setupDependencies()
	controller := pkg.NewController(shoutService, whistleService, cache)

	mux := http.NewServeMux()
	mux.Handle("/", controller)
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("can't listen %v", err)
	}
}

func initMemcache() *memcache.Client {
	return memcache.New("127.0.0.1:11211")
}
