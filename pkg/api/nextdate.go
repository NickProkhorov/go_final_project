package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

// NextDate вычисляет следующую дату задачи по простым правилам repeat.
// Поддерживает форматы: d <число> и y.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	// Парсим дату начала
	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", errors.New("некорректная дата начала")
	}

	// Разделяем правило, например: "d 7"
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("некорректный формат правила")
	}

	switch parts[0] {

	case "d":
		if len(parts) < 2 {
			return "", errors.New("не указано количество дней")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > 400 {
			return "", errors.New("некорректное значение количества дней (1–400)")
		}
		// увеличиваем хотя бы один раз
		for {
			start = start.AddDate(0, 0, n)
			if start.After(now) {
				break
			}
		}
		return start.Format(DateLayout), nil

	case "y":
		// ежегодное повторение
		for {
			start = start.AddDate(1, 0, 0)
			if start.After(now) {
				break
			}
		}
		return start.Format(DateLayout), nil

	default:
		return "", errors.New("неподдерживаемое правило повторения")
	}
}
