package game

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateSudoku genera un sudoku nuevo de acuerdo a la dificultad
func GenerateSudoku(difficulty string) ([9][9]int, [9][9]int) {
	//Crear una solución válida
	solution := [9][9]int{}
	FillBoard(&solution)

	// Copiar la solución para el puzzle
	puzzle := solution

	// Remover números según dificultad
	cellsToRemove := getCellsToRemove(difficulty)
	removeNumbers(&puzzle, cellsToRemove)

	return puzzle, solution
}

// Llena el tablero con una solución válida
func FillBoard(board *[9][9]int) bool {
	row, col, found := FindEmpty(board)

	if !found {
		return true // Tablero lleno
	}

	// Generar números aleatorios del 1 al 9
	numbers := rand.Perm(9)

	for _, num := range numbers {
		number := num + 1 // Convertir 0-8 a 1-9

		if IsValidMove(board, row, col, number) {
			board[row][col] = number

			if FillBoard(board) {
				return true
			}

			// Backtrack
			board[row][col] = 0
		}
	}

	return false
}

// Remueve números del tablero para crear el puzzle
func removeNumbers(board *[9][9]int, count int) {
	removed := 0

	for removed < count {
		row := rand.Intn(9)
		col := rand.Intn(9)

		if board[row][col] != 0 {
			board[row][col] = 0
			removed++
		}
	}
}

// Retorna cuántas celdas remover según dificultad
func getCellsToRemove(difficulty string) int {
	switch difficulty {
	case "easy":
		return 30 // Más números = más fácil
	case "medium":
		return 40
	case "hard":
		return 50 // Menos números = más difícil
	default:
		return 40
	}
}
