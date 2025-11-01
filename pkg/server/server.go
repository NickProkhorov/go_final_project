package server

import (
	"fmt"
	"net/http"

	"github.com/NickProkhorov/go_final_project/pkg/api"
)

// Run запускает статический файловый сервер и API
func Run(webDir, port string) error {
	// создаём маршрутизатор
	mux := http.NewServeMux()

	// регистрируем обработчик API
	mux.HandleFunc("/api/nextdate", api.NextDateHandler)

	// раздаём фронтенд
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	addr := ":" + port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
