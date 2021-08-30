package main

import (
	"context"
	"github.com/ktsstudio/selectel-exporter/pkg/config"
	"github.com/ktsstudio/selectel-exporter/pkg/exporter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
	exporterConfig, err := config.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	exp, err := exporter.Init(exporterConfig, 60)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	handler := promhttp.HandlerFor(exporter.Registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	srv := &http.Server{Addr: "0.0.0.0:9100", Handler: nil}
	go func() {
		log.Println(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	exp.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
