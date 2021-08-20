package selapi

import (
	"encoding/json"
	"github.com/ktsstudio/selectel-exporter/pkg/apperrors"
	"net/http"
)

type BalanceResponse struct {
	Status string `json:"status"`
	Data   struct {
		Currency  string `json:"currency"`
		IsPostpay bool   `json:"is_postpay"`
		Discount  int    `json:"discount"`
		Primary   struct {
			Main  int `json:"main"`
			Bonus int `json:"bonus"`
			VkRub int `json:"vk_rub"`
			Ref   int `json:"ref"`
			Hold  struct {
				Main  int `json:"main"`
				Bonus int `json:"bonus"`
				VkRub int `json:"vk_rub"`
			} `json:"hold"`
		} `json:"primary"`
		Storage struct {
			Main       int         `json:"main"`
			Bonus      int         `json:"bonus"`
			VkRub      int         `json:"vk_rub"`
			Prediction interface{} `json:"prediction"`
			Debt       int         `json:"debt"`
			Sum        int         `json:"sum"`
		} `json:"storage"`
		Vpc struct {
			Main       int         `json:"main"`
			Bonus      int         `json:"bonus"`
			VkRub      int         `json:"vk_rub"`
			Prediction interface{} `json:"prediction"`
			Debt       int         `json:"debt"`
			Sum        int         `json:"sum"`
		} `json:"vpc"`
		Vmware struct {
			Main       int         `json:"main"`
			Bonus      int         `json:"bonus"`
			VkRub      int         `json:"vk_rub"`
			Prediction interface{} `json:"prediction"`
			Debt       int         `json:"debt"`
			Sum        int         `json:"sum"`
		} `json:"vmware"`
	} `json:"data"`
}

func FetchBalance(token string) (*BalanceResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://my.selectel.ru/api/v3/billing/balance", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Token", token)

	data, err := fetch(client, req)
	if err != nil {
		return nil, err
	}

	resp := &BalanceResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, apperrors.NewResponseFormatError("BalanceResponse")
	}
	return resp, nil
}
