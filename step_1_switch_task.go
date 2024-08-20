package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// обработчик запросов по методу
// перенаправляет запрос в специализированный обработчик для каждой задачи
func switchTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//метод POST
	case http.MethodPost:
		//выбор обработчика в зависимости от адреса
		if r.URL.Path == "/api/task" {
			handlePostTask(w, r) //добавление задачи, файл step_4
		} else if r.URL.Path == "/api/task/done" {
			handleTaskDone(w, r) //выполнение задачи, файл step_7
		}
	//метод GET
	case http.MethodGet:
		//выбор обработчика в зависимости от адреса
		if r.URL.Path == "/api/task" {
			handleGetTask(w, r) //получение задачи по id, файл step_6
		} else if r.URL.Path == "/api/tasks" {
			////по ТЗ с авторизацией только получение списка задач
			auth(handleGetList)(w, r) //получение списка задач, файл step_5
		}
	//метод PUT
	case http.MethodPut:
		handleUpdateTask(w, r) //редактирование задачи, файл step_6
	//метод DELETE
	case http.MethodDelete:
		handleDeleteTask(w, r) //удаление задачи, файл step_7

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
