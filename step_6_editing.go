package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type IDTaskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
	Error   string `json:"error,omitempty"`
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	//извлечение id из параметров запроса
	id := r.URL.Query().Get("id")
	//получение таски по запрошенному id, обработка ошибки
	task, err := getTaskById(id)
	if err != nil {
		//если задача не получена, возврат респонса с полем error
		response := IDTaskResponse{
			Error: "Задача не найдена",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	//при успешном получении задачи, запись в ответ
	response := IDTaskResponse{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	json.NewEncoder(w).Encode(response)

}

// getTaskById возвращает задачу из БД по заданному id, либо ошибку
func getTaskById(id string) (Task, error) {
	var task Task
	//если в id не только цифры, то не беспокоим БД
	if !isNumeric(id) {
		return Task{}, errors.New("ID must be a number")
	}
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()
	//nowFormat := now.Format(DATE_FORMAT)
	// извлечение тела запроса в структуру task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//проверка формата поля ID
	if _, err = getTaskById(task.ID); err != nil {
		respondWithError(w, "Invalid ID formal", http.StatusBadRequest)
		return
	}

	//проверки полей Date, Title task на соответствие требованиям БД
	task, err = prepareTaskTitleAndDate(task, now)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//обновление параметров задачи, обработка ошибки
	err = updateTask(task)
	if err != nil {
		respondWithError(w, "Task update failed", http.StatusNotImplemented)
		return
	}
	//после успешного обновления задачи формируется ответ на http запрос
	//в виде пустой структуры Task в формате JSON
	json.NewEncoder(w).Encode(Task{})
}

// обновляет запись в базе данных с помощью SQL-запроса UPDATE
func updateTask(task Task) error {
	_, err := DB.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	return nil
}
