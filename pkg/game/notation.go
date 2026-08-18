package game

import (
	"fmt"
	"regexp"
	"strings"
)

type NotationParser struct {
	size int
}

func NewNotationParser(size int) *NotationParser {
	return &NotationParser{size: size}
}

func (np *NotationParser) ParseInput(input string, legalMoves []Move) (Move, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return Move{}, fmt.Errorf("entrada vazia")
	}

	norm := strings.ToLower(trimmed)
	norm = strings.ReplaceAll(norm, "para", " ")
	norm = strings.ReplaceAll(norm, "to", " ")
	norm = strings.ReplaceAll(norm, "->", " ")
	norm = strings.ReplaceAll(norm, "-", " ")
	norm = strings.ReplaceAll(norm, "x", " ")
	norm = strings.ReplaceAll(norm, ":", " ")

	re := regexp.MustCompile(`[a-z][0-9]+`)
	matches := re.FindAllString(norm, -1)

	if len(matches) < 2 {
		return Move{}, fmt.Errorf("formato invalido. Exemplo: 'E3 para F4' ou 'C3 D4'")
	}

	positions := make([]Position, 0, len(matches))
	for _, m := range matches {
		pos, err := ParseAlgebraic(m, np.size)
		if err != nil {
			return Move{}, err
		}
		positions = append(positions, pos)
	}

	from := positions[0]
	to := positions[len(positions)-1]

	var matchedMoves []Move
	for _, lm := range legalMoves {
		if lm.From == from && lm.To == to {
			if len(positions) > 2 {
				if np.pathMatches(lm.Path, positions) {
					matchedMoves = append(matchedMoves, lm)
				}
			} else {
				matchedMoves = append(matchedMoves, lm)
			}
		}
	}

	if len(matchedMoves) == 0 {
		pieceHasLegalMove := false
		for _, lm := range legalMoves {
			if lm.From == from {
				pieceHasLegalMove = true
				break
			}
		}

		if !pieceHasLegalMove {
			hasCaptures := false
			for _, lm := range legalMoves {
				if lm.IsCapture {
					hasCaptures = true
					break
				}
			}
			if hasCaptures {
				return Move{}, fmt.Errorf("jogada ilegal: ha capturas obrigatorias a serem feitas (Lei da Maioria)")
			}
			return Move{}, fmt.Errorf("jogada ilegal: a peca em %s nao possui movimentos validos", from.ToAlgebraic(np.size))
		}

		return Move{}, fmt.Errorf("jogada invalida de %s para %s", from.ToAlgebraic(np.size), to.ToAlgebraic(np.size))
	}

	bestMove := matchedMoves[0]
	for _, m := range matchedMoves[1:] {
		if m.CaptureCount() > bestMove.CaptureCount() {
			bestMove = m
		}
	}

	return bestMove, nil
}

func (np *NotationParser) pathMatches(fullPath []Position, inputPath []Position) bool {
	if len(inputPath) > len(fullPath) {
		return false
	}
	for i, pos := range inputPath {
		if i == 0 && fullPath[0] != pos {
			return false
		}
		if i == len(inputPath)-1 && fullPath[len(fullPath)-1] != pos {
			return false
		}
	}
	return true
}

func (np *NotationParser) FormatMove(m Move) string {
	return m.String(np.size)
}
