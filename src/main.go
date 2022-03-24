package main

import (
	"log"

	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/global"
	"github.com/sudoblockio/icon-go-api/healthcheck"
	"github.com/sudoblockio/icon-go-api/logging"
	"github.com/sudoblockio/icon-go-api/metrics"
	_ "github.com/sudoblockio/icon-go-api/models" // for swagger docs
	"github.com/sudoblockio/icon-go-api/redis"
)

func main() {
	config.ReadEnvironment()

	logging.Init()
	log.Printf("Main: Starting logging with level %s", config.Config.LogLevel)

	// Start Prometheus client
	// Go routine starts in function
	metrics.Start()

	// Start Redis Client
	// NOTE: redis is used for websockets
	redis.GetBroadcaster().Start()
	redis.GetRedisClient().StartSubscriber()

	// Start API server
	// Go routine starts in function
	routes.Start()

	// Start Health server
	// Go routine starts in function
	healthcheck.Start()

	global.WaitShutdownSig()
}
