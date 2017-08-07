package main

import (
	"os"

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

	h, err := handler.New()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	apex.Handle(proxy.NewHandler(h))
}
