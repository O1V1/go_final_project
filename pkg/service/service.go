package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/O1V1/go_final_project/pkg/entities"
	"github.com/O1V1/go_final_project/pkg/storage"
)

var dateFormat = entities.DATE_FORMAT

// интерфейс TaskService определяет методы для работы с Task
type TaskService interface {
	PrepareTaskTitleAndDate(task entities.Task, now time.Time) (entities.Task, error)
	NextDate(now time.Time, date string, repeat string) (string, error)
}

// структура для работы taskService
type taskService struct {
	taskStorage storage.TaskStorage
}

// конструктор нового экземпляра структуры taskService
func NewTaskService(taskStorage storage.TaskStorage) TaskService {
	return &taskService{
		taskStorage: taskStorage,
	}
}

// метод структуры taskService, бывшая функция prepareTaskTitleAndDate
// осуществляет подготовку структуры task для дальшейшего использования
func (s *taskService) PrepareTaskTitleAndDate(task entities.Task, now time.Time) (entities.Task, error) {
	nowFormat := now.Format(dateFormat)
	//проводятся проверки различных полей task на соответствие требованиям БД
	//проверка формата поля Title, оно не должно быть пустое
	if task.Title == "" {
		return entities.Task{}, errors.New("invalid title format")
	}

	//Если поле `date` не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = nowFormat
	}

	//проверка формата поля Date
	_, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return entities.Task{}, err
	}

	//Обработка поля Date, если указана прошедшая дата
	if task.Date < nowFormat {
		//используется текущая дата, если нет правила повторов
		if task.Repeat == "" {
			task.Date = nowFormat
			//используется функция nextDate, если указано правило повторов
		} else {
			nextDate, err := s.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return entities.Task{}, err
			}
			task.Date = nextDate
		}
	}
	return task, nil
}

