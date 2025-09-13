package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	// Парсим исходную дату
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	// Разбиваем правило на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	// Обрабатываем разные типы правил
	switch parts[0] {
	case "d":
		return handleDailyRule(now, date, parts)
	case "y":
		return handleYearlyRule(now, date)
	case "w":
		return handleWeeklyRule(now, date, parts)
	case "m":
		return handleMonthlyRule(now, date, parts)
	default:
		return "", fmt.Errorf("unsupported repeat format: %s", parts[0])
	}
}

// afterNow проверяет, что дата больше текущей (игнорируя время)
func afterNow(date, now time.Time) bool {
	// Нормализуем даты (убираем время)
	dateNorm := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowNorm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateNorm.After(nowNorm)
}

// handleDailyRule обрабатывает правило с днями
func handleDailyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid daily rule format")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid interval: %s", parts[1])
	}

	if interval <= 0 || interval > 400 {
		return "", fmt.Errorf("interval out of range (1-400): %d", interval)
	}

	// Увеличиваем дату пока она не станет больше now
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format("20060102"), nil
}

// handleYearlyRule обрабатывает ежегодное правило
func handleYearlyRule(now, date time.Time) (string, error) {
	// Увеличиваем на год пока дата не станет больше now
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format("20060102"), nil
}

// handleWeeklyRule обрабатывает правило с днями недели
func handleWeeklyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid weekly rule format")
	}

	// Парсим дни недели
	daysStr := strings.Split(parts[1], ",")
	weekdays := make(map[int]bool)

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("invalid weekday: %s", dayStr)
		}
		if day < 1 || day > 7 {
			return "", fmt.Errorf("weekday out of range (1-7): %d", day)
		}
		weekdays[day] = true
	}

	if len(weekdays) == 0 {
		return "", fmt.Errorf("no valid weekdays specified")
	}

	// Ищем ближайший подходящий день
	current := date
	for i := 0; i < 400; i++ { // Защита от бесконечного цикла
		current = current.AddDate(0, 0, 1)

		// Convert to ISO weekday (1=Monday, 7=Sunday)
		weekday := int(current.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}

		if weekdays[weekday] && afterNow(current, now) {
			return current.Format("20060102"), nil
		}
	}

	return "", fmt.Errorf("no valid date found within reasonable range")
}

// handleMonthlyRule обрабатывает правило с днями месяца
func handleMonthlyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid monthly rule format")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	monthDays := make(map[int]bool)

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("invalid month day: %s", dayStr)
		}
		if day != -1 && day != -2 && (day < 1 || day > 31) {
			return "", fmt.Errorf("month day out of range: %d", day)
		}
		monthDays[day] = true
	}

	// Парсим месяцы (если указаны)
	months := make(map[int]bool)
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				return "", fmt.Errorf("invalid month: %s", monthStr)
			}
			if month < 1 || month > 12 {
				return "", fmt.Errorf("month out of range (1-12): %d", month)
			}
			months[month] = true
		}
	}

	// Ищем ближайший подходящий день
	current := date
	for i := 0; i < 366*3; i++ { // 3 года максимум
		current = current.AddDate(0, 0, 1)

		// Проверяем месяц
		if len(months) > 0 {
			currentMonth := int(current.Month())
			if !months[currentMonth] {
				continue
			}
		}

		// Проверяем день
		currentDay := current.Day()
		lastDay := getLastDayOfMonth(current.Year(), int(current.Month()))

		validDay := false
		for day := range monthDays {
			switch day {
			case -1: // последний день месяца
				if currentDay == lastDay {
					validDay = true
				}
			case -2: // предпоследний день месяца
				if currentDay == lastDay-1 {
					validDay = true
				}
			default:
				if currentDay == day {
					validDay = true
				}
			}
			if validDay {
				break
			}
		}

		if validDay && afterNow(current, now) {
			return current.Format("20060102"), nil
		}
	}

	return "", fmt.Errorf("no valid date found within reasonable range")
}

// getLastDayOfMonth возвращает последний день месяца
func getLastDayOfMonth(year, month int) int {
	// Создаем дату на первый день следующего месяца и вычитаем 1 день
	if month == 12 {
		return 31 // Декабрь всегда имеет 31 день
	}
	firstOfNextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}
