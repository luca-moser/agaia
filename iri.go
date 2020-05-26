package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/iotaledger/iota.go/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type IRIMetrics struct {
	TotalMemory prometheus.Gauge
	FreeMemory  prometheus.Gauge
	MaxMemory   prometheus.Gauge
}

func NewIRIMetrics(id string) *IRIMetrics {
	return &IRIMetrics{
		TotalMemory: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "total",
			Help:      "Total memory",
		}),
		FreeMemory: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "free",
			Help:      "Free memory",
		}),
		MaxMemory: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "max",
			Help:      "Max memory",
		}),
	}
}

func spawnIRICollector(id string, host string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	log.Printf("connecting to %s", host)
	iotaAPI, err := api.ComposeAPI(api.HTTPClientSettings{
		URI: host,
	})
	if err != nil {
		log.Fatal(err)
	}

	iriMetrics := NewIRIMetrics(id)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				info, err := iotaAPI.GetNodeInfo()
				if err != nil {
					log.Println("unable to fetch node info", err)
					continue
				}
				iriMetrics.TotalMemory.Set(float64(info.JREFreeMemory))
				iriMetrics.MaxMemory.Set(float64(info.JREMaxMemory))
				iriMetrics.FreeMemory.Set(float64(info.JREFreeMemory))
			case <-interrupt:
				return
			}
		}
	}()
	<-interrupt
}
