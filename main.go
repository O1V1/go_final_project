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

	/* для проверки
	"20231106", "m 13", "20240213"},
		{"20240120", "m 40,11,19", ""},
		{"20240116", "m 16,5", "20240205"},
		{"20240126", "m 25,26,7", "20240207"},
		{"20240409", "m 31", "20240531"},
		{"20240329", "m 10,17 12,8,1", "20240810"},
		{"20230311", "m 07,19 05,6", "20240507"},
		{"20230311", "m 1 1,2", "20240201"},
		{"20240127", "m -1", "20240131"},
		{"20240222", "m -2", "20240228"},
		{"20240222", "m -2,-3", ""},
		{"20240326", "m -1,-2", "20240330"},
		{"20240201", "m -1,18", "20240218"},

		{"20240125", "w 1,2,3", "20240129"},
		{"20240126", "w 7", "20240128"},
		{"20230126", "w 4,5", "20240201"},

		{"20230226", "w 8,4,5", ""},


	now := time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)
	fmt.Println("время")
	a, _ := NextDate(now, "20240125", "w 1,2,3")
	fmt.Println(a)
	b, _ := NextDate(now, "20240126", "w 7")
	fmt.Println(b)
	c, _ := NextDate(now, "20240201", "m -1,18")
	fmt.Println(c)
	d, _ := NextDate(now, "20230126", "w 4,5")
	fmt.Println(d)
	fmt.Println("расчет")

	fmt.Println(a, b, c, d)

	// */

	// Запускаем сервер
	fmt.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
