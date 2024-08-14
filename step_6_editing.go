package main

import (
	"database/sql"
	"encoding/json"

	//"encoding/json"
	"net/http"

	"strconv"
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

	id := r.URL.Query().Get("id")
	task, err := getTaskById(id)

	if err != nil {
		//http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		//return
		response := IDTaskResponse{
			Error: "Задача не найдена",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := IDTaskResponse{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	json.NewEncoder(w).Encode(response)

	//fmt.Println(task)
	//fmt.Println(json.NewEncoder(w).Encode(task))

	/*
		if err != nil {
			//http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			//respondWithError(w, "Failed to add task", http.StatusNotFound)
			respondWithJSON(w, IDTaskResponse{ID: task.ID, Error: ""}, http.StatusNotFound)
			return
		}

		respondWithJSON(w, IDTaskResponse{ID: task.ID, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat, Error: ""}, http.StatusOK)

		//json.NewEncoder(w).Encode(task)
	} */
}

func getTaskById(id string) (Task, error) {
	var task Task
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		//http.Error(w, `{"error": "Некорректные данные"}`, http.StatusBadRequest)
		return
	}
	//проверка, что заголовок не пустой
	if task.Title == "" || task.ID == "" {
		respondWithError(w, "INVALID FORMAT", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		respondWithError(w, "INVALID FORMAT", http.StatusBadRequest)
		return
	}

	date := ""
	row := DB.QueryRow("SELECT date FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	err = row.Scan(&date)
	if err != nil {
		respondWithError(w, "INVALID FORMAT", http.StatusBadRequest)
		return
	}

	//если дата не заполнена, заполняем текущей датой
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	//превращаем дату в объект time.Time, если ошибка - обрабатываем и выходим
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		respondWithError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	/*
		if !strings.Contains("1234567890", task.ID) || task.ID == "" {
			respondWithError(w, "Invalid date format", http.StatusBadRequest)
			return
		} */

	//Если указана прошедшая дата, используем функцию из шага 3
	//если нет праавила повторения, то ставим текущую дату
	//if date.Before(now) {
	//с датами не получалось before, переделала на сравнение строк
	// можно попробовать truncate обе даты, если время будет
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

	//err = updateTask(task)
	_, err = DB.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error": "не обновлено"}`, http.StatusNotImplemented)
		return
	}
	json.NewEncoder(w).Encode(Task{})

}

/*
func updateTask(task Task) error {
	_, err := DB.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

/*
// респонд в случае ошибки
func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, TaskResponse{Error: message}, code)
}

// общий формат респонда
func respondWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

id := r.URL.Query().Get("id")
    if id != "" {
        // запрос на получение задачи по идентификатору
        task, err := getTaskById(id)
        if err != nil {
            http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(task)
    } else {
        // запрос на поиск задач по заголовку или комментарию
        search := r.URL.Query().Get("search")
        limit := r.URL.Query().Get("limit")
        tasks, err := searchTasks(search, limit)
        if err != nil {
            http.Error(w, `{"error": "Ошибка при поиске задач"}`, http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(tasks)
    }

*/
