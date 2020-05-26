package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HornetMetrics struct {
	MemSysTotal prometheus.Gauge
	MemHeap     prometheus.Gauge
	MemSys      prometheus.Gauge
	MemIdle     prometheus.Gauge
	MemReleased prometheus.Gauge
	MemObjects  prometheus.Gauge
}

func NewHornetMetrics(id string) *HornetMetrics {
	return &HornetMetrics{
		MemSysTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "total",
			Help:      "Sys total memory",
		}),
		MemHeap: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "heap",
			Help:      "Used heap memory",
		}),
		MemSys: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "sys",
			Help:      "Sys memory",
		}),
		MemIdle: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "heap_idle",
			Help:      "Idle heap memory",
		}),
		MemReleased: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "released",
			Help:      "Released memory",
		}),
		MemObjects: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_%s", id),
			Subsystem: "memory",
			Name:      "objects",
			Help:      "Objects count",
		}),
	}
}

func spawnHornetCollector(id string, host string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	hornetMetrics := NewHornetMetrics(id)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			wsMessage := &WebSocketMsg{}
			if err := json.Unmarshal(message, wsMessage); err != nil {
				log.Println("read:", err)
				return
			}

			if wsMessage.Type != MsgTypeNodeStatus {
				continue
			}

			d, err := json.Marshal(wsMessage.Data)
			if err != nil {
				continue
			}

			status := &NodeStatus{}
			if err := json.Unmarshal(d, status); err != nil {
				continue
			}

			hornetMetrics.MemHeap.Set(float64(status.Mem.HeapInuse))
			hornetMetrics.MemIdle.Set(float64(status.Mem.HeapIdle))
			hornetMetrics.MemObjects.Set(float64(status.Mem.HeapObjects))
			hornetMetrics.MemReleased.Set(float64(status.Mem.HeapReleased))
			hornetMetrics.MemSys.Set(float64(status.Mem.Sys))
			hornetMetrics.MemSysTotal.Set(float64(status.Mem.HeapSys))
		}
	}()

	<-interrupt
}
