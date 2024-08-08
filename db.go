package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase(dbFile string) {
	if !fileExists(dbFile) {
		createDatabase(dbFile)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

func createDatabase(dbFile string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

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

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database and table created successfully")
}
