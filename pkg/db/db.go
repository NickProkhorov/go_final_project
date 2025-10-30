package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная ссылка на открытую БД
var DB *sql.DB

// SQL-схема: таблица scheduler и индекс по date.
// date: CHAR(8) — формат YYYYMMDD ("20060102" в Go).
// title/repeat: VARCHAR; comment: TEXT.
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// открываем/создаём БД.
func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return err
		}
	}

	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	DB = d

	if install {
		if _, err := DB.Exec(schema); err != nil {
			_ = DB.Close()
			return err
		}
	}
	return nil
}

// закрываем соединение.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
