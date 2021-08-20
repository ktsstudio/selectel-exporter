package selapi

import (
	"github.com/ktsstudio/selectel-exporter/pkg/apperrors"
	"io/ioutil"
	"net/http"
)

func fetch(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, apperrors.NewRequestError(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apperrors.NewRequestError(err.Error())
	}
	if resp.StatusCode >= 300 {
		if resp.StatusCode == 401 {
			return nil, &apperrors.SelectelApiError{
				Code: resp.StatusCode,
				Body: "authorization required, check token",
			}
		}
		return nil, &apperrors.SelectelApiError{Code: resp.StatusCode, Body: string(body)}
	}
	return body, nil
}
