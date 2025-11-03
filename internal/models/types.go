package models

import "time"

// Board representa el tablero de 9x9
type Board [9][9]int

// Game contiene la información de una partida
type Game struct {
	ID         string    `json:"id"`
	Board      Board     `json:"board"`
	Solution   Board     `json:"solution"`
	Original   Board     `json:"original"`
	Difficulty string    `json:"difficulty"`
	StartTime  time.Time `json:"start_time"`
	Moves      int       `json:"moves"`
	IsComplete bool      `json:"is_complete"`
}

// NewGame crea una nueva partida
func NewGame(difficulty string) *Game {
	return &Game{
		ID:         generateID(),
		Board:      Board{},
		Solution:   Board{},
		Original:   Board{},
		Difficulty: difficulty,
		StartTime:  time.Now(),
		Moves:      0,
		IsComplete: false,
	}
}

// generateID genera un ID simple para la partida
func generateID() string {
	return time.Now().Format("20060102150405")
}
