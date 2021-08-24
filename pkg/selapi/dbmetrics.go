package selapi

import (
	"encoding/json"
	"fmt"
	"github.com/ktsstudio/selectel-exporter/pkg/apperrors"
	"net/http"
	"strconv"
)

type DatabaseMetric struct {
	DatastoreMetric
	DbName string `json:"db_name"`
}

type DatabaseMetricsResponses struct {
	Metrics struct {
		Connections        []DatabaseMetric `json:"memory_percent"`
		MaxTxDuration      []DatabaseMetric `json:"max_tx_duration"`
		XActCommitRollback []DatabaseMetric `json:"xact_commit_rollback"`
		XActCommit         []DatabaseMetric `json:"xact_commit"`
		XActRollback       []DatabaseMetric `json:"xact_rollback"`
		TupDeleted         []DatabaseMetric `json:"tup_deleted"`
		TupFetched         []DatabaseMetric `json:"tup_fetched"`
		TupInserted        []DatabaseMetric `json:"tup_inserted"`
		TupReturned        []DatabaseMetric `json:"tup_returned"`
		TupUpdated         []DatabaseMetric `json:"tup_updated"`
		CacheHitRatio      []DatabaseMetric `json:"cache_hit_ratio"`
		Deadlocks          []DatabaseMetric `json:"deadlocks"`
		Locks              []DatabaseMetric `json:"locks"`
	} `json:"metrics"`
}

func FetchDatabaseMetrics(token, region, datastoreId string, start, end int64) (*DatabaseMetricsResponses, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.dbaas.selcloud.ru/v1/datastores/%s/database-metrics", region, datastoreId),
		nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", token)
	q := req.URL.Query()
	q.Add("start", strconv.FormatInt(start, 10))
	q.Add("end", strconv.FormatInt(end, 10))
	req.URL.RawQuery = q.Encode()

	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &DatabaseMetricsResponses{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, apperrors.NewResponseFormatError("DatabaseMetricsResponses")
	}
	return resp, nil
}
