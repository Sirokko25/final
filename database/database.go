package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

)

type DB struct { // 1. Структура DB хранит в себе коннект к базе данных
	conn *sql.DB
}

func Createdatabase() (DB, error) { // 2. Функция CreateDB создает экземпляр структуры DB
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		log.Fatal(err)
		return DB{nil}, err
	}
	if install {
		createTableSql := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`
		_, err = db.Exec(createTableSql)
		if err != nil {
			log.Fatal(err)
			return DB{nil}, err
		}
		return DB{db}, nil
	}
	return DB{db}, nil
}

func (db *DB) Addtasktodb(task task.Task) (int64, error) { // 3. Функция AddTaskToDB становится методом структуры DB, и получает коннект из неё
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	return id, nil
}