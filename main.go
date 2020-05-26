package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cfg = readConfig()

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// spawn hornet collectors
	for _, hornetCfg := range cfg.HornetNodes {
		log.Printf("spawning hornet collector for %s - %s", hornetCfg.ID, hornetCfg.Host)
		go spawnHornetCollector(hornetCfg.ID, hornetCfg.Host)
	}

	// spawn folder collectors
	for _, folderCfg := range cfg.FolderSizes {
		log.Printf("spawning folder collector for %s - %s", folderCfg.ID, folderCfg.Path)
		go spawnFolderSizeCollector(folderCfg.ID, folderCfg.Path)
	}

	// spawn IRI collectors
	for _, iriCfg := range cfg.IRINodes {
		log.Printf("spawning IRI collector for %s - %s", iriCfg.ID, iriCfg.Host)
		go spawnIRICollector(iriCfg.ID, iriCfg.Host)
	}

	registerPrometheusMetricsEndpoint()
}

func registerPrometheusMetricsEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		log.Println(err)
	}
}
