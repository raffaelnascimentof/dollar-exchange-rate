package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const BR_UTC = -3 * 60 * 60
const YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"

func openConnection() (*sql.DB, error) {
	var connInformation = "host=localhost port=5432 user=root password=root dbname=dollar-quotation sslmode=disable"
	database, err := sql.Open("postgres", connInformation)
	if err != nil {
		fmt.Println("not found database")
	}

	return database, nil
}

func InitDB() {
	database, err := openConnection()
	if err != nil {
		fmt.Println("Connection failed : " + err.Error())
	}
	defer database.Close()

	sql := `CREATE TABLE IF NOT EXISTS quotation (id SERIAL PRIMARY KEY , value TEXT, data TEXT)`
	_, err = database.Exec(sql)
	if err != nil {
		fmt.Println("Table 'quotation' not created: " + err.Error())
		return
	}

	fmt.Println("Table 'quotation' created.")
}

func InsertQuotationValue(quotationValue string) {
	database, err := openConnection()
	if err != nil {
		fmt.Println("Not found database")
	}
	defer database.Close()

	stmt, err := database.Prepare(`INSERT INTO quotation (value, data) VALUES ($1, $2)`)
	if err != nil {
		fmt.Println("Erro prepare SQL")
	}
	defer stmt.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	_, err = stmt.ExecContext(ctx, quotationValue, getCurrentDate())
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
