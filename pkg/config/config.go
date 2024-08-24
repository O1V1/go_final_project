package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	DATE_FORMAT = "20060102"
	TASKS_LIMIT = "15"
)

var (
	// чтение переменных окружения
	todoPort     = os.Getenv("TODO_PORT")
	todoDBFile   = os.Getenv("TODO_DBFILE")
	TodoPassword = os.Getenv("TODO_PASSWORD")

	//переменные для работы приложения
	Port       string
	DBFile     string
	workingDir string

	WebDir = "./web" // Определение директории для статических файлов

	// формируем переменную секретного ключа для подписи токена
	// строка получена с помощью openssl rand -base64 32
	SecretKey = []byte("9wsz2ew8lF2pxS4LEg1pHxq9jVhztkKQD5O/5OfvPdE=")
)

func Init() {
	// получение номера порта
	Port = getPort()
	//получение пути к базе данных
	DBFile = getDBFile()

	// инициализация переменной workingDir значением текущего рабочего каталога
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v\n", err)
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
		return filepath.Join(workingDir, "scheduler.db")
	}
	return todoDBFile
}
