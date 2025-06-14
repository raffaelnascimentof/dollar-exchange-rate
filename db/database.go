package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenConnection() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./quotation.db")
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

	sql := `CREATE TABLE IF NOT EXISTS quotation (id INTEGER PRIMARY KEY AUTOINCREMENT ,name TEXT, code TEXT, codein TEXT, value TEXT, data TEXT)`
	_, err = database.Exec(sql)
	if err != nil {
		log.Println("Table 'quotation' not created: " + err.Error())
		return
	}
}
