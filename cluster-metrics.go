package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type clusterMetrics struct {
	numOfWorkers     prometheus.Gauge
	numOfWorkersUsed prometheus.Gauge
	numOfSupervisors prometheus.Gauge
}

func NewClusterMetrics(reg prometheus.Registerer) *clusterMetrics {
	m := &clusterMetrics{
		numOfWorkers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "storm_cluster_total_worker_processes",
			Help: "Number of worker processes in cluster",
		}),
		numOfWorkersUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "storm_cluster_worker_processes_used",
			Help: "Number of worker processes used in cluster",
		}),
		numOfSupervisors: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "storm_cluster_total_supervisor",
			Help: "Number of supervisors in cluster",
		}),
	}

	reg.MustRegister(m.numOfWorkers)
	reg.MustRegister(m.numOfWorkersUsed)
	reg.MustRegister(m.numOfSupervisors)

	return m
}

func collectClusterMetrics(m *clusterMetrics, stormUIHost string) func() {
	type Summary struct {
		SlotsUsed   int `json:"slotsUsed,omitempty"`
		SlotsTotal  int `json:"slotsTotal,omitempty"`
		Supervisors int `json:"supervisors,omitempty"`
	}

	return func() {
		summary, err := FetchAndDecode[Summary](
			fmt.Sprintf("http://%s/api/v1/cluster/summary", stormUIHost),
		)
		if err != nil {
			log.Fatal(err)
		}

		m.numOfWorkers.Set(float64(summary.SlotsTotal))
		m.numOfWorkersUsed.Set(float64(summary.SlotsUsed))
		m.numOfSupervisors.Set(float64(summary.Supervisors))

		log.Println("Updated cluster metrics")
	}
}
