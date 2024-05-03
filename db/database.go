package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func OpenConnection() (*sql.DB, error) {
	var connInformation = "host=localhost port=5432 user=root password=root dbname=dollar-quotation sslmode=disable"
	database, err := sql.Open("postgres", connInformation)
	if err != nil {
		fmt.Println("not found database")
	}

	return database, nil
}

func InitDB() {
	database, err := OpenConnection()
	if err != nil {
		log.Println("Connection failed : " + err.Error())
	}
	defer database.Close()

	sql := `CREATE TABLE IF NOT EXISTS quotation (id SERIAL PRIMARY KEY ,name TEXT, code TEXT, codein TEXT, value TEXT, data TEXT)`
	_, err = database.Exec(sql)
	if err != nil {
		log.Println("Table 'quotation' not created: " + err.Error())
		return
	}
}
