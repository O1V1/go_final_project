package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"

	"github.com/O1V1/go_final_project/pkg/config"
	"github.com/O1V1/go_final_project/pkg/entities"
)

var DATE_FORMAT = config.DATE_FORMAT

// интерфейс определяет методы для работы с сущностью Task
type TaskRepository interface {
	GetTaskByID(id string) (entities.Task, error)
	GetTasks(search string) ([]entities.Task, error)
	AddTask(task entities.Task) (string, error)
	UpdateTask(task entities.Task) error
	DeleteTask(id string) error
}

// taskRepository является структурой, которая реализует интерфейс TaskRepository
type taskRepository struct {
	db *sql.DB
}

// конструктор для структуры taskRepository
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

// метод структуры taskRepository, бывшая функция  addTaskToDB
func (r *taskRepository) AddTask(task entities.Task) (string, error) {
	//подготовка запроса
	stmt, err := r.db.Prepare("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)")
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
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	//id получен в формате int64
	idStr := strconv.Itoa(int(id))
	return idStr, nil
}

// метод структуры taskRepository, бывшая функция getTasksFromDB
func (r *taskRepository) GetTasks(search string) ([]entities.Task, error) {
	//подготовка формы
	limit := config.TASKS_LIMIT
	query := ""
	var rows *sql.Rows
	var err error
	t, errT := time.Parse("02.01.2006", search)
	emptyTaskSlice := make([]entities.Task, 0)
	//отдельные запросы для разных ситуаций
	switch {
	//если поисковый запрос - пустая строка
	case search == "":
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = r.db.Query(query, limit)
		if err != nil {
			return emptyTaskSlice, err
		}
	//выбор задачи на конкретную дату
	case search != "" && errT == nil:
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		search = t.Format(DATE_FORMAT)
		rows, err = r.db.Query(query, search, limit)
		if err != nil {
			return emptyTaskSlice, err
		}
	//поиск подстроки в базе данных
	default:
		//поисковую подстроку ищем везде
		search = "%" + search + "%"
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		rows, err = r.db.Query(query, search, search, limit)
		if err != nil {
			return emptyTaskSlice, err
		}
	}

	//обрабатываем полученные записи
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return emptyTaskSlice, err
		}
		tasks = append(tasks, task)
	}
	//проверка наличия ошибки, возвращаемой функцией rows.Next
	if err := rows.Err(); err != nil {
		return emptyTaskSlice, err
	}

	if tasks == nil {
		tasks = emptyTaskSlice
	}

	return tasks, nil
}

// метод структуры taskRepository, бывшая функция getTasksById
// получения экземпляра Task по его id из БД
func (r *taskRepository) GetTaskByID(id string) (entities.Task, error) {
	var task entities.Task
	//если в id не только цифры, то не беспокоим БД
	if !isNumeric(id) {
		return entities.Task{}, errors.New("ID must be a number")
	}
	err := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return entities.Task{}, err
	}
	return task, nil
}

// метод структуры taskRepository, бывшая функция updateTask
// обновление экземпляра Task в БД
func (r *taskRepository) UpdateTask(task entities.Task) error {
	_, err := r.db.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	return nil
}

// метод структуры taskRepository, бывшая функция DeleteTask
// удаляет запись из БД с помощью SQL-запроса DELETE
func (r *taskRepository) DeleteTask(id string) error {
	//проверка, что id содержит только цифры (и что не пустая строка)
	if !isNumeric(id) {
		return errors.New("invalid task ID")
	}
	result, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
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

// isNumeric проверяет, что в переданной строке только цифры
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
