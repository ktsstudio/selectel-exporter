package exporter

import (
	"fmt"
	"github.com/ktsstudio/selectel-exporter/pkg/selapi"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func buildGaugeKey(metricName, ip, dbName string) string {
	return fmt.Sprintf("%s_%s_%s", metricName, ip, dbName)
}

type databaseCollector struct {
	project selapi.Project
	datastore selapi.Datastore
	metrics   map[string]prometheus.Gauge
}

func NewDatabaseCollector(project selapi.Project, datastore selapi.Datastore) *databaseCollector {
	return &databaseCollector{
		project:   project,
		datastore: datastore,
		metrics:   make(map[string]prometheus.Gauge),
	}
}

func (col *databaseCollector) registerGauge(metricName string, metric selapi.DatabaseMetric)  {
	key := buildGaugeKey(metricName, metric.Ip, metric.DbName)
	g, ok := col.metrics[key]
	if !ok {
		instance := col.datastore.GetInstance(metric.Ip)
		g = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metricName,
			ConstLabels: prometheus.Labels{
				"project": col.project.Name,
				"datastore": col.datastore.Name,
				"ip": metric.Ip,
				"database": metric.DbName,
				"role": instance.Role,
			},
		})
		Registry.MustRegister(g)
		col.metrics[key] = g
	}
	g.Set(metric.Last)
}

func (col *databaseCollector) loadMetrics(metricName string, metrics []selapi.DatabaseMetric) {
	for _, m := range metrics {
		col.registerGauge(metricName, m)
	}
}

func (col *databaseCollector) GetInfo() string {
	return fmt.Sprintf(
		"project: %s, datastore: %s - collect database metrics", col.project.Name, col.datastore.Name)
}

func (col *databaseCollector) Collect(e *exporter) error {
	start := time.Now().Add(-1 * time.Minute).Unix()
	end := time.Now().Unix()
	res, err := selapi.FetchDatabaseMetrics(e.openstackAccountToken, e.region, col.datastore.Id, start, end)
	if err != nil {
		return err
	}

	col.loadMetrics("selectel_database_locks", res.Metrics.Locks)
	col.loadMetrics("selectel_database_deadlocks", res.Metrics.Deadlocks)
	col.loadMetrics("selectel_database_cache_hit_ratio", res.Metrics.CacheHitRatio)
	col.loadMetrics("selectel_database_tup_updated", res.Metrics.TupUpdated)
	col.loadMetrics("selectel_database_tup_returned", res.Metrics.TupReturned)
	col.loadMetrics("selectel_database_tup_inserted", res.Metrics.TupInserted)
	col.loadMetrics("selectel_database_tup_fetched", res.Metrics.TupFetched)
	col.loadMetrics("selectel_database_tup_deleted", res.Metrics.TupDeleted)
	col.loadMetrics("selectel_database_xact_rollback", res.Metrics.XActRollback)
	col.loadMetrics("selectel_database_xact_commit", res.Metrics.XActCommit)
	col.loadMetrics("selectel_database_xact_commit_rollback", res.Metrics.XActCommitRollback)
	col.loadMetrics("selectel_database_max_tx_duration", res.Metrics.MaxTxDuration)
	col.loadMetrics("selectel_database_connections", res.Metrics.Connections)

	col.loadMetrics("selectel_database_total_connections", res.Metrics.TotalConnections)
	col.loadMetrics("selectel_database_commands_total_delete", res.Metrics.CommandsTotalDelete)
	col.loadMetrics("selectel_database_commands_total_insert", res.Metrics.CommandsTotalInsert)
	col.loadMetrics("selectel_database_commands_total_select", res.Metrics.CommandsTotalSelect)
	col.loadMetrics("selectel_database_commands_total_update", res.Metrics.CommandsTotalUpdate)
	col.loadMetrics("selectel_database_innodb_buffer_pool_hit_ratio", res.Metrics.InnodbBufferPoolHitRatio)
	col.loadMetrics("selectel_database_slow_queries", res.Metrics.SlowQueries)
	col.loadMetrics("selectel_database_threads_cached", res.Metrics.ThreadsCached)
	col.loadMetrics("selectel_database_threads_connected", res.Metrics.ThreadsConnected)
	col.loadMetrics("selectel_database_threads_running", res.Metrics.ThreadsRunning)

	return nil
}
