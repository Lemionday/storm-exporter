package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type spoutSummary struct {
	ID        string `json:"spoutId"`
	Executors int    `json:"executors"`
	Tasks     int    `json:"tasks"`

	MessagesEmitted     int    `json:"emitted"`
	MessagesTransferred int    `json:"transferred"`
	MessagesAcked       int    `json:"acked"`
	MessagesFailed      int    `json:"failed"`
	MessagesLatency     string `json:"completeLatency"`
}

type spoutMetrics struct {
	numOfExecutors *prometheus.GaugeVec
	numOfTasks     *prometheus.GaugeVec

	messagesEmitted    *prometheus.GaugeVec
	messagesTransfered *prometheus.GaugeVec
	messagesAcked      *prometheus.GaugeVec
	messagesFailed     *prometheus.GaugeVec
	messagesLatency    *prometheus.GaugeVec
}

func NewSpoutMetrics(reg prometheus.Registerer) *spoutMetrics {
	m := &spoutMetrics{
		numOfExecutors: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_executors",
			Help: "Number of executors of the spout",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
		numOfTasks: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_tasks",
			Help: "Number of spout's tasks",
		}, []string{"TopoName", "TopoID", "SpoutID"}),

		messagesLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_messages_latency",
			Help: "Spout's latency for processing messages",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
		messagesEmitted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_messages_emitted",
			Help: "Number of messages emitted",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
		messagesTransfered: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_messages_transferred",
			Help: "Number of messages transferred",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
		messagesAcked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_messages_acked",
			Help: "Number of messages acked",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
		messagesFailed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_spout_messages_failed",
			Help: "Number of messages failed",
		}, []string{"TopoName", "TopoID", "SpoutID"}),
	}

	reg.MustRegister(m.numOfExecutors)
	reg.MustRegister(m.numOfTasks)

	reg.MustRegister(m.messagesLatency)
	reg.MustRegister(m.messagesEmitted)
	reg.MustRegister(m.messagesTransfered)
	reg.MustRegister(m.messagesAcked)
	reg.MustRegister(m.messagesFailed)

	return m
}

func collectSpoutMetrics(m *spoutMetrics, summary []spoutSummary, topoName, topoID string) {
	for _, spout := range summary {
		id := spout.ID

		m.numOfExecutors.WithLabelValues(topoName, topoID, id).Set(float64(spout.Executors))
		m.numOfTasks.WithLabelValues(topoName, topoID, id).Set(float64(spout.Tasks))

		m.messagesAcked.WithLabelValues(topoName, topoID, id).
			Set(float64(spout.MessagesAcked))
		latency := StringToFloat(spout.MessagesLatency)
		m.messagesLatency.WithLabelValues(topoName, topoID, id).
			Set(latency)
		m.messagesEmitted.WithLabelValues(topoName, topoID, id).
			Set(float64(spout.MessagesEmitted))
		m.messagesTransfered.WithLabelValues(topoName, topoID, id).
			Set(float64(spout.MessagesTransferred))
		m.messagesFailed.WithLabelValues(topoName, topoID, id).
			Set(float64(spout.MessagesFailed))
	}
}
