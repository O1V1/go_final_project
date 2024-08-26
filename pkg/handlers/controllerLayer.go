package handler

import (
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/O1V1/go_final_project/pkg/config"
	"github.com/O1V1/go_final_project/pkg/entities"
	"github.com/O1V1/go_final_project/pkg/service"
	"github.com/O1V1/go_final_project/pkg/storage"
)

var DATE_FORMAT = config.DATE_FORMAT

type (
	//структура для обслуживания различных задач
	TaskHandlerImpl struct {
		taskService    service.TaskService
		taskRepository storage.TaskRepository
	}
	//структура для работы с датами
	DateHandlerImpl struct {
		taskService service.TaskService
	}
)

// интерфейсы для
type TaskHandler interface {
	HandlePostTask(w http.ResponseWriter, r *http.Request)
	HandleGetList(w http.ResponseWriter, r *http.Request)
	HandleGetTask(w http.ResponseWriter, r *http.Request)
	HandleDeleteTask(w http.ResponseWriter, r *http.Request)
	HandleTaskDone(w http.ResponseWriter, r *http.Request)
	NextDateHandler(w http.ResponseWriter, r *http.Request)
}

// конструктор для структуры TaskHandlerImpl
func NewTaskHandler(taskService service.TaskService, taskRepository storage.TaskRepository) *TaskHandlerImpl {
	return &TaskHandlerImpl{taskService: taskService, taskRepository: taskRepository}
}

// конструктор для структуры DateHandlerImpl
func NewDateHandler(taskService service.TaskService) *DateHandlerImpl {
	return &DateHandlerImpl{taskService: taskService}
}

// метод структуры TaskHandlerImpl для удаления задачи
func (h *TaskHandlerImpl) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	//получение поля id из запроса
	id := r.URL.Query().Get("id")

	//удаление записи
	err := h.taskRepository.DeleteTask(id)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	//отправка пустого JSON-объекта в качестве ответа на HTTP-запрос
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// метод структуры TaskHandlerImpl добавление новой задачи в базу данных
func (h *TaskHandlerImpl) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	var task entities.Task
	//фиксирую текущую дату
	now := time.Now()

	//десериализация реквеста в структуру, обработка ошибки и выход в случае ошибки
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//проверка и заполнение полей
	task, err := h.taskService.PrepareTaskTitleAndDate(task, now)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr, err := h.taskRepository.AddTask(task)
	if err != nil {
		respondWithError(w, "Failed to add task", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, entities.TaskResponse{ID: idStr, Error: ""}, http.StatusOK)
}

// метод структуры TaskHandlerImpl для получения задачи по id
func (h *TaskHandlerImpl) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	//извлечение id из параметров запроса
	id := r.URL.Query().Get("id")
	//получение таски по запрошенному id, обработка ошибки
	task, err := h.taskRepository.GetTaskByID(id)
	if err != nil {
		//если задача не получена, возврат респонса с полем error
		response := entities.IDTaskResponse{
			Error: "Задача не найдена",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	//при успешном получении задачи, запись в ответ
	response := entities.IDTaskResponse{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	json.NewEncoder(w).Encode(response)
}

// метод структуры TaskHandlerImpl для  получения списка задач
func (h *TaskHandlerImpl) HandleGetList(w http.ResponseWriter, r *http.Request) {
	//извлекаем из реквеста параметры search, limit
	search := r.URL.Query().Get("search")
	emptyTaskSlice := make([]entities.Task, 0)
	//передаем в функцию поисковый запрос и лимит
	tasks, err := h.taskRepository.GetTasks(search)
	if err != nil {
		respondWithJSON(w, entities.MultiTasksResponse{
			Tasks: emptyTaskSlice, Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, entities.MultiTasksResponse{Tasks: tasks}, http.StatusOK)
}

// метод структуры TaskHandlerImpl для выполнения задачи
func (h *TaskHandlerImpl) HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	//получение поля id
	id := r.URL.Query().Get("id")
	//получение записи из БД по полю id
	task, err := h.taskRepository.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}
	//Одноразовая задача с пустым полем `repeat` удаляется.
	switch task.Repeat {
	case "":
		err = h.taskRepository.DeleteTask(id)
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
		nextDate, err := h.taskService.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			respondWithError(w, "Invalid repeat rule", http.StatusBadRequest)
			return
		}
		task.Date = nextDate
		//обновление измененных полей задачи
		err = h.taskRepository.UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error": "не обновлено"}`, http.StatusNotImplemented)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// метод структуры TaskHandlerImpl для редактирования задачи
func (h *TaskHandlerImpl) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entities.Task
	now := time.Now()
	//nowFormat := now.Format(DATE_FORMAT)
	// извлечение тела запроса в структуру task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//проверка формата поля ID
	if _, err = h.taskRepository.GetTaskByID(task.ID); err != nil {
		respondWithError(w, "Invalid ID formal", http.StatusBadRequest)
		return
	}

	//проверки полей Date, Title task на соответствие требованиям БД
	task, err = h.taskService.PrepareTaskTitleAndDate(task, now)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//обновление параметров задачи, обработка ошибки
	err = h.taskRepository.UpdateTask(task)
	if err != nil {
		respondWithError(w, "Task update failed", http.StatusNotImplemented)
		return
	}
	//после успешного обновления задачи формируется ответ на http запрос
	//в виде пустой структуры Task в формате JSON
	json.NewEncoder(w).Encode(entities.Task{})
}

// метод структуры DateHandlerImpl для получения даты в соответствии с указанным правилом повторений
func (h *DateHandlerImpl) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || date == "" || repeat == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Парсим дату `now`
	now, err := time.Parse(DATE_FORMAT, nowStr)
	if err != nil {
		http.Error(w, "Некорректный формат даты 'now'", http.StatusBadRequest)
		return
	}

	// Вызываем функцию NextDate для вычисления следующей даты
	nextDate, err := h.taskService.NextDate(now, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
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
	respondWithJSON(w, entities.TaskResponse{Error: message}, code)
}
