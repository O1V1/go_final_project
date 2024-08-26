package main

import (
	"fmt"

	"github.com/O1V1/go_final_project/pkg/api"
	"github.com/O1V1/go_final_project/pkg/config"
	"github.com/O1V1/go_final_project/pkg/server"
	"github.com/O1V1/go_final_project/pkg/storage"
)

func main() {
	//инициализация различных параметоров для работы программы
	config.Init()

	//создается экземпляр структуры storage с помощью конструктора
	dbRep := storage.NewStorage(nil)
	//вызывается метод InitDatabase для открытия базы данных.
	dbRep.Init(config.DBFile)
	//выделяю переменную для удобства
	db := dbRep.DB()

	defer db.Close()

	// Вывод сообщения о подключении к БД
	fmt.Printf("Using database file: %s\n", config.DBFile)

	// Настройка API-обработчиков
	api.NewServer(db)

	// запуск сервера
	server.Start()

}
