package main

import (
	"log"
	"os"

	"github.com/NickProkhorov/go_final_project/pkg/db"
	"github.com/NickProkhorov/go_final_project/pkg/server"
)

func main() {

	webDir := "./web"

	// Читаем порт из переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Путь к БД: по умолчанию scheduler.db, иначе берем из TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализация БД (создаст таблицу/индекс, если файла не было)
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db close error:", err)
		}
	}()

	// Запускаем сервер
	if err := server.Run(webDir, port); err != nil {
		log.Fatal(err)
	}
}
