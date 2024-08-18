package main

import (
	//"database/sql"
	"encoding/json"

	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Task структура для объекта базы данных
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

/*/ обработчик по методу запроса
func switchTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/api/task" {
			handlePostTask(w, r)
		} else if r.URL.Path == "/api/task/done" {
			handleTaskDone(w, r)
		}

	case http.MethodGet:
		if r.URL.Path == "/api/task" {
			handleGetTask(w, r)
		} else if r.URL.Path == "/api/tasks" {
			handleGetList(w, r)
		}
	case http.MethodPut:
		handleUpdateTask(w, r)

	case http.MethodDelete:
		handleDeleteTask(w, r)

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
} */

// Обработчик POST-запроса, добавление новой задачи в базу данных
func handlePostTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	//фиксирую текущую дату
	now := time.Now()

	//десериализация реквеста в структуру, обработка ошибки и выход в случае ошибки
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//проверка, что заголовок не пустой
	if task.Title == "" {
		respondWithError(w, "Title is required", http.StatusBadRequest)
		return
	}

	//если дата не заполнена, заполняем текущей датой
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	//превращаем дату в объект time.Time, если ошибка - обрабатываем и выходим
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		respondWithError(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	//Если указана прошедшая дата, используем функцию из шага 3
	//если нет праавила повторения, то ставим текущую дату
	if task.Date < now.Format("20060102") {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				respondWithError(w, "Invalid repeat rule", http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	idStr, err := addTaskToDB(task)
	if err != nil {
		respondWithError(w, "Failed to add task", http.StatusInternalServerError)
		return
	}
	//idStr := strconv.Itoa(id)

	respondWithJSON(w, TaskResponse{ID: idStr, Error: ""}, http.StatusOK)
}

// addTaskToDatabase Добавляет задаччу в БД, возвращает id добавленной записи или ошибку
func addTaskToDB(task Task) (string, error) {

	//подготовка запроса
	stmt, err := DB.Prepare("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	//при возврате функции запрос будет закрыт
	defer stmt.Close()

	// исполнение подготовленного запроса с полями задачи task в качестве аргументов
	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	// получаем id последней удачно добавленной записи
	// при ошибках выполняется возврат функции
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	//id получена в формате int64, п
	idStr := strconv.Itoa(int(id))
	return idStr, nil
}

// общий формат респонда в формате JSON
// respondWithJSON записывает в переменную w нужные заголовки, передаваемые данные и код ответа
func respondWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	//
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

// респонд в случае ошибки
// respondWithError записывает в w сообщение и код ошибки
func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, TaskResponse{Error: message}, code)
}
