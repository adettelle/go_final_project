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
	"github.com/go-chi/chi/v5"
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

	log.Println("Before nextDay")
	nextDay, err := dateutil.NextDate(dtNow, date, repeat)

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

func (a *Api) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		idToSearch := r.URL.Query().Get("id") // это параметр запроса
		if idToSearch != "" {
			id, err := strconv.Atoi(idToSearch)
			if err != nil {
				log.Println("error", err)
				errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
				http.Error(w, errorJson, http.StatusBadRequest)
				return // иначе пойдем в a.GetTask(w, r, id) это стиль с гардами (защитниками). иначе надо написать else {a.GetTask(w, r, id)}
			}
			a.GetTask(w, r, id)
		} else {
			errorJson := fmt.Sprintf("{\"error\":\"%s\"}", "ID should not be empty.")
			http.Error(w, errorJson, http.StatusBadRequest)
		}

	case r.Method == http.MethodPost:
		a.CreateTask(w, r)
	case r.Method == http.MethodPut:
		a.UpdateTask(w, r)
	}
}

// http://localhost:7540/api/tasks/257
func (a *Api) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	idToSearch := chi.URLParam(r, "id")
	if idToSearch != "" {
		id, err := strconv.Atoi(idToSearch)
		if err != nil {
			log.Println("error", err)
			errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
			http.Error(w, errorJson, http.StatusBadRequest)
			return // иначе пойдем в a.GetTask(w, r, id) это стиль с гардами (защитниками). иначе надо написать else {a.GetTask(w, r, id)}
		}
		a.GetTask(w, r, id)
	}
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("search") != "" {
		s := r.URL.Query().Get("search")
		a.SearchTasks(w, r, s)
	} else {
		a.GetAllTasks(w)
	}
}

func (a *Api) GetAllTasks(w http.ResponseWriter) {
	foundTasks, err := a.repo.GetAllTasks()
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task) // для тестов
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp) // Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *Api) SearchTasks(w http.ResponseWriter, r *http.Request, search string) {
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
	_, err = w.Write(resp) //
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}
}

// CreateTask posts task into DB
func (a *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest) // 400
		return
	}
	log.Println("received:", buf.String())

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	id, err := a.repo.AddTask(parseBody)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", id))) //
}

// UpdateTask updates task in DB
func (a *Api) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest) // 400
		return
	}

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}
	idToSearch, err := strconv.Atoi(parseBody.ID)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", "Can not convert id") // ("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	_, err = a.repo.GetTask(idToSearch)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", "Can not convert id") // ("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}

	err = a.repo.UpdateTaskInBd(parseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonItem, err := json.Marshal(parseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonItem) //
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

func (a *Api) GetTask(w http.ResponseWriter, r *http.Request, id int) {
	foundTask, err := a.repo.GetTask(id)
	log.Println("we are in GetTask", "foundTask:", foundTask)
	if err != nil {
		errorJson := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		http.Error(w, errorJson, http.StatusInternalServerError) // 500
		return
	}

	resp, err := json.Marshal(foundTask)
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
		http.Error(w, errorJson, http.StatusBadRequest)
		return
	}
}
