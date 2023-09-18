package main

import (
	"log"
	"net/http"

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

//https://michaelheap.com/golang-using-memcached/

func main() {
	cache := initMemcache()
	defer cache.Close()
	if err := cache.Ping(); err != nil {
		log.Fatalf("can't ping memcached %v", err)
	}

	counter := pkg.NewSafeCounter()

	myJob := job.NewJob(cache, counter)
	job.RegisterJob(myJob)

	shoutService, whistleService := setupDependencies()
	controller := pkg.NewController(shoutService, whistleService, cache)

	mux := http.NewServeMux()
	mux.Handle("/", controller)
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("can't listen %v", err)
	}
}

func initMemcache() *memcache.Client {
	return memcache.New("0.0.0.0:11211")
}
