package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type topologySummary struct {
	Name           string `json:"name"`
	ID             string `json:"id"`
	UpTime         int    `json:"uptimeSeconds"`
	TasksTotal     int    `json:"tasksTotal"`
	WorkersTotal   int    `json:"workersTotal"`
	ExecutorsTotal int    `json:"executorsTotal"`

	RequestedMemoryOnHeep  float64 `json:"requestedMemOnHeap"`
	RequestedMemoryOffHeep float64 `json:"requestedMemoryOffHeep"`
	RequestedTotalMemory   float64 `json:"requestedTotalMemory"`
	RequestedCPU           float64 `json:"requestedCpu"`

	AssignedMemoryOnHeep  float64 `json:"assignedMemOnHeap"`
	AssignedMemoryOffHeep float64 `json:"assignedMemoryOffHeep"`
	AssignedTotalMemory   float64 `json:"assignedTotalMemory"`
	AssignedCPU           float64 `json:"assignedCpu"`

	StatsTransfered int     `json:"statsTransfered"`
	StatsEmitted    int     `json:"statsEmitted"`
	StatsAcked      int     `json:"statsAcked"`
	StatsFailed     int     `json:"statsFailed"`
	StatsLatency    float64 `json:"statsLatency"`
}

type topologiesSummary struct {
	Topologies []topologySummary `json:"topologies,omitempty"`
}

func main() {
	// flag.Parse()
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	stormUIHost := os.Getenv("STORM_UI_HOST")
	if stormUIHost == "" {
		stormUIHost = "localhost:8081"
	}

	refreshRateStr := os.Getenv("REFRESH_RATE")
	var refreshRate int64
	if refreshRateStr == "" {
		refreshRate = 5
	} else {
		var err error
		refreshRate, err = strconv.ParseInt(refreshRateStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a non-global registry
	reg := prometheus.NewRegistry()

	clusterMetric := NewClusterMetrics(reg)
	go func() {
		for {
			collectClusterMetrics(clusterMetric, stormUIHost)
			time.Sleep(time.Duration(refreshRate))
		}
	}()

	topologyMetrics := NewTopologyMetrics(reg)
	spoutMetrics := NewSpoutMetrics(reg)
	boltMetrics := NewBoltMetrics(reg)
	go func() {
		for {
			func() {
				topologies, err := FetchAndDecode[topologiesSummary](
					fmt.Sprintf("http://%s/api/v1/topology/summary", stormUIHost),
				)
				if err != nil {
					log.Println(err)
					return
				}

				for _, topo := range topologies.Topologies {
					collectTopologyMetrics(topologyMetrics, topo)

					data, err := FetchAndDecode[struct {
						Spouts []spoutSummary `json:"spouts"`
						Bolts  []boltSummary  `json:"bolts"`
					}](
						fmt.Sprintf(
							"http://%s/api/v1/topology/%s?windowSize=5",
							stormUIHost,
							topo.ID,
						),
					)
					if err != nil {
						log.Println(err)
						return
					}
					collectSpoutMetrics(spoutMetrics, data.Spouts, topo.Name, topo.ID)
					collectBoltMetrics(boltMetrics, data.Bolts, topo.Name, topo.ID)
				}
				defer func() {
					time.Sleep(time.Duration(refreshRate))
				}()
			}()
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(addr, nil))
}
