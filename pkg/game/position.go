package game

import (
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	Row int
	Col int
}

func Pos(row, col int) Position {
	return Position{Row: row, Col: col}
}

func (p Position) IsValid(size int) bool {
	return p.Row >= 0 && p.Row < size && p.Col >= 0 && p.Col < size
}

func (p Position) IsPlayable(size int) bool {
	if !p.IsValid(size) {
		return false
	}
	return (p.Row+p.Col)%2 == 1
}

func (p Position) ToAlgebraic(size int) string {
	if !p.IsValid(size) {
		return "??"
	}
	colChar := string(rune('A' + p.Col))
	rowNum := size - p.Row
	return fmt.Sprintf("%s%d", colChar, rowNum)
}

func ParseAlgebraic(s string, size int) (Position, error) {
	clean := strings.TrimSpace(strings.ToUpper(s))
	if len(clean) < 2 || len(clean) > 3 {
		return Position{}, fmt.Errorf("posicao invalida: %s", s)
	}

	colChar := clean[0]
	if colChar < 'A' || colChar >= byte('A'+size) {
		return Position{}, fmt.Errorf("coluna invalida: %c", colChar)
	}
	col := int(colChar - 'A')

	rowStr := clean[1:]
	rowNum, err := strconv.Atoi(rowStr)
	if err != nil || rowNum < 1 || rowNum > size {
		return Position{}, fmt.Errorf("linha invalida: %s", rowStr)
	}
	row := size - rowNum

	pos := Pos(row, col)
	if !pos.IsValid(size) {
		return Position{}, fmt.Errorf("posicao fora dos limites: %s", s)
	}

	return pos, nil
}
