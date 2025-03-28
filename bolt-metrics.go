package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type boltSummary struct {
	ID        string `json:"boltId"`
	Executors int    `json:"executors"`
	Tasks     int    `json:"tasks"`

	Capacity string `json:"capacity"`

	MessagesEmitted     int    `json:"emitted"`
	MessagesTransferred int    `json:"transferred"`
	MessagesAcked       int    `json:"acked"`
	MessagesFailed      int    `json:"failed"`
	ProcessLatency      string `json:"processLatency"`
	ExecuteLatency      string `json:"executeLatency"`
}

type boltMetrics struct {
	numOfExecutors *prometheus.GaugeVec
	numOfTasks     *prometheus.GaugeVec

	capacity *prometheus.GaugeVec

	messagesEmitted     *prometheus.GaugeVec
	messagesTransferred *prometheus.GaugeVec
	messagesAcked       *prometheus.GaugeVec
	messagesFailed      *prometheus.GaugeVec
	processLatency      *prometheus.GaugeVec
	executeLatency      *prometheus.GaugeVec
}

func NewBoltMetrics(reg prometheus.Registerer) *boltMetrics {
	m := &boltMetrics{
		numOfExecutors: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_executors",
			Help: "Number of executors of the bolt",
		}, []string{"TopoName", "TopoID", "boltID"}),
		numOfTasks: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_tasks",
			Help: "Number of bolt's tasks",
		}, []string{"TopoName", "TopoID", "boltID"}),

		capacity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_capacity",
			Help: "Number of messages executed * average execute latency / time window",
		}, []string{"TopoName", "TopoID", "boltID"}),

		processLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_messages_latency",
			Help: "Average latency for processing messages",
		}, []string{"TopoName", "TopoID", "boltID"}),
		executeLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_execute_latency",
			Help: "Average latency for executing bolt's methods",
		}, []string{"TopoName", "TopoID", "boltID"}),
		messagesEmitted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_messages_emitted",
			Help: "Number of messages emitted",
		}, []string{"TopoName", "TopoID", "boltID"}),
		messagesTransferred: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_messages_transferred",
			Help: "Number of messages transferred",
		}, []string{"TopoName", "TopoID", "boltID"}),
		messagesAcked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_messages_acked",
			Help: "Number of messages acked",
		}, []string{"TopoName", "TopoID", "boltID"}),
		messagesFailed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_bolt_messages_failed",
			Help: "Number of messages failed",
		}, []string{"TopoName", "TopoID", "boltID"}),
	}

	reg.MustRegister(m.numOfExecutors)
	reg.MustRegister(m.numOfTasks)

	reg.MustRegister(m.executeLatency)
	reg.MustRegister(m.processLatency)
	reg.MustRegister(m.messagesEmitted)
	reg.MustRegister(m.messagesTransferred)
	reg.MustRegister(m.messagesAcked)
	reg.MustRegister(m.messagesFailed)

	return m
}

func collectBoltMetrics(m *boltMetrics, summary []boltSummary, topoName, topoID string) {
	for _, bolt := range summary {
		id := bolt.ID

		m.numOfExecutors.WithLabelValues(topoName, topoID, id).Set(float64(bolt.Executors))
		m.numOfTasks.WithLabelValues(topoName, topoID, id).Set(float64(bolt.Tasks))

		m.messagesAcked.WithLabelValues(topoName, topoID, id).
			Set(float64(bolt.MessagesAcked))
		m.executeLatency.WithLabelValues(topoName, topoID, id).
			Set(StringToFloat(bolt.ExecuteLatency))
		m.processLatency.WithLabelValues(topoName, topoID, id).
			Set(StringToFloat(bolt.ProcessLatency))
		m.messagesEmitted.WithLabelValues(topoName, topoID, id).
			Set(float64(bolt.MessagesEmitted))
		m.messagesTransferred.WithLabelValues(topoName, topoID, id).
			Set(float64(bolt.MessagesTransferred))
		m.messagesFailed.WithLabelValues(topoName, topoID, id).
			Set(float64(bolt.MessagesFailed))
	}
}
