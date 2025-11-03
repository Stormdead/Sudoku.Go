package main

import (
	"fmt"
	"log"
	"net/http"
	"sudoku/internal/web"
)

func main() {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rutas web
	http.HandleFunc("/", homeHandler)

	// Rutas API
	http.HandleFunc("/api/game/new", web.HandleNewGame)
	http.HandleFunc("/api/game/move", web.HandleValidateMove)
	http.HandleFunc("/api/game/status", web.HandleGameStatus)
	http.HandleFunc("/api/game/validate", web.HandleValidateGame)
	http.HandleFunc("/api/game/solve", web.HandleSolveGame)

	fmt.Println("Servidor de Sudoku corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}
