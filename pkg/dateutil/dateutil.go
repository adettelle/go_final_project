package dateutil

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adettelle/go_final_project/pkg/parser"
)

// NextDate calculates the next task date using the specified repeat rule.
// now — время от которого ищется ближайшая дата;
// date — исходное время в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения в одном из форматов:
// d <число> - задача переносится на указанное число дней;
// y - задача выполняется ежегодно;
// w <через запятую от 1 до 7> - задача назначается в указанные дни недели,
// где 1 — понедельник, 7 — воскресенье;
// m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12] -
// задача назначается в указанные дни месяца.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Expected repeat, got an empty string.")
	}

	_, err := strconv.Atoi(date)
	if err != nil {
		log.Printf("Wrong type of date: %v\n", err)
		return "", err
	}

	dt, err := time.Parse("20060102", date)
	if err != nil {
		log.Printf("Wrong date: %v\n", err)
		return "", err
	}

	rule := strings.Split(repeat, " ")

	var parsedRepeat parser.Repeat

	log.Printf("rule[0] before switch is: %v", rule[0])
	switch {
	case rule[0] == "y":
		parsedRepeat, err = parser.ParseYRepeat(rule)
		if err != nil {
			return "", err // чтобы тест не падал на строке log.Fatal(err), надо вместо этого возвращать ошибку
		}
	case rule[0] == "d":
		parsedRepeat, err = parser.ParseDRepeat(rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "w":
		parsedRepeat, err = parser.ParseWRepeat(rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "m":
		parsedRepeat, err = parser.ParseMRepeat(rule, now, date)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Unkown repeat identifier %s", rule[0])
	}

	d, err := parsedRepeat.GetNextDate(now, dt)
	if err != nil {
		return "", err
	}
	return d.Format("20060102"), nil
}
