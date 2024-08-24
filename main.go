package main

import (
	"fmt"
	"net/http"

	"myproject/pkg/config"
	handler "myproject/pkg/handlers"
	auth "myproject/pkg/midware"
	"myproject/pkg/router"
	"myproject/pkg/server"
	"myproject/pkg/service"
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

	taskRepository := storage.NewTaskRepository(DB)                                   //новый экземпляр репозитория задач
	taskService := service.NewTaskService(taskRepository)                             //новый экземпляр сервиса задач
	taskHandler := handler.NewTaskHandler(taskService, taskRepository)                //новый экземпляр обработчика задач
	dateHandler := handler.NewDateHandler(taskService)                                //новый экземпляр обработчика дат
	authHandler := auth.NewAuthHandler(config.TodoPassword, []byte(config.SecretKey)) //новый экземпляр обработчика аутентификации
	authMiddleware := auth.NewAuthMiddleware()
	//новый экземпляр маршрутизатора
	router := router.NewRouter(taskHandler, dateHandler, authHandler, authMiddleware)

	//для шага 8 - обработчик для аутентификации
	http.HandleFunc("/api/signin", router.SwitchTaskHandler)

	// маршрутизатор для обслуживания статических файлов (интерфейс пользователя, отображение контента, пр)
	fs := http.FileServer(http.Dir(config.WebDir))
	http.Handle("/", fs)

	//API-обработчик для адреса "/api/nextdate"
	http.HandleFunc("/api/nextdate", router.SwitchTaskHandler)

	// для шагов 4 (добавление), 6 (редактирование) и 7(удаление) задачи
	//API-обработчик для адреса "/api/task" (по ТЗ все методы с авторизацией)
	http.HandleFunc("/api/task", authMiddleware.Middleware(router.SwitchTaskHandler))

	// для шага 5 (список задач, поиск)
	//API-обработчик для адреса "/api/tasks" (по ТЗ не все задачи с авторизацией)
	http.HandleFunc("/api/tasks", router.SwitchTaskHandler)

	// для шага 7(выполнение задачи)
	//API-обработчик для адреса "/api/task/done" по ТЗ с авторизацией
	http.HandleFunc("/api/task/done", authMiddleware.Middleware(router.SwitchTaskHandler))

	// запуск сервера
	server.Start()

}
