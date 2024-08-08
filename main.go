package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	//"path/filepath"
)

func main() {

	// Получаем значение переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	}

	// Указываем директорию, из которой будут обслуживаться статические файлы
	webDir := "./web"

	// Получаем значение переменной окружения TODO_DBFILE
	// иначе формируем путь из текущей директории
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // Файл базы данных по умолчанию
	}

	// Инициализируем базу данных
	initDatabase(dbFile)

	// Подключение к БД
	fmt.Printf("Using database file: %s\n", dbFile)

	// Настраиваем маршрутизатор для обслуживания статических файлов
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// Запускаем сервер
	fmt.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
