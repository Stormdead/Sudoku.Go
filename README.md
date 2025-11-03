# Sudoku.Go

Una aplicación web moderna de Sudoku construida con **Go** en el backend y **JavaScript/HTML/CSS** en el frontend. Aprende Go resolviendo Sudokus.

## 🎮 Características

- ✨ **Interfaz moderna y responsiva** - Diseño gradiente con animaciones suaves
- 🎯 **Tres niveles de dificultad** - Fácil, Medio y Difícil
- ⏱️ **Cronómetro integrado** - Mide tu velocidad de resolución
- 🤖 **Resolvedor automático** - Usa backtracking para resolver Sudokus
- 💾 **Historial de partidas** - Guardado en localStorage
- 📊 **Estadísticas** - Tasa de éxito, partidas totales, partidas completadas
- ⌨️ **Navegación con teclado** - Usa flechas para moverte entre celdas
- 📱 **Responsive design** - Funciona en móvil, tablet y escritorio
- 🎨 **Validación visual** - Celdas con feedback de error/éxito

## 🏗️ Arquitectura

### Backend (Go)
- `internal/game/` - Lógica del Sudoku (generación, validación, resolución)
- `internal/models/` - Estructuras de datos
- `internal/web/` - Handlers HTTP y API REST
- `main.go` - Punto de entrada

### Frontend (JavaScript/HTML/CSS)
- `templates/index.html` - Estructura HTML
- `static/js/app.js` - Lógica del juego
- `static/css/style.css` - Estilos modernos

## 🚀 Instalación

### Requisitos
- Go 1.16+
- Navegador moderno (Chrome, Firefox, Safari, Edge)

### Pasos

1. **Clonar el repositorio**
```bash
git clone https://github.com/Stormdead/Sudoku.Go.git
cd Sudoku.Go
