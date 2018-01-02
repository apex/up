package main

import (
	"os"
	"time"

	"github.com/apex/go-apex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"

	"github.com/apex/up"
	"github.com/apex/up/handler"
	"github.com/apex/up/internal/proxy"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/lambda/runtime"
)

func main() {
	stage := os.Getenv("UP_STAGE")

	if s := os.Getenv("LOG_LEVEL"); s != "" {
		log.SetLevelFromString(s)
	}

	log.SetHandler(json.Default)
	log.Info("initializing")

	// read config
	c, err := up.ReadConfig("up.json")
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	// init project
	p := runtime.New(c)

	// init runtime
	start := time.Now()
	if err := p.Init(stage); err != nil {
		log.Fatalf("error initializing: %s", err)
	}
	log.WithField("duration", util.MillisecondsSince(start)).Info("initialized")

	// init handler
	h, err := handler.New(c)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// serve
	apex.Handle(proxy.NewHandler(h))
}
