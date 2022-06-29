package main

import (
	"github.com/jerryan999/goapp/internal/api"
	"github.com/jerryan999/goapp/internal/configs"
	"github.com/jerryan999/goapp/internal/pkg/cachestore"
	"github.com/jerryan999/goapp/internal/pkg/datastore"
	"github.com/jerryan999/goapp/internal/pkg/logger"
	"github.com/jerryan999/goapp/internal/server/http"
	"github.com/jerryan999/goapp/internal/users"
)

func main() {
	l := logger.New("goapp", "v1.0.0", 1)

	jsonFilePath := "./config.json"
	cfg, err := configs.NewService(jsonFilePath)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	dscfg, err := cfg.Datastore()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	mongoClient, err := datastore.NewService(dscfg)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	cacheCfg, err := cfg.Cachestore()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	redispool, err := cachestore.NewService(cacheCfg)
	if err != nil {
		// Cache could be something we'd be willing to tolerate if not available
		// Though this is strictly based on how critical cache is to your application
		l.Error(err)
		return
	}

	us, err := users.NewService(l, mongoClient, redispool)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	a, err := api.NewService(l, us)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	h, err := http.NewService(
		httpCfg,
		a,
	)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	h.Start()
}
