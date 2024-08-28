package main

import (
	"fmt"

	"github.com/O1V1/go_final_project/pkg/controller/api"
	"github.com/O1V1/go_final_project/pkg/controller/config"
	"github.com/O1V1/go_final_project/pkg/storage"
)

func main() {
	//инициализация различных параметоров для работы программы
	config.Init()

	//получение указателя на открытое подключение к базе данных, создание бд при необходимости
	db := storage.Init(config.DBFile)
	defer db.Close()
	fmt.Printf("Using database file: %s\n", config.DBFile)

	// Настройка и маршрутизация API-обработчиков
	api.NewServer(db)

	// запуск сервера
	api.StartServer(config.Port)
}
