package scheduler

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// daysForWrepeat gets slice of numbers of days in week for reapeat rule with letter "w"
func daysForWrepeat(rule []string) ([]int, error) {
	// var week: days of week in rule repeat
	week := []int{}
	x := strings.Split(rule[1], ",")
	for i := 0; i < len(x); i++ {
		num, err := strconv.Atoi(x[i])
		if err != nil || num > 7 || num < 1 {
			return nil, fmt.Errorf("Can not parse days for repeat value %s\n", x[i])
		}
		week = append(week, num)
	}
	return week, nil
}

// ParseYrule returns result in time.Time type when repeat rule starts with "y"
func ParseYrule(now time.Time, date time.Time, rule []string) (time.Time, error) {
	i := 1
	if rule[0] == "y" {
		for {
			result := date.AddDate(i, 0, 0)
			if result.After(now) {
				return result, nil
			}
			i++
		}
	}
	// an empty time.Time struct literal will return Go's zero date
	return time.Time{}, fmt.Errorf("expected 'y', got '%s'", rule[0])
}

// ParseDrule returns result in time.Time type when repeat rule starts with "d"
func ParseDrule(now time.Time, date time.Time, rule []string) (time.Time, error) {
	if rule[0] == "d" {
		next, err := strconv.Atoi(rule[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("expected number of days less than 400, got '%s'", rule[1])
		}
		if next > 0 && next <= 400 {
			result := date
			for {
				result = result.AddDate(0, 0, next)
				if result.After(now) {
					return result, nil
				}
			}
		}
	}
	return time.Time{}, fmt.Errorf("Error in checking days in repeat rule, got '%s'", rule)
}

// ParseWrule returns result in time.Time type when repeat rule starts with "w"
func ParseWrule(now time.Time, date time.Time, rule []string) (time.Time, error) {
	startdate := startDateForMWrule(now, date)

	if rule[0] == "w" {
		week, err := daysForWrepeat(rule)
		if err != nil {
			return time.Time{}, fmt.Errorf("Error in checking days in repeat rule 'w'. Got '%s'", rule[1])
		}
		todayWeekday := startdate.Weekday()

		sort.Ints(week) // сортируем, чтобы сразу взять тот день, что больше номером, чем сегодняшний

		numDay := int(todayWeekday)
		if numDay == 7 {
			numDay = 0
		}

		for _, n := range week {
			if n > numDay {
				result := startdate.AddDate(0, 0, n-numDay)

				return result, nil
			}
		}

		increment := 7 - int(startdate.Weekday())

		result := startdate.AddDate(0, 0, increment+week[0])
		return result, nil

	}
	return time.Time{}, fmt.Errorf("Error in checking days in repeat rule, got '%s'", rule)
}

// checkMruleDays checks the the second part of repeat rule string for "m" rule.
// -1 and -2 are converted to the last day of the month and day before last.
func checkMruleDays(now time.Time, d []string) ([]int, error) {
	// Сначала всегда рассматриваем сегодняшний месяц.
	// Смотрим со всеми днями, а уж если не подходят предложенные правилом дни (все < сегодня),
	// то месяц берем следующий месяц.
	// Если в правиле 31ое число, а в рассматриваемом месяце 30 дней,
	// то проверяется, чтобы в следующем месяце был 31 день и рассматриваем уже следующий месяц
	days := []int{}

	for _, day := range d {
		num, err := strconv.Atoi(day)
		if err != nil {
			return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", day)
		}
		if num >= 1 && num <= 31 {
			days = append(days, num)
		} else if num == -1 {
			// time.Date принимает значения вне их обычных диапазонов, то есть
			// значения нормализуются во время преобразования
			// Чтобы рассчитать количество дней текущего месяца (t), смотрим на день следующего месяца
			t := Date(now.Year(), int(now.Month()+1), 0)
			days = append(days, int(t.Day()))
		} else if num == -2 {
			// time.Date принимает значения вне их обычных диапазонов, то есть
			// значения нормализуются во время преобразования
			// Чтобы рассчитать количество дней текущего месяца (t), смотрим на день следующего месяца
			t := Date(now.Year(), int(now.Month()+1), 0)
			days = append(days, int(t.Day())-1)
		} else {
			return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", day)
		}

	}
	return days, nil
}

// checkMruleMonths checks the the third part of repeat rule string for "m" rule.
func checkMruleMonths(m []string) ([]int, error) {
	months := []int{}
	for _, month := range m {
		num, err := strconv.Atoi(month)
		if err != nil || num < 1 || num > 12 {
			return nil, fmt.Errorf("Error in checking days in repeat rule 'm', got '%s'", month)
		}
		months = append(months, num)
	}
	return months, nil
}

// Date returns time type from the int types of year, month and day.
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// ruleMwithMonth gets sorted mDays and mMonths.
func ruleMwithMonth(startdate time.Time, mDays []int, mMonths []int) time.Time {
	var nextDay time.Time

	for _, month := range mMonths {
		if month == int(startdate.Month()) {
			startdate = Date(startdate.Year(), month, 01)

			// dayInMonth is number of days in the current month.
			t := Date(startdate.Year(), int(startdate.Month())+1, 0) // день до следующего месяца
			dayInMonth := t.Day()

			for _, day := range mDays {
				if day > int(startdate.Day()) && day <= dayInMonth {
					gotDay := Date(startdate.Year(), int(startdate.Month()), day)
					nextDay = gotDay
					return nextDay
				} else if day > int(startdate.Day()) && day > dayInMonth {
					startdate = Date(startdate.Year(), int(startdate.Month())+1, 01)
				}
			}
		} else if month > int(startdate.Month()) { // else сделан для того,
			// чтобы 01 число следующего месяца тоже учитывалось в поиске
			startdate = Date(startdate.Year(), month, 01)

			// dayInMonth is number of days in the current month.
			t := Date(startdate.Year(), int(startdate.Month())+1, 0) // день до следующего месяца
			dayInMonth := t.Day()

			for _, day := range mDays {
				if day >= int(startdate.Day()) && day <= dayInMonth {
					gotDay := Date(startdate.Year(), int(startdate.Month()), day)
					nextDay = gotDay
					return nextDay
				} else if day > int(startdate.Day()) && day > dayInMonth {
					startdate = Date(startdate.Year(), int(startdate.Month())+1, 01)
				}
			}
		}
	}
	return nextDay
}

// startDateForMWrule selects startdate from now and date:
// selects a later date
func startDateForMWrule(now time.Time, date time.Time) time.Time {
	if date.After(now) {
		return date
	}
	return now
}

// ParseMrule returns result in time.Time type when repeat rule starts with "m".
func ParseMrule(now time.Time, date time.Time, rule []string) (time.Time, error) {
	startdate := startDateForMWrule(now, date)
	if rule[0] == "m" {
		days := strings.Split(rule[1], ",")
		mDays, err := checkMruleDays(startdate, days)
		if err != nil {
			return time.Time{}, fmt.Errorf("Error in checking days in repeat rule 'm'. Got '%s'", rule[1])
		}

		sort.Ints(mDays)

		// ниже проверяем, что день startdate не является больше, чем последнее число из mDays
		// если же больше, то startmonth надо сделать следующим месяцем
		var nextDay time.Time

		if len(rule) == 2 {
			for _, day := range mDays {
				if day > int(startdate.Day()) {
					nextDay = startdate.AddDate(0, 0, day-int(startdate.Day()))
					if nextDay.Day() != day {
						nextDay = Date(startdate.Year(), int(startdate.Month())+1, day)
					}
					return nextDay, nil
				}
			}

			if nextDay == Date(0001, 01, 01) { // 0001-01-01 00:00:00 +0000 UTC нулевой вариант времени
				startdate = Date(int(startdate.Year()), int(startdate.Month())+1, 01)
				for _, day := range mDays {
					if day >= int(startdate.Day()) {
						nextDay = startdate.AddDate(0, 0, day-int(startdate.Day()))
						return nextDay, nil
					}
				}
			}
		}

		if len(rule) == 3 {
			months := strings.Split(rule[2], ",")
			mMonths, err := checkMruleMonths(months)
			if err != nil {
				return time.Time{}, fmt.Errorf("Error in checking months in repeat rule 'm'. Got '%s'", rule[2])
			}
			sort.Ints(mMonths)

			nextDay = ruleMwithMonth(startdate, mDays, mMonths)
			return nextDay, nil

		} else if len(rule) > 3 {
			return time.Time{}, fmt.Errorf("Error in repeat rule 'm'. Got '%s'", rule)
		}
	}
	return time.Time{}, fmt.Errorf("Error in checking days in 'm' repeat rule, got '%s'", rule)
}

// NextDate calculates the next task date using the specified repeat rule
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Expected repeat, got an empty string.\n")
	}

	_, err := strconv.Atoi(date)
	if err != nil {
		log.Printf("Wrong type of date: %v\n", err)
		return "", err
	}

	dt, err := time.Parse("20060102", date)
	if err != nil {
		log.Printf("Wrong date: %v\n", err)
	}

	rule := strings.Split(repeat, " ")
	var result time.Time

	switch {
	case rule[0] == "y":
		result, err = ParseYrule(now, dt, rule)
		if err != nil {
			return "", err // чтобы тест не падал на строке log.Fatal(err), надо вместо этого возвращать ошибку
		}
	case rule[0] == "d":
		result, err = ParseDrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "w":
		result, err = ParseWrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "m":
		result, err = ParseMrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	}

	return result.Format("20060102"), nil
}
