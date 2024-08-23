package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Определение директории для статических файлов, объявление перем для рабочей директории
// чтение переменных окружения
var (
	webDir       = "./web"
	workingDir   string
	todoPort     = os.Getenv("TODO_PORT")
	todoDBFile   = os.Getenv("TODO_DBFILE")
	todoPassword = os.Getenv("TODO_PASSWORD")
)

// инициализация переменной workingDir значением текущего рабочего каталога
func init() {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v\n", err)
	}
}

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

	// для шага 5 (список задач, поиск)
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
	if todoPort == "" {
		return "7540" //значение по умочанию
	}
	//проверка, что получено цифровое значение для порта
	var err error
	_, err = strconv.Atoi(todoPort)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	return todoPort
}

// getDBFile() возвращает путь к файлу базы данных из переменной окружения TODO_DBFILE, либо значение по умолчанию
func getDBFile() string {
	//если переменная окружения - пустая строка, формируем путь по умолчанию (scheduler.db в текущей директории)
	if todoDBFile == "" {
		//dbFile = filepath.Join(workingDir, "scheduler.db")
		return filepath.Join(workingDir, "scheduler.db")
	}
	return todoDBFile
}
