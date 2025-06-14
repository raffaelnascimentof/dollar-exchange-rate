package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/raffaelnascimentof/dollar-exchange-rate/db"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/domain"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/dto"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/quotationrepository"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/response"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/responseerror"
)

const URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	db.InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	quotation, responseError := getQuotation(ctx)
	if responseError.Message != "" {
		jsonResponse, _ := json.Marshal(responseError)
		http.Error(w, string(jsonResponse), responseError.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation)
}

func getQuotation(ctx context.Context) (*dto.QuotationDTO, responseerror.ResponseError) {
	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, responseerror.CreateError("Internal server error", http.StatusInternalServerError)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, responseerror.CreateError("request timeout reached", http.StatusRequestTimeout)
		}
		return nil, responseerror.CreateError("Internal server error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, responseerror.CreateError("Internal server error", http.StatusInternalServerError)
	}

	var quotationResponse response.QuotationResponse
	error := json.Unmarshal(responseData, &quotationResponse)
	if error != nil {
		return nil, responseerror.CreateError("Internal server error", http.StatusInternalServerError)
	}

	quotationDomain := domain.ToDomain(quotationResponse)
	quotationrepository.Save(quotationDomain)

	return dto.ToDTO(quotationDomain), responseerror.CreateError("", http.StatusOK)
}
