package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/raffaelnascimentof/dollar-exchange-rate/config"
)

type QuotationResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type QuotationDTO struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	CodeIn string `json:"codein"`
	Bid    string `json:"bid"`
}

type ResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"statuscode"`
}

func main() {
	config.InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	quotation, responseError := getQuotation(w, r, ctx)
	if responseError.Message != "" {
		jsonResponse, _ := json.Marshal(responseError)
		http.Error(w, string(jsonResponse), responseError.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation)
}

func getQuotation(w http.ResponseWriter, r *http.Request, ctx context.Context) (*QuotationDTO, ResponseError) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, createErrorResponse("Internal server error", http.StatusInternalServerError)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, createErrorResponse("request timeout reached", http.StatusRequestTimeout)
		}
		return nil, createErrorResponse("Internal server error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, createErrorResponse("Internal server error", http.StatusInternalServerError)
	}

	var quotationResponse QuotationResponse
	error := json.Unmarshal(responseData, &quotationResponse)
	if error != nil {
		return nil, createErrorResponse("Internal server error", http.StatusInternalServerError)
	}

	quotationDTO := convertToQuotationDTO(quotationResponse)
	saveQuotationExchange(quotationDTO.Bid)

	return quotationDTO, createErrorResponse("", http.StatusOK)
}

func createErrorResponse(err string, code int) ResponseError {
	var responseError = ResponseError{
		Message: err,
		Code:    code,
	}
	return responseError
}

func saveQuotationExchange(exchangeQuotation string) {
	config.InsertQuotationValue(exchangeQuotation)
}

func convertToQuotationDTO(quotationResponse QuotationResponse) *QuotationDTO {
	quotationDTO := QuotationDTO{
		Name:   quotationResponse.USDBRL.Name,
		Code:   quotationResponse.USDBRL.Code,
		CodeIn: quotationResponse.USDBRL.Codein,
		Bid:    quotationResponse.USDBRL.Bid,
	}

	return &quotationDTO
}
