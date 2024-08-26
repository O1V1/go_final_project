package router

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	handler "github.com/O1V1/go_final_project/pkg/handlers"
	auth "github.com/O1V1/go_final_project/pkg/middleware"
)

type Router struct {
	taskHandler    *handler.TaskHandlerImpl
	dateHandler    *handler.DateHandlerImpl
	authHandler    *auth.AuthHandlerImpl
	authMiddleware *auth.AuthMiddleware
}

// конструктор для маршрутизации обработчиков
func NewRouter(taskHandler *handler.TaskHandlerImpl, dateHandler *handler.DateHandlerImpl, authHandler *auth.AuthHandlerImpl, authMiddleware *auth.AuthMiddleware) *Router {
	return &Router{taskHandler: taskHandler, dateHandler: dateHandler, authHandler: authHandler, authMiddleware: authMiddleware}
}

// метод структуры router для маршрутизации обработчиков
// перенаправляет запрос в специализированный обработчик для каждой задачи
func (r *Router) SwitchTaskHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	//метод POST
	case http.MethodPost:
		//выбор обработчика в зависимости от адреса
		if req.URL.Path == "/api/task" {
			r.taskHandler.HandlePostTask(w, req) //добавление задачи, файл step_4
		} else if req.URL.Path == "/api/task/done" {
			r.taskHandler.HandleTaskDone(w, req) //выполнение задачи, файл step_7
		} else if req.URL.Path == "/api/signin" {
			r.authHandler.SigninHandler(w, req)
		}
	//метод GET
	case http.MethodGet:
		//выбор обработчика в зависимости от адреса
		if req.URL.Path == "/api/task" {
			r.taskHandler.HandleGetTask(w, req) //получение задачи по id, файл step_6
		} else if req.URL.Path == "/api/tasks" {
			////по ТЗ с авторизацией только получение списка задач
			//auth(handleGetList)(w, r) //получение списка задач, файл step_5
			r.authMiddleware.Middleware(r.taskHandler.HandleGetList)(w, req)
		} else if req.URL.Path == "/api/nextdate" {
			r.dateHandler.NextDateHandler(w, req)
		}
	//метод PUT
	case http.MethodPut:
		r.taskHandler.HandleUpdateTask(w, req) //редактирование задачи, файл step_6
	//метод DELETE
	case http.MethodDelete:
		r.taskHandler.HandleDeleteTask(w, req) //удаление задачи, файл step_7

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
