package main

import (
	"log"
	"os"

	"github.com/NickProkhorov/go_final_project/pkg/server"
)

func main() {
	// Каталог с фронтендом
	webDir := "./web"

	// Читаем порт из переменной окружения TODO_PORT (опционально)
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Запускаем сервер
	if err := server.Run(webDir, port); err != nil {
		log.Fatal(err)
	}
}
