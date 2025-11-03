package tests

import (
	"sudoku/internal/game"
	"testing"
)

// Prueba la validación de movimientos
func TestIsValidMove(t *testing.T) {
	board := [9][9]int{}

	// Tablero vacío, debe ser válido
	if !game.IsValidMove(&board, 0, 0, 5) {
		t.Error("Expected true, got false for empty cell")
	}

	// Agregar un número y verificar que no se pueda repetir en la misma fila
	board[0][0] = 5
	if game.IsValidMove(&board, 0, 1, 5) {
		t.Error("Expected false, number already in row")
	}

	// Número en diferente fila pero misma columna debe ser inválido
	if game.IsValidMove(&board, 1, 0, 5) {
		t.Error("Expected false, number already in column")
	}

	// Número en diferente fila y columna fuera de la misma caja 3x3 debe ser válido
	if !game.IsValidMove(&board, 5, 5, 5) {
		t.Error("Expected true for different row, column and box")
	}

	// Agregar un 3 en [1][1] (misma caja que [0][0])
	board[1][1] = 3
	// Intentar agregar 5 en [5][5] debe ser válido (diferente caja)
	if !game.IsValidMove(&board, 5, 5, 3) {
		t.Error("Expected true, 3 is not in box [5][5]")
	}

	// Ahora agregar un 5 en [1][2] (misma caja 3x3 que [0][0])
	board[1][2] = 5
	// Intentar agregar 5 en [2][1] debe ser inválido (misma caja)
	if game.IsValidMove(&board, 2, 1, 5) {
		t.Error("Expected false, 5 already in this 3x3 box")
	}
}

// Prueba si el tablero está completo
func TestIsComplete(t *testing.T) {
	board := [9][9]int{}

	// Tablero vacío no está completo
	if game.IsComplete(&board) {
		t.Error("Expected false for empty board")
	}

	// Llenar el tablero
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			board[i][j] = ((i*9 + j) % 9) + 1
		}
	}

	// Ahora debe estar completo
	if !game.IsComplete(&board) {
		t.Error("Expected true for full board")
	}
}

// Prueba la generación de sudokus
func TestGenerateSudoku(t *testing.T) {
	puzzle, solution := game.GenerateSudoku("medium")

	// Verificar que la solución es válida
	if !game.IsBoardValid(&solution) {
		t.Error("Solution is not valid")
	}

	// Verificar que el puzzle no está vacío
	nonEmpty := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if puzzle[i][j] != 0 {
				nonEmpty++
			}
		}
	}

	if nonEmpty == 0 {
		t.Error("Puzzle is completely empty")
	}
}

// Prueba el resolvedor
func TestSolveSudoku(t *testing.T) {
	// Crear un sudoku simple (casi completo)
	board := [9][9]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	if !game.SolveSudoku(&board) {
		t.Error("Failed to solve sudoku")
	}

	// Verificar que está completo
	if !game.IsComplete(&board) {
		t.Error("Solved board is not complete")
	}

	// Verificar que es válido
	if !game.IsBoardValid(&board) {
		t.Error("Solved board is not valid")
	}
}
