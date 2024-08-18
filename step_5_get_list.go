package main

import (
	"database/sql"

	//"encoding/json"
	"net/http"
	"strconv"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

// определяет структуру для списка задач
type MultiTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error,omitempty"`
}

// по ТЗ лимит должен быть от 10 до 50, определяем любое число
const TAX_LIMIT = "15"

// Обработчик GET-запроса для получения списка задач
func handleGetList(w http.ResponseWriter, r *http.Request) {
	//извлекаем из реквеста параметры search, limit
	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")

	//по ТЗ лимит должен быть от 10 до 50
	//если лимит задан в запросе, то проверяем на цифры и на соответствие диапазону из ТЗ
	//var limit int
	//var err error
	if limitStr == "" {
		limitStr = TAX_LIMIT
	} else {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if limit < 10 || limit > 50 {
			limitStr = TAX_LIMIT
		}
	}

	//передаем в функцию поисковый запрос и лимит
	tasks, err := getTasksFromDB(search, limitStr)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, MultiTasksResponse{Tasks: tasks}, http.StatusOK)
}

// getTasksFromDB возвращает список задач из БД, удовлетворяющих условию search или ошибку
func getTasksFromDB(search string, limit string) ([]Task, error) {

	//подготовка формы
	query := ""
	var rows *sql.Rows
	var err error
	//отдельные запросы для разных ситуаций
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			search = t.Format("20060102")
			rows, err = DB.Query(query, search, limit)
			if err != nil {
				return nil, err
			}
		} else {
			//поисковую подстроку ищем везде
			search = "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, search, search, limit)
			if err != nil {
				return nil, err
			}
		}
	} else {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
		if err != nil {
			return nil, err
		}
	}

	//обрабатываем полученные записи
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = make([]Task, 0)
	}

	return tasks, nil
}

// isNumeric проверяет, что в переданной строке только цифры
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

/*func getLimit(limitStr string) string {

//по ТЗ лимит должен быть от 10 до 50
	//если лимит задан в запросе, то проверяем на цифры и на соответствие диапазону из ТЗ
	var limit int
	var err error
	if limitStr == "" {
		limitStr = TAX_LIMIT
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if limit < 10 || limit > 50 {
			limitStr = TAX_LIMIT
		}
	} */
