package entities

// константы для обработки задач в бд, нужны для бизнес-логики и не зависят от внешнего окружения
const (
	DATE_FORMAT = "20060102" //установка предпочтительного формата для даты
	TASKS_LIMIT = "15"       //установка максимального количества возвращаемых записей
)

type (
	// Task структура для объекта базы данных
	Task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	// структура для ответа с ошибкой
	TaskResponse struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	}

	//структара ответа с полями задачи и ошибкой
	IDTaskResponse struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
		Error   string `json:"error,omitempty"`
	}

	// определяет структуру для списка задач
	MultiTasksResponse struct {
		Tasks []Task `json:"tasks"`
		Error string `json:"error,omitempty"`
	}
)
