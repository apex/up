package main

import (
	stdjson "encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	apex "github.com/apex/go-apex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"

	"github.com/apex/up"
	"github.com/apex/up/handler"
	"github.com/apex/up/internal/logs"
	"github.com/apex/up/internal/proxy"
	"github.com/apex/up/internal/util"
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

	// // init project
	// p := runtime.New(c)
	//
	// // init runtime
	// if err := p.Init(stage); err != nil {
	// 	ctx.Fatalf("error initializing: %s", err)
	// }

	// read environment variables
	if err := loadEnvironment(ctx); err != nil {
		ctx.Fatalf("error loading environment variables: %s", err)
	}

	// overrides
	if err := c.Override(stage); err != nil {
		ctx.Fatalf("error overriding: %s", err)
	}

	// create handler
	h, err := handler.FromConfig(c)
	if err != nil {
		ctx.Fatalf("error creating handler: %s", err)
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

// loadEnvironment loads environment variables.
func loadEnvironment(ctx log.Interface) error {
	var m map[string]string

	ctx.Info("loading environment variables")
	b, err := ioutil.ReadFile("up-env.json")
	if err != nil {
		return err
	}

	err = stdjson.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for k, v := range m {
		ctx.WithFields(log.Fields{
			"name":  k,
			"value": strings.Repeat("*", len(v)),
		}).Info("set variable")
		os.Setenv(k, v)
	}

	return nil
}
