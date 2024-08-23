package main

import (
	//"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// обрабатывает запрос на выполнение задачи
func handleTaskDone(w http.ResponseWriter, r *http.Request) {
	//получение поля id
	id := r.URL.Query().Get("id")
	//получение записи из БД по полю id
	task, err := getTaskById(id)
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}
	//Одноразовая задача с пустым полем `repeat` удаляется.
	switch task.Repeat {
	case "":
		err = DeleteTask(id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
	default:
		//переводим дату задачи в формат time.Time
		now, err := time.Parse(DATE_FORMAT, task.Date)
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
		//обновление измененных полей задачи
		err = updateTask(task)
		if err != nil {
			http.Error(w, `{"error": "не обновлено"}`, http.StatusNotImplemented)
			return
		}

	}
	/*
		if task.Repeat == "" {
			err = DeleteTask(id)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
				return
			}
		//Для периодической задачи изменяется дата следующего выполнения
		} else {
			//переводим дату задачи в формат time.Time
			now, err := time.Parse(DATE_FORMAT, task.Date)
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
			//обновление измененных полей задачи
			err = updateTask(task)
			if err != nil {
				http.Error(w, `{"error": "не обновлено"}`, http.StatusNotImplemented)
				return
			}

		}
	*/
	json.NewEncoder(w).Encode(map[string]interface{}{})
	//respondWithJSON(w, interface{}, http.StatusOK)
}

// Обработчик DELETE-запроса
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	//получение поля id из запроса
	id := r.URL.Query().Get("id")
	//удаление записи
	err := DeleteTask(id)
	//обработка ошибки
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	//отправка пустого JSON-объекта в качестве ответа на HTTP-запрос
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// удаляет запись из БД с помощью SQL-запроса DELETE
func DeleteTask(id string) error {
	//проверка, что id содержит только цифры (и что не пустая строка)
	if !isNumeric(id) {
		return errors.New("invalid task ID")
	}
	result, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	//проверка с помощью RowsAffected()
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	//если затронуто более 1 записи, возвращается ошибка
	if rowsAffected != 1 {
		return fmt.Errorf("failed to delete task with id %s: %d rows affected", id, rowsAffected)
	}

	return nil
}
