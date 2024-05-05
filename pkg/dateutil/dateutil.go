package dateutil

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adettelle/go_final_project/pkg/parser"
)

// NextDate calculates the next task date using the specified repeat rule
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
	var result time.Time

	log.Printf("rule[0] before switch is: %v", rule[0])
	switch {
	case rule[0] == "y":
		result, err = parser.ParseYrule(now, dt, rule)
		if err != nil {
			return "", err // чтобы тест не падал на строке log.Fatal(err), надо вместо этого возвращать ошибку
		}
	case rule[0] == "d":
		result, err = parser.ParseDrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "w":
		result, err = parser.ParseWrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	case rule[0] == "m":
		result, err = parser.ParseMrule(now, dt, rule)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Unkown repeat identifier %s", rule[0])
	}

	return result.Format("20060102"), nil
}
