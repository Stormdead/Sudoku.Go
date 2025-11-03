package game

//Resuelve el sudoku usando backtracking
func SolveSudoku(board *[9][9]int) bool {
	row, col, found := FindEmpty(board)

	//Si no hay espacios vacíos, el sudoku está resuelto
	if !found {
		return true
	}

	//Intentar números del 1 al 9
	for num := 1; num <= 9; num++ {
		if IsValidMove(board, row, col, num) {
			board[row][col] = num

			//Recursion
			if SolveSudoku(board) {
				return true
			}

			//backtrack
			board[row][col] = 0
		}
	}
	return false
}

//Encuentra un espacio vacío en el tablero
func FindEmpty(board *[9][9]int) (int, int, bool) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				return i, j, true
			}
		}
	}
	return -1, -1, false
}

//Copia el tablero completo
func CopyBoard(source *[9][9]int) [9][9]int {
	var dest [9][9]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			dest[i][j] = source[i][j]
		}
	}
	return dest
}

//Limpia todas las celdas
func ClearBoard(board *[9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			board[i][j] = 0
		}
	}
}
