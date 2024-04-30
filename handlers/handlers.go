package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adettelle/go_final_project/pkg/scheduler"
)

// checkRepeatRule checks if the reapeat rule starts with wright letter
func checkRepeatRule(repeat string) bool {
	result := strings.Split(repeat, " ")
	match, err := regexp.MatchString("[d, m, w, y]", result[0])
	if err != nil {
		return false
	} else if result[0] == "d" {
		if len(result) == 1 {
			return false
		}
	}
	return match
}

func GetNextDay(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	//check date
	_, err := strconv.Atoi(date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = time.Parse("20060102", date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//check now
	_, err = strconv.Atoi(now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	dtNow, err := time.Parse("20060102", now)
	if err != nil {
		log.Printf("Wrong date: %v\n", err)
	}

	// check repeat
	if repeat == "" {
		err := fmt.Errorf("Empty repeat value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if !checkRepeatRule(repeat) {
		err := fmt.Errorf("Wrong repeat value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	nextDay, err := scheduler.NextDate(dtNow, date, repeat)

	if err != nil {
		err := fmt.Errorf("Wrong repeat value")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDay))
}
