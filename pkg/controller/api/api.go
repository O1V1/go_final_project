package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/O1V1/go_final_project/pkg/controller/config"
	handler "github.com/O1V1/go_final_project/pkg/controller/handlers"
	auth "github.com/O1V1/go_final_project/pkg/controller/middleware"
	"github.com/O1V1/go_final_project/pkg/service"
	"github.com/O1V1/go_final_project/pkg/storage"
)

func NewServer(db *sql.DB) {
	taskStorage := storage.NewTaskStorage(db)                                          //новый экземпляр репозитория задач
	taskService := service.NewTaskService(taskStorage)                                 //новый экземпляр сервиса задач
	taskHandler := *handler.NewTaskHandler(taskService, taskStorage)                   //новый экземпляр обработчика задач
	dateHandler := *handler.NewDateHandler(taskService)                                //новый экземпляр обработчика дат
	authHandler := *auth.NewAuthHandler(config.TodoPassword, []byte(config.SecretKey)) //новый экземпляр обработчика аутентификации
	authFunc := auth.NewAuthMiddleware()

	// маршрутизатор для обслуживания статических файлов (интерфейс пользователя, отображение контента, пр)
	fs := http.FileServer(http.Dir(config.WebDir))
	http.Handle("/", fs)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		//маршрутизатор для обслуживания остальных запросов
		switch r.URL.Path {
		case "/api/signin": //обработчик для аутентификации step_8
			authHandler.SigninHandler(w, r)
		case "/api/nextdate": //обработчик для адреса "/api/nextdate" step_3
			dateHandler.NextDateHandler(w, r)
		case "/api/task":
			switch r.Method {
			case http.MethodPost: //добавление задачи, step_4
				authFunc.Middleware(taskHandler.HandlePostTask)(w, r)
			case http.MethodGet: //получение задачи по id, step_6
				authFunc.Middleware(taskHandler.HandleGetTask)(w, r)
			case http.MethodPut: //редактирование задачи, step_6
				authFunc.Middleware(taskHandler.HandleUpdateTask)(w, r)
			case http.MethodDelete: //удаление задачи, step_7
				authFunc.Middleware(taskHandler.HandleDeleteTask)(w, r)
			default:
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
		case "/api/tasks": //получение списка задач, step_5
			authFunc.Middleware(taskHandler.HandleGetList)(w, r)
		case "/api/task/done": //выполнение задачи, step_7
			authFunc.Middleware(taskHandler.HandleTaskDone)(w, r)

		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})
}

func StartServer(port string) {
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
