package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	//"path/filepath"
)

// Указываем директорию, из которой будут обслуживаться статические файлы
var webDir = "./web"

func main() {

	// Получаем значение переменной окружения TODO_PORT
	port := getPort()
	/* os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	} */

	// Получаем значение переменной окружения TODO_DBFILE
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

	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Запускаем сервер
	/* fmt.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil)) */

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

// Функция для получения порта из переменной окружения
func getPort() string {
	// По умолчанию используем порт 7540
	defaultPort := "7540"

	// Проверяем переменную окружения TODO_PORT и что это цифры
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return strconv.Itoa(port)
		}
	}

	return defaultPort
}
