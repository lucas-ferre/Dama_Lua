package game

import "strings"

type Move struct {
	From      Position
	To        Position
	Path      []Position
	Captures  []Position
	IsCapture bool
	Promotion bool
}

func NewSimpleMove(from, to Position) Move {
	return Move{
		From:      from,
		To:        to,
		Path:      []Position{from, to},
		Captures:  nil,
		IsCapture: false,
		Promotion: false,
	}
}

func NewCaptureMove(from, to Position, path []Position, captures []Position) Move {
	return Move{
		From:      from,
		To:        to,
		Path:      path,
		Captures:  captures,
		IsCapture: true,
		Promotion: false,
	}
}

func (m Move) CaptureCount() int {
	return len(m.Captures)
}

func (m Move) String(size int) string {
	if len(m.Path) == 0 {
		return m.From.ToAlgebraic(size) + " -> " + m.To.ToAlgebraic(size)
	}

	parts := make([]string, len(m.Path))
	for i, p := range m.Path {
		parts[i] = p.ToAlgebraic(size)
	}

	sep := " -> "
	if m.IsCapture {
		sep = " x "
	}
	return strings.Join(parts, sep)
}

func (m Move) Equals(other Move) bool {
	if m.From != other.From || m.To != other.To || m.IsCapture != other.IsCapture {
		return false
	}
	if len(m.Captures) != len(other.Captures) {
		return false
	}
	for i := range m.Captures {
		if m.Captures[i] != other.Captures[i] {
			return false
		}
	}
	return true
}
