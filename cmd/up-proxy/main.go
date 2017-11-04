package main

import (
	"os"
	"time"

	"github.com/apex/go-apex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"

	"github.com/apex/up/handler"
	"github.com/apex/up/internal/proxy"
)

func main() {
	if s := os.Getenv("LOG_LEVEL"); s != "" {
		log.SetLevelFromString(s)
	}

	log.SetHandler(json.Default)

<<<<<<< HEAD
=======
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
	log.Infof("initialized in %s", time.Since(start))

	// init handler
>>>>>>> add initialized time
	h, err := handler.New()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	apex.Handle(proxy.NewHandler(h))
}
