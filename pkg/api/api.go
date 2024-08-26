package api

import (
	"database/sql"
	"net/http"

	"github.com/O1V1/go_final_project/pkg/config"
	handler "github.com/O1V1/go_final_project/pkg/handlers"
	auth "github.com/O1V1/go_final_project/pkg/middleware"
	"github.com/O1V1/go_final_project/pkg/router"
	"github.com/O1V1/go_final_project/pkg/service"
	"github.com/O1V1/go_final_project/pkg/storage"
)

func NewServer(db *sql.DB) {
	taskRepository := storage.NewTaskRepository(db)                                   //новый экземпляр репозитория задач
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
}
