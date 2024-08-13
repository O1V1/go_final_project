package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func initDatabase(dbFile string) {

	var err error
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if !fileExists(dbFile) {
		createDatabase()
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// func createDatabase(dbFile string) {
func createDatabase() {
	var err error
	/* DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer DB.Close() */

	createTableSQL := `
    CREATE TABLE scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date CHAR(8) NOT NULL DEFAULT "",
        title VARCHAR(256) NOT NULL DEFAULT "",
        comment TEXT,
        repeat VARCHAR(128) NOT NULL DEFAULT ""
    );
    CREATE INDEX idx_date ON scheduler (date);
    `

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database and table created successfully")
}

/* func fileExists() bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbile)

	return err != nil
} */
