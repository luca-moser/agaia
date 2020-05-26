package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type FolderSizeMetrics struct {
	Size prometheus.Gauge
}

func NewFolderSizeMetrics(id string) *FolderSizeMetrics {
	return &FolderSizeMetrics{
		Size: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: fmt.Sprintf("iota_benchmark_folder_%s", id),
			Subsystem: "",
			Name:      "size",
			Help:      "Folder Size",
		}),
	}
}

func spawnFolderSizeCollector(id string, dir string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	folderSizeMetrics := NewFolderSizeMetrics(id)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				size, err := DirSize(dir)
				if err != nil {
					continue
				}
				folderSizeMetrics.Size.Set(float64(size))
			case <-interrupt:
				return
			}
		}
	}()
	<-interrupt
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
