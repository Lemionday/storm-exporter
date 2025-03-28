package main

import "github.com/prometheus/client_golang/prometheus"

type topologyMetrics struct {
	upTime           *prometheus.GaugeVec
	numOfTasks       *prometheus.GaugeVec
	numOfWorkers     *prometheus.GaugeVec
	numOfExecutors   *prometheus.GaugeVec
	numOfReplication *prometheus.GaugeVec

	requestedMemoryOnHeep  *prometheus.GaugeVec
	requestedMemoryOffHeep *prometheus.GaugeVec
	requestedTotalMemory   *prometheus.GaugeVec
	requestedCPU           *prometheus.GaugeVec

	assignedMemoryOnHeep  *prometheus.GaugeVec
	assignedMemoryOffHeep *prometheus.GaugeVec
	assignedTotalMemory   *prometheus.GaugeVec
	assignedCPU           *prometheus.GaugeVec

	statsTransfered *prometheus.GaugeVec
	statsEmitted    *prometheus.GaugeVec
	statsAcked      *prometheus.GaugeVec
	statsFailed     *prometheus.GaugeVec
	statsLatency    *prometheus.GaugeVec
}

func NewTopologyMetrics(reg prometheus.Registerer) *topologyMetrics {
	m := &topologyMetrics{
		upTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_uptime",
			Help: "How long the topology has been running (seconds)",
		}, []string{"TopoName", "TopoID"}),

		numOfTasks: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_number_of_tasks",
			Help: "Number of tasks in the topology",
		}, []string{"TopoName", "TopoID"}),
		numOfWorkers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_total_workers",
			Help: "Number of workers currently used by the topology",
		}, []string{"TopoName", "TopoID"}),
		numOfExecutors: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_total_executors",
			Help: "Number of executors in the topology",
		}, []string{"TopoName", "TopoID"}),
		numOfReplication: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_replication_count",
			Help: "Number of nimbus hosts on which this topology is replicated",
		}, []string{"TopoName", "TopoID"}),

		requestedMemoryOnHeep: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_requested_mem_on_heap",
			Help: "Requested on-heap memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		requestedMemoryOffHeep: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_requested_mem_off_heap",
			Help: "Requested off-heap memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		requestedTotalMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_requested_total_mem",
			Help: "Requested total memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		requestedCPU: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_requested_cpu",
			Help: "Requested cpu by user (%)",
		}, []string{"TopoName", "TopoID"}),

		assignedMemoryOnHeep: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_assigned_mem_on_heap",
			Help: "Assigned on-heap memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		assignedMemoryOffHeep: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_assigned_mem_off_heap",
			Help: "Assigned off-heap memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		assignedTotalMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_assigned_total_mem",
			Help: "Assigned total memory by user (MB)",
		}, []string{"TopoName", "TopoID"}),
		assignedCPU: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_topology_assigned_cpu",
			Help: "Assigned cpu by user (%)",
		}, []string{"TopoName", "TopoID"}),

		statsTransfered: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_stats_transfered",
			Help: "Number of messages transfered",
		}, []string{"TopoName", "TopoID"}),
		statsEmitted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_stats_emitted",
			Help: "Number of messages emitted",
		}, []string{"TopoName", "TopoID"}),
		statsLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_stats_total_latency",
			Help: "Total messages processing latency",
		}, []string{"TopoName", "TopoID"}),
		statsAcked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_stats_acked",
			Help: "Number of messages acknowledged",
		}, []string{"TopoName", "TopoID"}),
		statsFailed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "storm_stats_failed",
			Help: "Number of messages failed",
		}, []string{"TopoName", "TopoID"}),
	}

	reg.MustRegister(m.upTime)

	reg.MustRegister(m.numOfTasks)
	reg.MustRegister(m.numOfWorkers)
	reg.MustRegister(m.numOfExecutors)
	reg.MustRegister(m.numOfReplication)

	reg.MustRegister(m.requestedCPU)
	reg.MustRegister(m.requestedMemoryOffHeep)
	reg.MustRegister(m.requestedMemoryOnHeep)
	reg.MustRegister(m.requestedTotalMemory)

	reg.MustRegister(m.assignedCPU)
	reg.MustRegister(m.assignedMemoryOffHeep)
	reg.MustRegister(m.assignedMemoryOnHeep)
	reg.MustRegister(m.assignedTotalMemory)

	reg.MustRegister(m.statsAcked)
	reg.MustRegister(m.statsFailed)
	reg.MustRegister(m.statsEmitted)
	reg.MustRegister(m.statsTransfered)
	reg.MustRegister(m.statsLatency)

	return m
}

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

func collectTopologyMetrics(m *topologyMetrics, topo topologySummary) {
	name := topo.Name
	id := topo.ID

	m.upTime.WithLabelValues(name, id).Set(float64(topo.UpTime))
	m.numOfTasks.WithLabelValues(name, id).Set(float64(topo.TasksTotal))
	m.numOfWorkers.WithLabelValues(name, id).Set(float64(topo.WorkersTotal))
	m.numOfExecutors.WithLabelValues(name, id).Set(float64(topo.ExecutorsTotal))

	m.requestedMemoryOnHeep.WithLabelValues(name, id).Set(topo.RequestedMemoryOnHeep)
	m.requestedMemoryOffHeep.WithLabelValues(name, id).Set(topo.RequestedMemoryOffHeep)
	m.requestedTotalMemory.WithLabelValues(name, id).Set(topo.RequestedTotalMemory)
	m.requestedCPU.WithLabelValues(name, id).Set(topo.RequestedCPU)

	m.assignedMemoryOnHeep.WithLabelValues(name, id).Set(topo.AssignedMemoryOnHeep)
	m.assignedMemoryOffHeep.WithLabelValues(name, id).Set(topo.AssignedMemoryOffHeep)
	m.assignedTotalMemory.WithLabelValues(name, id).Set(topo.AssignedTotalMemory)
	m.assignedCPU.WithLabelValues(name, id).Set(topo.AssignedCPU)

	m.statsAcked.WithLabelValues(name, id).Set(float64(topo.StatsAcked))
	m.statsFailed.WithLabelValues(name, id).Set(float64(topo.StatsFailed))
	m.statsTransfered.WithLabelValues(name, id).Set(float64(topo.StatsTransfered))
	m.statsEmitted.WithLabelValues(name, id).Set(float64(topo.StatsEmitted))
	m.statsLatency.WithLabelValues(name, id).Set(float64(topo.StatsLatency))
}
