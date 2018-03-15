package main

import (
	"os"
	"time"

	"github.com/apex/go-apex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"

	"github.com/apex/up"
	"github.com/apex/up/handler"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/proxy"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/aws/runtime"
)

func main() {
	start := time.Now()
	stage := os.Getenv("UP_STAGE")

	// setup logging
	log.SetHandler(json.Default)
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		log.SetLevelFromString(s)
	}

	log.Log = log.WithFields(logs.Fields())
	log.Info("initializing")

	// read config
	c, err := up.ReadConfig("up.json")
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	ctx := log.WithFields(log.Fields{
		"name": c.Name,
		"type": c.Type,
	})

	// init project
	p := runtime.New(c)

	// init runtime
	if err := p.Init(stage); err != nil {
		ctx.Fatalf("error initializing: %s", err)
	}

	// overrides
	if err := c.Override(stage); err != nil {
		ctx.Fatalf("error overriding: %s", err)
	}

	// select handler
	h, err := handler.FromConfig(c)
	if err != nil {
		ctx.Fatalf("error selecting handler: %s", err)
	}

	// init handler
	h, err = handler.New(c, h)
	if err != nil {
		ctx.Fatalf("error initializing handler: %s", err)
	}

	// metrics
	// err = p.Metric("initialize", float64(util.MillisecondsSince(start)))
	// if err != nil {
	// 	ctx.WithError(err).Warn("putting metric")
	// }

	// serve
	log.WithField("duration", util.MillisecondsSince(start)).Info("initialized")
	apex.Handle(proxy.NewHandler(h))
}
