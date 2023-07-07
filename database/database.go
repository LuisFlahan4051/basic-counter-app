package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DATABASE_NAME = "database"
)

func GetConnection(databaseName string) *sql.DB {
	database, err := sql.Open("sqlite3", fmt.Sprintf("./database/%s.db", databaseName))
	if err != nil {
		panic(err)
	}
	return database
}

func InitDatabaseIfNotExists() {
	database := GetConnection(DATABASE_NAME)
	defer database.Close()

	createTableIncomes := `CREATE TABLE IF NOT EXISTS incomes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value REAL CHECK(value > 0) NOT NULL, 
		type VARCHAR(64) NOT NULL, 
		created_at DATE NOT NULL);`
	statement, err := database.Prepare(createTableIncomes)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	createTableExpenses := `CREATE TABLE IF NOT EXISTS expenses(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value REAL CHECK(value > 0) NOT NULL, 
		type VARCHAR(64) NOT NULL, 
		created_at DATE NOT NULL);`
	statement, err = database.Prepare(createTableExpenses)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}
