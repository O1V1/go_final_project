package main

import (
	//"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	//"strconv"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

func handleTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := getTaskById(id)
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}
	if task.Repeat == "" {
		err = DeleteTask(id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		//переводим дату задачи в формат time.Time
		now, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "time parse error}`, http.StatusNotImplemented)
			return
		}
		//передаем запрос с двумя одинаковыми датами
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			respondWithError(w, "Invalid repeat rule", http.StatusBadRequest)
			return
		}
		task.Date = nextDate
		err = UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error": "не обновлено"}`, http.StatusNotImplemented)
			return
		}

	}

	json.NewEncoder(w).Encode(map[string]interface{}{})
	//respondWithJSON(w, interface{}, http.StatusOK)
}

// Обработчик DELETE-запроса
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = DeleteTask(id)

	if err != nil {
		http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		//respondWithError(w, "Invalid repeat rule", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{})

}

func DeleteTask(id string) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}
