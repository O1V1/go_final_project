package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// интерфейс для работы с бд
type Storage interface {
	Init(dbFile string)
	//Create()
	DB() *sql.DB
}

// структура для интерфейса
type storage struct {
	db *sql.DB
}

// конструктор нового экземпляра структуры
func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

// метод возвращает указатель на экземпляр *sql.DB
func (r *storage) DB() *sql.DB {
	return r.db
}

// метод initDatabase готовит базу данных к использованию
func (r *storage) Init(dbFile string) {
	var err error
	//установка соединения с базой данных dbFile
	r.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	//создание файла базы данных в случае его отсутствия
	if !fileExists(dbFile) {
		createDatabase(r.db)
	}
}

// метод createDatabase создает базу данных с таблицей scheduler
func createDatabase(db *sql.DB) {
	//формируется текст команды для создания и индексирования таблицы
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
	//исполнение команды, обработтка ошибки
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	//сообщение об успешном выполнении задачи
	fmt.Println("Database and table created successfully")
}

// fileExists проверяет, существует ли файл filename
func fileExists(filename string) bool {
	//os.Stat ищет файл либо по абсол. пути, если он указан в filename, либо в текущей раб директории
	info, err := os.Stat(filename)
	//Если err равна os.ErrNotExist, файл не найден, возвращается false
	if os.IsNotExist(err) {
		return false
	}
	//если filename не является директорией, возращается true
	return !info.IsDir()
}
