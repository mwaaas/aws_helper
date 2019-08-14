package main

import (
	"github.com/gorilla/mux"
	"github.com/hellofresh/health-go"
	"time"
)

func status(r *mux.Router) {
	health.Register(health.Config{
		Name:      "server",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: func() error {
			// rabbitmq health check implementation goes here
			return nil
		},
	})
	r.Handle("/status", health.Handler())
}
