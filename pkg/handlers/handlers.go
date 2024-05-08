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

func (a *Api) MyHandle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost: // http.MethodPost - это константа
		a.CreateTask(w, r)
	case r.Method == http.MethodGet:
		if r.URL.Query().Get("search") != "" {
			s := r.URL.Query().Get("search")
			a.SearchTasks(w, r, s)
		} else {
			a.GetAllTasks(w)
		}
	}
}

func (a *Api) GetAllTasks(w http.ResponseWriter) {
	foundTasks, err := a.repo.GetAllTasks()
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusInternalServerError) // 500
		return
	}
	//w.Write([]byte(fmt.Sprintf("%v", foundTasks)))

	result := make(map[string][]models.Task) // для тестов
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp) // Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//w.Write([]byte(fmt.Sprintf("%v", resp)))
}

func (a *Api) SearchTasks(w http.ResponseWriter, r *http.Request, search string) {
	// r.URL.Query().Get(search)
	foundTasks, err := a.repo.SearchTasks(search)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task) // для тестов
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp) // Write(resp)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		// http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}
}

// CreateTask posts task into DB
func (a *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Println("err4:", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest) // 400
		return
	}
	fmt.Println("received:", buf.String())

	var parseBody models.Task
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

/*
// Чтение строки по заданному id.
// Из таблицы должна вернуться только одна строка.
func (tr TasksRepository) GetTask(id int) (models.TaskCreationRequest, error) {
	s := models.TaskCreationRequest{}
	row := tr.db.QueryRow("SELECT id, date, title, comment, repeat from task WHERE id = :id",
		sql.Named("id", id))

	// заполняем объект TaskCreationRequest данными из таблицы
	err := row.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
	if err != nil {
		return s, err
	}
	return s, nil
}
*/
