package selapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Token struct {
	Id string `json:"id"`
}

type TokensResponse struct {
	Token Token `json:"token"`
}

type tokenRequest struct {
	ProjectId string `json:"project_id"`
}

type tokensRequest struct {
	Token tokenRequest `json:"token"`
}

func ObtainToken(token, projectId string) (*TokensResponse, error) {
	body := tokensRequest{Token: tokenRequest{ProjectId: projectId}}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://my.selectel.ru/api/vpc/resell/v2/tokens",
		bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Token", token)
	req.Header.Add("Content-Type", "application/json")
	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &TokensResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
