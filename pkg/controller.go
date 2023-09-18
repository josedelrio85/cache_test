package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/josedelrio85/test_cache/pkg/service"
)

type controller struct {
	bussinesShout   service.Actionable
	businessWhistle service.Actionable
	cacheClient     *memcache.Client
}

func NewController(bussinesShout, bussinesWhistle service.Actionable, cache *memcache.Client) controller {
	return controller{
		bussinesShout:   bussinesShout,
		businessWhistle: bussinesWhistle,
		cacheClient:     cache,
	}
}

func (c controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	value, err := c.cacheClient.Get("counter")
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorStr := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		json.NewEncoder(w).Encode(errorStr)
		return
	}

	var counter string
	if value != nil {
		counter = string(value.Value)
	}

	response := struct {
		Shout   string `json:"shout"`
		Whistle string `json:"whistle"`
		Counter string `json:"counter"`
	}{
		Shout:   c.bussinesShout.DoSomething("THE CLASH"),
		Whistle: c.businessWhistle.DoSomething("THE SPECIALS"),
		Counter: counter,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
