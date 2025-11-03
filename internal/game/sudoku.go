package game

// Verifica si un número es válido en una posición
func IsValidMove(board *[9][9]int, row, col, num int) bool {
	// Validar rango
	if row < 0 || row > 8 || col < 0 || col > 8 || num < 1 || num > 9 {
		return false
	}

	// Si la celda ya tiene un número, no es válido
	if board[row][col] != 0 {
		return false
	}

	// Verificar fila
	for i := 0; i < 9; i++ {
		if board[row][i] == num {
			return false
		}
	}

	// Verificar columna
	for i := 0; i < 9; i++ {
		if board[i][col] == num {
			return false
		}
	}

	// Verificar caja 3x3
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3

	for i := boxRow; i < boxRow+3; i++ {
		for j := boxCol; j < boxCol+3; j++ {
			if board[i][j] == num {
				return false
			}
		}
	}

	return true
}

// Verifica si el tablero está completo
func IsComplete(board *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

// Verifica si el tablero completo es válido (sin duplicados)
func IsBoardValid(board *[9][9]int) bool {
	// Verificar filas
	for i := 0; i < 9; i++ {
		seen := make(map[int]bool)
		for j := 0; j < 9; j++ {
			num := board[i][j]
			if num != 0 {
				if seen[num] {
					return false
				}
				seen[num] = true
			}
		}
	}

	// Verificar columnas
	for j := 0; j < 9; j++ {
		seen := make(map[int]bool)
		for i := 0; i < 9; i++ {
			num := board[i][j]
			if num != 0 {
				if seen[num] {
					return false
				}
				seen[num] = true
			}
		}
	}

	// Verificar cajas 3x3
	for boxRow := 0; boxRow < 9; boxRow += 3 {
		for boxCol := 0; boxCol < 9; boxCol += 3 {
			seen := make(map[int]bool)
			for i := boxRow; i < boxRow+3; i++ {
				for j := boxCol; j < boxCol+3; j++ {
					num := board[i][j]
					if num != 0 {
						if seen[num] {
							return false
						}
						seen[num] = true
					}
				}
			}
		}
	}

	return true
}

// Verifica si hay duplicados en una fila
func CheckRow(board *[9][9]int, row int, num int) bool {
	for i := 0; i < 9; i++ {
		if board[row][i] == num {
			return false
		}
	}
	return true
}

// Verifica si hay duplicados en una columna
func CheckCol(board *[9][9]int, col int, num int) bool {
	for i := 0; i < 9; i++ {
		if board[i][col] == num {
			return false
		}
	}
	return true
}

// Verifica si hay duplicados en la caja 3x3
func CheckBox(board *[9][9]int, row, col, num int) bool {
	boxRow := (row / 3) * 3
	boxCol := (col / 3) * 3

	for i := boxRow; i < boxRow+3; i++ {
		for j := boxCol; j < boxCol+3; j++ {
			if board[i][j] == num {
				return false
			}
		}
	}
	return true
}
