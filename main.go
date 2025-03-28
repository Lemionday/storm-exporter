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

type config struct {
	addr        string
	stormUIHost string
	refreshRate int64
}

var conf config

func main() {
	// flag.Parse()
	conf.addr = os.Getenv("EXPORTER_PORT")
	if conf.addr == "" {
		conf.addr = ":8080"
	}

	conf.stormUIHost = os.Getenv("STORM_UI_HOST")
	if conf.stormUIHost == "" {
		conf.stormUIHost = "localhost:8081"
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
			collectClusterMetrics(clusterMetric, conf.stormUIHost)
			log.Println("Updated cluster's metrics")
			time.Sleep(time.Duration(refreshRate))
		}
	}()

	topologyMetrics := NewTopologyMetrics(reg)
	spoutMetrics := NewSpoutMetrics(reg)
	boltMetrics := NewBoltMetrics(reg)
	go func() {
		for {
			func() {
				topologies, err := FetchAndDecode[struct {
					Topologies []topologySummary `json:"topologies,omitempty"`
				}](
					fmt.Sprintf("http://%s/api/v1/topology/summary", conf.stormUIHost),
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
							conf.stormUIHost,
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
					log.Println("Updated topologies's metrics")
					time.Sleep(time.Duration(refreshRate))
				}()
			}()
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(conf.addr, nil))
}
