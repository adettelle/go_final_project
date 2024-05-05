package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adettelle/go_final_project/pkg/dateutil"
	"github.com/adettelle/go_final_project/pkg/db/repo"
	"github.com/adettelle/go_final_project/pkg/models"
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

	fmt.Println("Before nextDay")
	nextDay, err := dateutil.NextDate(dtNow, date, repeat)
	fmt.Println("Error in oops:", err)

	if err != nil {
		err := fmt.Errorf("Wrong repeat value")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDay))
}

// это совокупность хэндлеров, часто называется api
type Api struct {
	repo *repo.TasksRepository
}

// это конструктор объекта api.
func NewApi(repo *repo.TasksRepository) *Api {
	return &Api{repo: repo} // создаем ссылку на объект api со свойством repo, равным repo из параметров функции
}

// CreateTask ???????????
func (a *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Println("err4:", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}
	fmt.Println("received:", buf.String())

	var parseBody models.TaskCreationRequest
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		fmt.Println("err3:", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	fmt.Println("parseBody:", parseBody)
	if err != nil {
		fmt.Println("err2:", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		fmt.Println(parseBody)
		fmt.Println(parseBody.Comment)
		return
	}

	id, err := a.repo.AddTask(parseBody)
	if err != nil {
		fmt.Println("err1:", err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", id))) //
}
