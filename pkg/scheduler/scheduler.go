package scheduler

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Scheduler struct {
	ID      int
	Date    int
	Title   string
	Comment string
	repeat  string
}

// --- чтобы оперировать Scheduler, нужна всегда ссылка на БД
// это непонятно, и как назавать непонятно
type SchedulerList struct {
	db *sql.DB
}

func (s *SchedulerList) getItems() {
	// s.db.Open
}

func checkFileExists(dbFile string) bool {
	log.Printf("Check file existance %s", dbFile)

	_, err := os.Stat(dbFile)
	// fmt.Println(dbFile) // /tmp/go-build1725643628/b001/exe/scheduler.db

	// если файла нет, будет такая строка: 1: stat /tmp/go-build2135713667/b001/exe/scheduler.db: no such file or directory
	// если файл бд есть, то err = nil и мы попадаем во второй else
	// fmt.Println("1:", err)

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("DB file %s doesn't exist.", dbFile)
			return false
		} else {
			log.Fatal(err)
			return false
		}
	} else {
		log.Printf("DB file %s exists.", dbFile)
		return true
	}
}

func dbCreate(dbFilePath string) {
	// формируем строку для дальнейшего создания таблицы
	scheduler := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date    CHAR(8) NOT NULL DEFAULT "",
		title   VARCHAR(128) NOT NULL DEFAULT "",
		comment VARCHAR(250),
		repeat  VARCHAR(128)
	);
	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
	`

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// в методе Exec я отправляю базе данных строку запроса scheduler на выполнение
	_, err = db.Exec(scheduler)
	if err != nil {
		log.Fatal(err)
	}
}

// dbConnection проверяет существование БД. Если ёё нет, создает БД.
// как лучше назвать функцию
func DbConnection() {
	appPath, err := os.Executable()
	// fmt.Println(appPath) // /tmp/go-build1725643628/b001/exe/main

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// Get the TODO_DBFILE environment variable
	pathDb := os.Getenv("TODO_DBFILE")
	if pathDb != "" {
		dbFile = pathDb
	}

	if err != nil {
		log.Fatal(err)
	}

	if !checkFileExists(dbFile) {
		dbCreate(dbFile)
	}
}
