package quotationrepository

import (
	"context"
	"log"
	"time"

	"github.com/raffaelnascimentof/dollar-exchange-rate/db"
	"github.com/raffaelnascimentof/dollar-exchange-rate/server/domain"
)

const BR_UTC = -3 * 60 * 60
const YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"

func Save(quotationDomain *domain.QuotationDomain) {
	database, err := db.OpenConnection()
	if err != nil {
		log.Println("Not found database")
	}
	defer database.Close()

	stmt, err := database.Prepare(`INSERT INTO quotation (name, code, codein, value, data) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Erro prepare SQL")
	}
	defer stmt.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	_, err = stmt.ExecContext(ctx, quotationDomain.Name, quotationDomain.Code, quotationDomain.CodeIn, quotationDomain.Bid, getCurrentDate())
	if err != nil {
		panic(err.Error())
	}
}

func getCurrentDate() string {
	timeZoneBRT := time.FixedZone("BRT", BR_UTC)
	currentDate := time.Now().In(timeZoneBRT)

	currentDateFormatted := currentDate.Format(YYYY_MM_DD_HH_MM_SS)

	return currentDateFormatted
}
