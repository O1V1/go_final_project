package main

import (
	"fmt"

	"myproject/pkg/api"
	"myproject/pkg/config"
	"myproject/pkg/server"
	"myproject/pkg/storage"
)

func main() {
	//инициализация различных параметоров для работы программы
	config.Init()

	//создается экземпляр структуры databaseRepository с помощью конструктора
	dbRep := storage.NewDatabaseRepository(nil)
	//вызывается метод InitDatabase для открытия базы данных.
	dbRep.InitDatabase(config.DBFile)
	//выделяю переменную для удобства
	DB := dbRep.DB()

	defer DB.Close()

	// Вывод сообщения о подключении к БД
	fmt.Printf("Using database file: %s\n", config.DBFile)

	// Настройка API-обработчиков
	api.SetupAPI(DB)

	// запуск сервера
	server.Start()

}
