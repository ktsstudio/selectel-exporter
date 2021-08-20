package selapi

import (
	"encoding/json"
	"fmt"
	"github.com/ktsstudio/selectel-exporter/pkg/apperrors"
	"net/http"
	"strconv"
)

type DatastoreMetric struct {
	DatastoreId string    `json:"datastore_id"`
	Ip          string    `json:"ip"`
	Timestamps  []float64 `json:"timestamps"`
	Values      []float64 `json:"values"`
	Max         float64   `json:"max"`
	Min         float64   `json:"min"`
	Avg         float64   `json:"avg"`
	Last        float64   `json:"last"`
}

type DatastoreMetricsResponses struct {
	Metrics struct {
		Step          float64           `json:"step"`
		MemoryPercent []DatastoreMetric `json:"memory_percent"`
		MemoryBytes   []DatastoreMetric `json:"memory_bytes"`
		Cpu           []DatastoreMetric `json:"cpu"`
		DiskPercent   []DatastoreMetric `json:"disk_percent"`
		DiskBytes     []DatastoreMetric `json:"disk_bytes"`
	} `json:"metrics"`
}

func FetchDatastoreMetrics(token, region, datastoreId string, start, end int64) (*DatastoreMetricsResponses, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.dbaas.selcloud.ru/v1/datastores/%s/metrics", region, datastoreId),
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

	resp := &DatastoreMetricsResponses{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, apperrors.NewResponseFormatError("DatastoreMetricsResponses")
	}
	return resp, nil
}
