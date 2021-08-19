package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"kts/selectel-exporter/pkg/selapi"
	"log"
	"time"
)

type datastoreMetrics struct {
	memoryPercent prometheus.Gauge
	memoryBytes   prometheus.Gauge
	cpu           prometheus.Gauge
	diskPercent   prometheus.Gauge
	diskBytes     prometheus.Gauge
}

type datastoreCollector struct {
	project selapi.Project
	datastore selapi.Datastore
	metrics   map[string]*datastoreMetrics
}

func NewDatastoreCollector(project selapi.Project, datastore selapi.Datastore) *datastoreCollector {
	return &datastoreCollector{
		project:   project,
		datastore: datastore,
		metrics:   make(map[string]*datastoreMetrics),
	}
}

func (col *datastoreCollector) registerGauge(name string, g prometheus.Gauge, ip string, value float64) prometheus.Gauge {
	if g == nil {
		g = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        name,
			ConstLabels: prometheus.Labels{"project": col.project.Name, "datastore": col.datastore.Name, "ip": ip},
		})
		prometheus.MustRegister(g)
	}
	g.Set(value)
	return g
}

func (col *datastoreCollector) initMetric(ip string) *datastoreMetrics {
	if _, ok := col.metrics[ip]; !ok {
		col.metrics[ip] = &datastoreMetrics{}
	}
	return col.metrics[ip]
}

func (col *datastoreCollector) loadMemoryPercent(data *selapi.DatastoreMetricsResponses) {
	for _, item := range data.Metrics.MemoryPercent {
		m := col.initMetric(item.Ip)
		m.memoryPercent = col.registerGauge(
			"selectel_datastore_memory_percent", m.memoryPercent, item.Ip, item.Last)
	}
}

func (col *datastoreCollector) loadMemoryBytes(data *selapi.DatastoreMetricsResponses) {
	for _, item := range data.Metrics.MemoryBytes {
		m := col.initMetric(item.Ip)
		m.memoryBytes = col.registerGauge(
			"selectel_datastore_memory_bytes", m.memoryBytes, item.Ip, item.Last)
	}
}

func (col *datastoreCollector) loadCpu(data *selapi.DatastoreMetricsResponses) {
	for _, item := range data.Metrics.Cpu {
		m := col.initMetric(item.Ip)
		m.cpu = col.registerGauge("selectel_datastore_cpu", m.cpu, item.Ip, item.Last)
	}
}

func (col *datastoreCollector) loadDiskPercent(data *selapi.DatastoreMetricsResponses) {
	for _, item := range data.Metrics.DiskPercent {
		m := col.initMetric(item.Ip)
		m.diskPercent = col.registerGauge(
			"selectel_datastore_disk_percent", m.diskPercent, item.Ip, item.Last)
	}
}

func (col *datastoreCollector) loadDiskBytes(data *selapi.DatastoreMetricsResponses) {
	for _, item := range data.Metrics.DiskBytes {
		m := col.initMetric(item.Ip)
		m.diskBytes = col.registerGauge(
			"selectel_datastore_disk_bytes", m.diskBytes, item.Ip, item.Last)
	}
}

func (col *datastoreCollector) Collect(e *exporter) error {
	log.Println("collect datastore metrics")
	start := time.Now().Add(-1 * time.Minute).Unix()
	end := time.Now().Unix()
	res, err := selapi.FetchDatastoreMetrics(e.openstackAccountToken, e.region, col.datastore.Id, start, end)
	if err != nil {
		return err
	}
	col.loadMemoryPercent(res)
	col.loadMemoryBytes(res)
	col.loadCpu(res)
	col.loadDiskPercent(res)
	col.loadDiskBytes(res)
	return nil
}
