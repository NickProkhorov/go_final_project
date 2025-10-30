package server

import (
	"fmt"
	"net/http"
)

// Run запускает статический файловый сервер
func Run(webDir, port string) error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	addr := ":" + port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
