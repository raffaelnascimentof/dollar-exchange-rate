package dto

import "github.com/raffaelnascimentof/dollar-exchange-rate/client/response"

type QuotationDTO struct {
	Value string `json:"dolar"`
}

func ToDTO(quotationResponse *response.QuotationResponse) *QuotationDTO {
	quotationDTO := QuotationDTO{
		Value: quotationResponse.Dolar,
	}

	return &quotationDTO
}
