package selapi

import (
	"encoding/json"
	"fmt"
	"kts/selectel-exporter/pkg/apperrors"
	"net/http"
)

type Datastore struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Enabled bool `json:"enabled"`
}

type DatastoresResponse struct {
	Datastores []Datastore `json:"datastores"`
}

func FetchDatastores(token, region string) (*DatastoresResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.dbaas.selcloud.ru/v1/datastores", region),
		nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", token)

	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &DatastoresResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, apperrors.NewResponseFormatError("DatastoresResponse")
	}
	return resp, nil
}

type Database struct {
	Id string `json:"id"`
	Name string `json:"name"`
	DatastoreId string `json:"datastore_id"`
	Status string `json:"status"`
}

type DatabasesResponse struct {
	Databases []Database
}

func FetchDatabases(token, region string) (*DatabasesResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.dbaas.selcloud.ru/v1/databases", region),
		nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", token)

	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &DatabasesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, apperrors.NewResponseFormatError("DatabasesResponse")
	}
	return resp, nil
}
