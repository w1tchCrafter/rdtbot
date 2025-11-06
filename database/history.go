package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/w1tchCrafter/arrays/pkg/arrays"
)

type History struct {
	db *sql.DB
}

func NewHistory() (*History, error) {
	h := &History{}
	err := h.InitDB()
	return h, err
}

func (h *History) InitDB() error {
	db, err := sql.Open("sqlite3", "history.db")

	if err != nil {
		return err
	}

	h.db = db
	return h.CreateTables()
}

func (h *History) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS history (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"url" TEXT
	);`

	_, err := h.db.Exec(query)
	return err
}

func (h *History) Close() {
	h.db.Close()
}

func (h *History) Insert(data string) error {
	query := `INSERT INTO history (url) VALUES (?)`
	_, err := h.db.Exec(query, data)
	return err
}

func (h *History) Get() (arrays.Array[string], error) {
	result := arrays.New[string]()
	rows, err := h.db.Query("SELECT url from history")

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var url string
		err = rows.Scan(&url)

		if err != nil {
			return result, err
		}

		result.Push(url)
	}

	return result, rows.Err()
}
