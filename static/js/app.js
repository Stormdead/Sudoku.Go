let currentGame = null;
let timerInterval = null;
let moveCount = 0;
let gameHistory = [];

const API_BASE = '/api';
const STORAGE_KEY = 'sudoku_history';
const MAX_HISTORY = 10;

// Inicializar cuando carga la página
document.addEventListener('DOMContentLoaded', () => {
    console.log('Sudoku.Go cargado');
    loadHistory();
    createEmptyBoard();
    displayStats();
});

// Crear tablero vacío
function createEmptyBoard() {
    const board = document.getElementById('sudokuBoard');
    board.innerHTML = '';
    
    for (let i = 0; i < 81; i++) {
        const input = document.createElement('input');
        input.type = 'text';
        input.className = 'sudoku-cell';
        input.maxLength = '1';
        input.inputMode = 'numeric';
        input.dataset.index = i;
        input.dataset.row = Math.floor(i / 9);
        input.dataset.col = i % 9;
        input.addEventListener('input', (e) => validateInput(e));
        input.addEventListener('keydown', (e) => handleKeydown(e));
        board.appendChild(input);
    }
}

async function newGame() {
    console.log('Nuevo juego');
    const difficulty = document.getElementById('difficultySelect').value;
    console.log('Dificultad:', difficulty);
    
    try {
        const response = await fetch(`${API_BASE}/game/new`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ difficulty })
        });

        console.log('Response status:', response.status);

        if (!response.ok) {
            throw new Error('Error al crear juego: ' + response.statusText);
        }

        currentGame = await response.json();
        console.log('Juego creado:', currentGame);

        // Renderizar el tablero
        renderBoard();

        moveCount = 0;
        document.getElementById('moves').textContent = '0';
        document.getElementById('message').textContent = '';
        
        startTimer();
    } catch (error) {
        console.error('Error:', error);
        showMessage('Error al crear juego: ' + error.message, 'error');
    }
}

function renderBoard() {
    if (!currentGame) return;

    const cells = document.querySelectorAll('.sudoku-cell');
    cells.forEach((cell, index) => {
        const row = Math.floor(index / 9);
        const col = index % 9;
        const value = currentGame.board[row][col];

        cell.value = value || '';
        
        if (currentGame.original[row][col] !== 0) {
            cell.readOnly = true;
            cell.className = 'sudoku-cell readonly given';
        } else {
            cell.readOnly = false;
            cell.className = 'sudoku-cell';
        }
    });
}

function validateInput(event) {
    const input = event.target;
    const value = input.value;
    
    if (value && !/^[1-9]$/.test(value)) {
        input.value = '';
        return;
    }

    if (!currentGame) return;

    moveCount++;
    document.getElementById('moves').textContent = moveCount;
}

function handleKeydown(event) {
    if (event.key === 'Enter') {
        submitMove(event.target);
    } else if (event.key === 'Delete' || event.key === 'Backspace') {
        event.target.value = '';
    } else if (event.key === 'ArrowUp' || event.key === 'ArrowDown' || event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
        handleArrowKey(event);
    }
}

function handleArrowKey(event) {
    event.preventDefault();
    const currentCell = event.target;
    const currentIndex = parseInt(currentCell.dataset.index);
    let newIndex = currentIndex;

    switch(event.key) {
        case 'ArrowUp':
            newIndex = currentIndex - 9 >= 0 ? currentIndex - 9 : currentIndex;
            break;
        case 'ArrowDown':
            newIndex = currentIndex + 9 < 81 ? currentIndex + 9 : currentIndex;
            break;
        case 'ArrowLeft':
            newIndex = currentIndex - 1 >= 0 ? currentIndex - 1 : currentIndex;
            break;
        case 'ArrowRight':
            newIndex = currentIndex + 1 < 81 ? currentIndex + 1 : currentIndex;
            break;
    }

    const cells = document.querySelectorAll('.sudoku-cell');
    cells[newIndex].focus();
}

async function submitMove(cell) {
    if (!currentGame || cell.readOnly) return;

    const row = parseInt(cell.dataset.row);
    const col = parseInt(cell.dataset.col);
    const num = parseInt(cell.value);

    if (!num || num < 1 || num > 9) {
        cell.value = '';
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/game/move`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                game_id: currentGame.id,
                row,
                col,
                num
            })
        });

        const data = await response.json();

        if (!data.success) {
            showMessage('No valido: ' + data.message, 'error');
            cell.classList.add('error');
            setTimeout(() => {
                cell.value = '';
                cell.classList.remove('error');
            }, 500);
            return;
        }

        currentGame = data.game_state;
        renderBoard();
        cell.classList.add('valid');
        setTimeout(() => {
            cell.classList.remove('valid');
        }, 300);
        showMessage('Movimiento valido', 'success');

        if (currentGame.is_complete) {
            showMessage('Tablero completo. Valida tu solucion.', 'info');
        }
    } catch (error) {
        console.error('Error:', error);
        showMessage('Error: ' + error.message, 'error');
        cell.value = '';
    }
}

async function validateGame() {
    if (!currentGame) {
        showMessage('Crea un juego primero', 'error');
        return;
    }

    console.log('Validar juego');
    
    try {
        const response = await fetch(`${API_BASE}/game/validate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                game_id: currentGame.id
            })
        });

        const data = await response.json();

        if (data.success) {
            showMessage(data.message, 'success');
            stopTimer();
            saveToHistory(true);
            displayStats();
        } else {
            showMessage(data.message, 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showMessage('Error: ' + error.message, 'error');
    }
}

