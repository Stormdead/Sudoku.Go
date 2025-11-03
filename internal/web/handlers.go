package web

import (
	"encoding/json"
	"net/http"
	"sudoku/internal/game"
	"sudoku/internal/models"
	"sync"
)

// GameStore almacena las partidas activas en memoria
type GameStore struct {
	mu    sync.RWMutex
	games map[string]*models.Game
}

var store = &GameStore{
	games: make(map[string]*models.Game),
}

// Request y Response structs para JSON

// NewGameRequest es la solicitud para crear un nuevo juego
type NewGameRequest struct {
	Difficulty string `json:"difficulty"`
}

// MoveRequest es la solicitud para hacer un movimiento
type MoveRequest struct {
	GameID string `json:"game_id"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Num    int    `json:"num"`
}

// MoveResponse es la respuesta después de un movimiento
type MoveResponse struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message"`
	GameState *models.Game `json:"game_state,omitempty"`
}

// SolveRequest es la solicitud para resolver el juego
type SolveRequest struct {
	GameID string `json:"game_id"`
}

// HandleNewGame crea un nuevo juego
func HandleNewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NewGameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validar dificultad
	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}

	// Generar sudoku
	puzzle, solution := game.GenerateSudoku(req.Difficulty)

	// Crear juego
	newGame := models.NewGame(req.Difficulty)
	newGame.Board = puzzle
	newGame.Solution = solution
	newGame.Original = puzzle // Copiar el puzzle original para saber qué celdas eran dadas

	// Guardar en store
	store.mu.Lock()
	store.games[newGame.ID] = newGame
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGame)
}

// HandleValidateMove valida un movimiento
func HandleValidateMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Obtener el juego
	store.mu.RLock()
	currentGame, exists := store.games[req.GameID]
	store.mu.RUnlock()

	if !exists {
		response := MoveResponse{
			Success: false,
			Message: "Game not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar que la celda esté vacía
	if currentGame.Board[req.Row][req.Col] != 0 {
		response := MoveResponse{
			Success: false,
			Message: "Cell is not empty",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar el movimiento - CONVERTIR BOARD A ARRAY
	boardArray := [9][9]int(currentGame.Board)
	if !game.IsValidMove(&boardArray, req.Row, req.Col, req.Num) {
		response := MoveResponse{
			Success: false,
			Message: "Invalid move",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hacer el movimiento
	store.mu.Lock()
	currentGame.Board[req.Row][req.Col] = req.Num
	currentGame.Moves++

	// Verificar si está completo
	boardArrayCheck := [9][9]int(currentGame.Board)
	if game.IsComplete(&boardArrayCheck) {
		currentGame.IsComplete = true
	}

	store.mu.Unlock()

	response := MoveResponse{
		Success:   true,
		Message:   "Move accepted",
		GameState: currentGame,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGameStatus obtiene el estado del juego
func HandleGameStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		http.Error(w, "Missing game_id", http.StatusBadRequest)
		return
	}

	store.mu.RLock()
	currentGame, exists := store.games[gameID]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

// HandleValidateGame valida si el juego está completo y correcto
func HandleValidateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GameID string `json:"game_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	store.mu.RLock()
	currentGame, exists := store.games[req.GameID]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Verificar si está completo - CONVERTIR A ARRAY
	boardArrayCheck := [9][9]int(currentGame.Board)
	if !game.IsComplete(&boardArrayCheck) {
		response := MoveResponse{
			Success: false,
			Message: "Board is not complete",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verificar si es válido - CONVERTIR A ARRAY
	boardArrayValid := [9][9]int(currentGame.Board)
	if !game.IsBoardValid(&boardArrayValid) {
		response := MoveResponse{
			Success: false,
			Message: "Board has errors",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Comparar con la solución
	boardMatches := true
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if currentGame.Board[i][j] != currentGame.Solution[i][j] {
				boardMatches = false
				break
			}
		}
		if !boardMatches {
			break
		}
	}

	if !boardMatches {
		response := MoveResponse{
			Success: false,
			Message: "Solution is incorrect",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := MoveResponse{
		Success:   true,
		Message:   "Congratulations! You solved it!",
		GameState: currentGame,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSolveGame resuelve el juego automáticamente
func HandleSolveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SolveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	store.mu.RLock()
	currentGame, exists := store.games[req.GameID]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Copiar el tablero para no modificar el original - CONVERTIR A ARRAY
	boardArray := [9][9]int(currentGame.Board)
	if !game.SolveSudoku(&boardArray) {
		response := MoveResponse{
			Success: false,
			Message: "Could not solve the puzzle",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Actualizar el juego
	store.mu.Lock()
	currentGame.Board = models.Board(boardArray)
	currentGame.IsComplete = true
	store.mu.Unlock()

	response := MoveResponse{
		Success:   true,
		Message:   "Puzzle solved!",
		GameState: currentGame,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
