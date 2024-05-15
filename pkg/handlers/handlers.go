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

const (
	MarshallingError    = "Error in marshalling JSON."
	UnMarshallingError  = "Error in unmarshalling JSON."
	ResponseWriteError  = "Error in writing data."
	ReadingError        = "Error in reading data."
	InvalidIdError      = "Invalid id."
	IdMissingError      = "ID is missing."
	InvalidDateError    = "Invalid date."
	InvalidNowDateError = "Invalid now date."
	InvalidRepeatError  = "Invalid repeat value."
	InternalServerError = "Internal server error."
	ValidatingDateError = "Error in validating date."
)

// Обработка ошибок для возврата ошибки в виде json.
type apiError struct {
	Error string `json:"error"`
}

func NewApiError(err error) apiError {
	return apiError{Error: err.Error()}
}

func (e apiError) ToJson() ([]byte, error) {
	res, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RenderApiError(w http.ResponseWriter, err error, status int) {
	apiErr := NewApiError(err)
	errorJson, _ := apiErr.ToJson()
	http.Error(w, string(errorJson), status)
}

// checkRepeatRule checks if the reapeat rule starts with correct letter
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

// GetNextDay find next day for the task
// почему мы можем сделать Get("now")?????????????? ведь у таска нет now
func GetNextDay(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now") // ????????????????
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	//check date
	_, err := strconv.Atoi(date)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidDateError), http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", date)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidDateError), http.StatusBadRequest)
		return
	}

	//check now
	_, err = strconv.Atoi(now)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidNowDateError), http.StatusBadRequest)
		return
	}

	dtNow, err := time.Parse("20060102", now)
	if err != nil {
		err := fmt.Errorf("Wrong date: %v\n", err)
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidNowDateError), http.StatusBadRequest)
		return
	}

	// check repeat
	if repeat == "" {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidRepeatError), http.StatusBadRequest)
		return
	} else if !checkRepeatRule(repeat) {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidRepeatError), http.StatusBadRequest)
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
		// idToSearch, err := AsInt(r.URL.Query().Get("id"))
		if idToSearch != "" {
			id, err := strconv.Atoi(idToSearch)
			if err != nil {
				log.Println("error:", err)
				RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
				return // иначе пойдем в a.GetTask(w, r, id) это стиль с гардами (защитниками). иначе надо написать else {a.GetTask(w, r, id)}
			}
			a.GetTask(w, r, id)
		} else {
			RenderApiError(w, fmt.Errorf(IdMissingError), http.StatusBadRequest)
			return
		}

	case r.Method == http.MethodPost:
		log.Println("We are in MethodPost")
		a.CreateTask(w, r)
	case r.Method == http.MethodPut:
		a.UpdateTask(w, r)
	case r.Method == http.MethodDelete:
		idToSearch := r.URL.Query().Get("id") // это параметр запроса
		if idToSearch != "" {
			a.DeleteTask(w, r)
		} else {
			RenderApiError(w, fmt.Errorf(IdMissingError), http.StatusBadRequest)
			return
		}
	}
}

// http://localhost:7540/api/tasks/257
func (a *Api) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	idToSearch := chi.URLParam(r, "id")
	if idToSearch != "" {
		id, err := strconv.Atoi(idToSearch)
		if err != nil {
			log.Println("error:", err)
			RenderApiError(w, fmt.Errorf(IdMissingError), http.StatusBadRequest)
			return
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
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task) // для тестов
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(ResponseWriteError), http.StatusBadRequest)
		return
	}
}

func (a *Api) SearchTasks(w http.ResponseWriter, r *http.Request, search string) {
	foundTasks, err := a.repo.SearchTasks(search)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task) // для тестов
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp) //
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(ResponseWriteError), http.StatusBadRequest)
		return
	}
}

// CreateTask posts task into DB
func (a *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf("Error in creating task"), http.StatusBadRequest) // 400
		return
	}
	log.Println("received:", buf.String())

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(UnMarshallingError), http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(ValidatingDateError), http.StatusBadRequest)
		return
	}

	id, err := a.repo.AddTask(parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError)
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
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(ReadingError), http.StatusBadRequest) // 400
		return
	}

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiError(w, fmt.Errorf(UnMarshallingError), http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(ValidatingDateError), http.StatusBadRequest)
		return
	}
	idToSearch, err := strconv.Atoi(parseBody.ID)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	_, err = a.repo.GetTask(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	err = a.repo.UpdateTaskInBd(parseBody)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError)
		return
	}

	jsonItem, err := json.Marshal(parseBody)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonItem)
}

func (a *Api) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in DeleteTask")
	idToSearch := r.URL.Query().Get("id")
	if idToSearch != "" {
		id, err := strconv.Atoi(idToSearch)
		if err != nil {
			log.Println("error:", err)
			RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
			return
		}

		err = a.repo.DeleteTask(id)
		if err != nil {
			log.Println("error:", err)
			RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusInternalServerError) // 500
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}")) // пустой JSON
			return
		}

	} else {
		RenderApiError(w, fmt.Errorf(IdMissingError), http.StatusInternalServerError) // 500
		return
	}
}

func (a *Api) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in TaskDoneHandler")
	idToSearch := r.URL.Query().Get("id")
	if idToSearch != "" {
		id, err := strconv.Atoi(idToSearch)
		if err != nil {
			log.Println("error:", err)
			RenderApiError(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
			return
		}

		newTask, err := a.repo.PostTaskDone(id)
		if newTask == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}")) // строка с пустым json
			return
		}
		if err != nil {
			log.Println("error:", err)
			RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}")) // w.Write(resp)
			return
		}
	}
}

func (a *Api) GetTask(w http.ResponseWriter, r *http.Request, id int) {
	foundTask, err := a.repo.GetTask(id)
	log.Println("we are in GetTask", "foundTask:", foundTask)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	resp, err := json.Marshal(foundTask)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp) // Write(resp)
	if err != nil {
		log.Println("error:", err)
		RenderApiError(w, fmt.Errorf(ResponseWriteError), http.StatusBadRequest)
		return
	}
}