async function solveGame() {
    if (!currentGame) {
        showMessage('Crea un juego primero', 'error');
        return;
    }

    console.log('Resolver juego');
    
    try {
        const response = await fetch(`${API_BASE}/game/solve`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                game_id: currentGame.id
            })
        });

        const data = await response.json();

        if (data.success) {
            currentGame = data.game_state;
            renderBoard();
            showMessage('Sudoku resuelto', 'success');
            stopTimer();
            saveToHistory(false);
            displayStats();
        } else {
            showMessage(data.message, 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showMessage('Error: ' + error.message, 'error');
    }
}

function showMessage(msg, type) {
    const messageEl = document.getElementById('message');
    messageEl.textContent = msg;
    messageEl.className = 'message ' + type;
    
    if (type === 'success') {
        setTimeout(() => {
            messageEl.textContent = '';
            messageEl.className = '';
        }, 3000);
    }
}

function startTimer() {
    if (timerInterval) clearInterval(timerInterval);
    
    let seconds = 0;
    timerInterval = setInterval(() => {
        seconds++;
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        document.getElementById('timer').textContent = 
            `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
    }, 1000);
}

function stopTimer() {
    if (timerInterval) clearInterval(timerInterval);
}

function saveToHistory(completed) {
    if (!currentGame) return;

    const timerText = document.getElementById('timer').textContent;
    const historyEntry = {
        id: currentGame.id,
        difficulty: currentGame.difficulty,
        completed: completed,
        time: timerText,
        moves: moveCount,
        timestamp: new Date().toLocaleString()
    };

    gameHistory.unshift(historyEntry);
    if (gameHistory.length > MAX_HISTORY) {
        gameHistory.pop();
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(gameHistory));
    displayHistory();
}

function loadHistory() {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
        gameHistory = JSON.parse(stored);
        displayHistory();
    }
}

function displayHistory() {
    const historyDiv = document.getElementById('history') || createHistorySection();
    
    if (gameHistory.length === 0) {
        historyDiv.innerHTML = '<p>No hay partidas guardadas aun</p>';
        return;
    }

    let html = '<div class="history-list">';
    gameHistory.forEach((entry, index) => {
        const status = entry.completed ? 'Completado' : 'Resuelto';
        html += `
            <div class="history-item">
                <div class="history-item-info">
                    <span class="difficulty-badge ${entry.difficulty}">${entry.difficulty}</span>
                    <span style="margin-left: 10px;">${status}</span>
                    <div class="history-item-time">${entry.time} - ${entry.timestamp}</div>
                </div>
                <div class="history-item-action">
                    <button onclick="deleteHistoryEntry(${index})">Eliminar</button>
                </div>
            </div>
        `;
    });
    html += '</div>';
    
    historyDiv.innerHTML = html;
}

function createHistorySection() {
    const container = document.querySelector('.container');
    const section = document.createElement('div');
    section.id = 'history';
    section.className = 'history-section';
    section.innerHTML = '<h3>Historial de Partidas</h3>';
    container.appendChild(section);
    return section;
}

function deleteHistoryEntry(index) {
    gameHistory.splice(index, 1);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(gameHistory));
    displayHistory();
    displayStats();
}

function displayStats() {
    const statsDiv = document.getElementById('stats') || createStatsSection();
    
    const completed = gameHistory.filter(g => g.completed).length;
    const total = gameHistory.length;
    const winRate = total > 0 ? Math.round((completed / total) * 100) : 0;

    statsDiv.innerHTML = `
        <div class="stats">
            <div class="stat-card">
                <div class="stat-card-value">${total}</div>
                <div class="stat-card-label">Partidas Totales</div>
            </div>
            <div class="stat-card">
                <div class="stat-card-value">${completed}</div>
                <div class="stat-card-label">Completadas</div>
            </div>
            <div class="stat-card">
                <div class="stat-card-value">${winRate}%</div>
                <div class="stat-card-label">Tasa Exito</div>
            </div>
        </div>
    `;
}

function createStatsSection() {
    const container = document.querySelector('.container');
    const section = document.createElement('div');
    section.id = 'stats';
    section.className = 'stats-section';
    container.insertBefore(section, document.querySelector('.history-section'));
    return section;
}

// Limpiar historial
function clearAllHistory() {
    if (confirm('Seguro que quieres eliminar todo el historial?')) {
        gameHistory = [];
        localStorage.removeItem(STORAGE_KEY);
        displayHistory();
        displayStats();
        showMessage('Historial eliminado', 'success');
    }
}