package selapi

import (
	"encoding/json"
	"net/http"
)

type Project struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

func FetchProjects(token string) (*ProjectsResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.selectel.ru/vpc/resell/v2/projects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Token", token)

	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &ProjectsResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
