package tests

import (
	"testing"

	"damas-go/pkg/game"
)

func TestParseAlgebraic(t *testing.T) {
	pos, err := game.ParseAlgebraic("E3", 10)
	if err != nil {
		t.Fatalf("erro inesperado ao parsear E3: %v", err)
	}
	if pos.Row != 7 || pos.Col != 4 {
		t.Fatalf("esperado Row=7, Col=4 para E3 em 10x10, obtido Row=%d, Col=%d", pos.Row, pos.Col)
	}

	algebraic := pos.ToAlgebraic(10)
	if algebraic != "E3" {
		t.Fatalf("esperado E3, obtido %s", algebraic)
	}
}

func TestNotationParser(t *testing.T) {
	b := game.NewBoard(10)
	rules := game.NewRulesEngine()
	parser := game.NewNotationParser(10)

	legalMoves := rules.GetLegalMoves(b, game.White)
	if len(legalMoves) == 0 {
		t.Fatalf("sem movimentos legais iniciais")
	}

	sampleMove := legalMoves[0]
	fromStr := sampleMove.From.ToAlgebraic(10)
	toStr := sampleMove.To.ToAlgebraic(10)

	inputs := []string{
		fromStr + " para " + toStr,
		fromStr + " to " + toStr,
		fromStr + " " + toStr,
		fromStr + " - " + toStr,
		fromStr + " -> " + toStr,
	}

	for _, in := range inputs {
		m, err := parser.ParseInput(in, legalMoves)
		if err != nil {
			t.Fatalf("falha ao interpretar '%s': %v", in, err)
		}
		if m.From != sampleMove.From || m.To != sampleMove.To {
			t.Fatalf("movimento parseado nao corresponde ao esperado para '%s'", in)
		}
	}
}
