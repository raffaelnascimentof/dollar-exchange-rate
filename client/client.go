package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raffaelnascimentof/dollar-exchange-rate/client/dto"
	"github.com/raffaelnascimentof/dollar-exchange-rate/client/response"
)

const URL = "http://localhost:8080/cotacao"
const PERMISSION_FILE = 0644

func main() {
	quotationResponse := getQuotation()

	var file = createQuotationHistoryFile()
	defer file.Close()

	quotationDTO := dto.ToDTO(quotationResponse)
	quotationJson, _ := json.Marshal(&quotationDTO)

	_, err := file.WriteString(string(quotationJson) + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever informações da cotação: %v\n", err)
	}
}

func getQuotation() *response.QuotationResponse {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar requisição: %v\n", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "Tempo limite excedido: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "Erro ao executar requisição: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	var quotationResponse response.QuotationResponse
	err = json.Unmarshal(data, &quotationResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
	}

	return &quotationResponse
}

func createQuotationHistoryFile() os.File {
	var fileName = "quotation-history.txt"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
		}

		return *file
	} else {
		file, err := os.OpenFile(fileName, os.O_APPEND, PERMISSION_FILE)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir arquivo: %v\n", err)
		}
		return *file
	}
}
