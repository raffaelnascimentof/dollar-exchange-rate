package dto

import "github.com/raffaelnascimentof/dollar-exchange-rate/server/domain"

type QuotationDTO struct {
	Bid string `json:"bid"`
}

func ToDTO(quotationDomain *domain.QuotationDomain) *QuotationDTO {
	quotationDTO := QuotationDTO{
		Bid: quotationDomain.Bid,
	}

	return &quotationDTO
}