// метод структуры taskService, бывшая функция NextDate
// NextDate вычисляет следующую дату на основе правил повторения
func (s *taskService) NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим исходную дату
	initialDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %s", date)
	}

	// БАЗОВОЕ ПРАВИЛО 1. Если правило не указано, отмеченная выполненной задача будет удаляться из таблицы
	// по ТЗ ошибка должна возвращаться, если в колонке`repeat` — пустая строка;
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	// Определяем переменную для следующей даты
	var nextDate time.Time

	//Выделяем первый символ правила и переменную для дальн инстркуции
	indicator := rune(repeat[0])
	//inStr := ""

	//Для наглядности выделяем правила черезе switch
	switch indicator {

	//БАЗОВОЕ ПРАВИЛО 2 — задача переносится на указанное число дней
	//спросить, нужно ли учитывать лишние пробелы
	case 'd':
		inStr, found := strings.CutPrefix(repeat, string(indicator)+" ")
		if !found || inStr == "" {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

		days, err := strconv.Atoi(inStr)
		if err != nil || days > 400 || days < 1 {
			return "", fmt.Errorf("некорректное значение дней: %s", inStr)
		}

		nextDate = getCorrectDate(now, initialDate, 0, 0, days)

	//БАЗОВОЕ ПРАВИЛО 3  — задача выполняется ежегодно
	// Уточнить, нужно ли учитывать лишние пробелы в строке
	case 'y':
		if repeat == "y" {
			nextDate = getCorrectDate(now, initialDate, 1, 0, 0)
		} else {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

	//ПРАВИЛО 1 СО ЗВЕЗДОЧКОЙ
	case 'w':
		inStr, found := strings.CutPrefix(repeat, string(indicator)+" ")
		if !found || inStr == "" {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

		weekDays := strings.Split(inStr, ",")
		if len(weekDays) > 7 {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

		//проверка, что используются правильные символы
		// спросить, нужно  ли проверять на повторы?
		for _, smb := range weekDays {
			if !strings.Contains("1234567", smb) {
				return "", fmt.Errorf("неверный формат дней недели: %s", inStr)
			}
		}

		nextDate = getNextWeekDay(now, weekDays)

		// проверка, чтобы новая дата была позже now
		for !nextDate.After(initialDate) {
			nextDate = getNextWeekDay(nextDate, weekDays)
		}

	//ПРАВИЛО 2 СО ЗВЕЗДОЧКОЙ — задача назначается в указанные дни месяца. вторая последовательность опциональна и указывает на определённые месяцы
	case 'm':
		inStr, found := strings.CutPrefix(repeat, string(indicator)+" ")
		if !found || inStr == "" {
			return "", fmt.Errorf("неверный формат, нет пробела после буквы: %s", repeat)
		}
		tasks := strings.Split(inStr, " ")
		lnTasks := len(tasks)
		if lnTasks > 2 || lnTasks < 1 {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

		daysS := strings.Split(tasks[0], ",")
		if len(daysS) > 31 || len(daysS) == 0 {
			return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
		}

		//проверка, что используются правильные символы
		// спросить, нужно  ли проверять на повторы?
		for _, smb := range daysS {
			numD, err := strconv.Atoi(smb)
			if err != nil {
				return "", fmt.Errorf("нечисловые символы для дней месяца: %s", smb)
			}
			if numD > 31 || numD < -2 || numD == 0 {
				return "", fmt.Errorf("некорректные числа для дней месяца: %s", inStr)
			}
		}
		days := parseToIntSlice(daysS)

		var months []int
		if lnTasks == 2 {
			monthsStr := strings.Split(tasks[1], ",")
			if len(monthsStr) > 12 {
				return "", fmt.Errorf("неверный формат повторений месяцев: %s", repeat)
			}
			//проверка, что используются правильные символы
			// спросить, нужно  ли проверять на повторы?
			for _, smb := range monthsStr {
				numM, err := strconv.Atoi(smb)
				if err != nil {
					return "", fmt.Errorf("неверный формат номеров месяцев: %s", smb)
				}
				if numM > 12 || numM < 1 {
					return "", fmt.Errorf("некорректные числа для номеров месяцев: %s", inStr)
				}
			}
			//после проверок можно перевети в числа
			months = parseToIntSlice(monthsStr)
		}

		// передаем данные в функцию для получения новой даты
		nextDate = getNextMonthDay(now, days, months)

		// проверка, чтобы новая дата была позже now
		for !nextDate.After(initialDate) {
			nextDate = getNextMonthDay(nextDate, days, months)
		}

	// если кейсы не подошли, возвращаем ошибку
	default:
		return "", fmt.Errorf("неверный формат правила повторения: %s", repeat)
	}

	return nextDate.Format(dateFormat), nil
}

// далее - вспомогательные функции, которые используются методами этого слоя

// isNumeric проверяет, что в переданной строке только цифры
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// находит дату переноса задачи, позже текущей
func getCorrectDate(now, inDate time.Time, y, m, d int) time.Time {
	nextD := inDate.AddDate(y, m, d)
	for !nextD.After(now) {
		nextD = nextD.AddDate(y, m, d)
	}
	return nextD
}

func getNextWeekDay(now time.Time, weekDays []string) time.Time {
	weekDaysI := parseToIntSlice(weekDays)
	for {
		now = now.AddDate(0, 0, 1)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for _, day := range weekDaysI {
			if weekday == day {
				return now
			}
		}
	}
}

// getNextMonthDay находит следующий день и месяц
func getNextMonthDay(now time.Time, days []int, months []int) time.Time {
	month := int(now.Month())
	if len(months) > 0 && !isContains(months, month) {
		now = time.Date(now.Year(), now.Month(), 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
		for {
			now = now.AddDate(0, 1, 0)
			month = int(now.Month())
			if isContains(months, month) {
				break
			}
		}
		now = now.AddDate(0, 0, -1)
	}

	for {
		now = now.AddDate(0, 0, 1)
		//month = int(now.Month())
		day := now.Day()
		for _, dayX := range days {
			if dayX == -1 && now.AddDate(0, 0, 1).Day() == 1 {
				return now
			} else if dayX == -2 && now.AddDate(0, 0, 2).Day() == 1 {
				return now
			} else if dayX == day {
				return now
			}
		}
	}
}

// parseToIntSliceparseMonths парсит строковый слайс в слайс чисел
func parseToIntSlice(strSlice []string) []int {
	intSlice := make([]int, 0, len(strSlice))
	for _, item := range strSlice {
		num, _ := strconv.Atoi(item)
		intSlice = append(intSlice, num)
	}
	return intSlice
}

// contains проверяет, содержит ли список элемент
func isContains(list []int, item int) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
