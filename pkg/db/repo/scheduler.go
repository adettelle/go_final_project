package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adettelle/go_final_project/pkg/models"
)

// чтобы оперировать Tasks (TaskCreationRequest), нужна всегда ссылка на БД
type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) TasksRepository {
	return TasksRepository{db: db}
}

func (tr TasksRepository) AddTask(t models.Task) (int, error) {
	task, err := tr.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := task.LastInsertId()
	if err != nil {
		return 0, err
	}

	// возвращаем идентификатор последней добавленной записи
	return int(id), nil
}

// put
func (tr TasksRepository) UpdateTaskInBd(t models.Task) error {
	_, err := tr.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment,"+
		"repeat = :repeat WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID))

	if err != nil {
		return err
	}

	return nil
}

// Чтение строки по заданному id.
// Из таблицы должна вернуться только одна строка.
func (tr TasksRepository) GetTask(id int) (models.Task, error) {
	s := models.Task{}
	row := tr.db.QueryRow("SELECT id, date, title, comment, repeat from scheduler WHERE id = :id",
		sql.Named("id", id))

	// заполняем объект TaskCreationRequest данными из таблицы
	err := row.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
	if err != nil {
		return s, err
	}
	return s, nil
}

// Из таблицы должны вернуться сроки с ближайшими датами.
func (tr TasksRepository) GetAllTasks() ([]models.Task, error) {
	limitConst := 20
	today := time.Now().Format("20060102")

	rows, err := tr.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= :today "+
		"ORDER BY date LIMIT :limit",
		sql.Named("today", today),
		sql.Named("limit", limitConst))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []models.Task{}
	// заполняем объект Task данными из таблицы
	for rows.Next() { // пока есть записи
		s := models.Task{} // создлаем новый объект  Task и заполняем его данными из текущего row
		err := rows.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

// Из таблицы должна вернуться срока в соответсвии с критерием поиска search.
func (tr TasksRepository) SearchTasks(search string) ([]models.Task, error) {
	limitConst := 20
	var rows *sql.Rows

	searchDate, err := time.Parse("02.01.2006", search)
	if err != nil {
		rows, err = tr.db.Query("SELECT id, date, title, comment, repeat FROM scheduler "+
			"WHERE title LIKE :search OR comment "+
			"LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", fmt.Sprintf("%%%s%%", search)),
			sql.Named("limit", limitConst))

		if err != nil {
			return nil, err
		}
	} else {
		rows, err = tr.db.Query("SELECT id, date, title, comment, repeat FROM scheduler "+
			"WHERE date LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", searchDate.Format("20060102")),
			sql.Named("limit", limitConst))

		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	result := []models.Task{}
	// заполняем объект Task данными из таблицы
	for rows.Next() { // пока есть записи
		s := models.Task{} // создлаем новый объект  Task и заполняем его данными из текущего row
		err := rows.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}
