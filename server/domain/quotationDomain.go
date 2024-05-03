package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/raffaelnascimentof/dollar-exchange-rate/server/response"
)

type QuotationDomain struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	CodeIn string `json:"codein"`
	Bid    string `json:"bid"`
}

func ToDomain(quotationResponse response.QuotationResponse) *QuotationDomain {
	quotationDomain := QuotationDomain{
		Name:   quotationResponse.USDBRL.Name,
		Code:   quotationResponse.USDBRL.Code,
		CodeIn: quotationResponse.USDBRL.Codein,
		Bid:    quotationResponse.USDBRL.Bid,
	}

	return &quotationDomain
}

func (q QuotationDomain) Value() (driver.Value, error) {
	return json.Marshal(q)
}
