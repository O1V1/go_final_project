package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Указываем директорию для статических файлов
var webDir = "./web"

func main() {

	// получение номера порта
	port := getPort()

	//получение пути к базе данных
	dbFile := getDBFile()

	// Инициализия базы данных
	initDatabase(dbFile)

	defer DB.Close()

	// Вывод сообщения о подключении к БД
	fmt.Printf("Using database file: %s\n", dbFile)

	//для шага 8 - обработчик для аутентификации
	http.HandleFunc("/api/signin", signinHandler)

	// маршрутизатор для обслуживания статических файлов (интерфейс пользователя, отображение контента, пр)
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	//для шага 3 - запросы к функции NextDate
	//API-обработчик для адреса "/api/nextdate"
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// для шагов 4 (добавление), 6 (редактирование) и 7(удаление) задачи
	//API-обработчик для адреса "/api/task"
	//по ТЗ все методы с авторизацией
	http.HandleFunc("/api/task", auth(switchTaskHandler))

	// для шага 5 (список задач и поиск)
	//API-обработчик для адреса "/api/tasks"
	//по ТЗ с авторизацией только получение списка задач
	http.HandleFunc("/api/tasks", switchTaskHandler)

	// для шага 7(выполнение задачи)
	//API-обработчик для адреса "/api/task/done"
	//по ТЗ с авторизацией
	http.HandleFunc("/api/task/done", auth(switchTaskHandler))

	// Запускаем сервер
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

// getPort() возвращает порт из переменной окружения, либо значение по умолчанию
func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" //значение по умочанию
	} else {
		//проверка, что получено цифровое значение для порта
		var err error
		_, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
	return port
}

// getDBFile() возвращает путь к файлу базы данных из переменной окружения TODO_DBFILE, либо значение по умолчанию
func getDBFile() string {
	//Получаем значение переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")
	//если переменная окружения - пустая строка, формируем путь по умолчанию (scheduler.db в текущей директории)
	if dbFile == "" {
		wd, err := os.Getwd()
		//os.Executable() - не подошла, так как возвращала временную директорию
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(wd, "scheduler.db")
	}
	return dbFile
}
