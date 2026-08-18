package tests

import (
	"testing"

	"damas-go/pkg/game"
)

func TestRulesSimpleMovesInitialBoard(t *testing.T) {
	b := game.NewBoard(10)
	rules := game.NewRulesEngine()

	whiteMoves := rules.GetLegalMoves(b, game.White)
	if len(whiteMoves) == 0 {
		t.Fatalf("deveria haver movimentos legais para as brancas no inicio")
	}

	for _, m := range whiteMoves {
		if m.IsCapture {
			t.Fatalf("nao deveria haver capturas no tabuleiro inicial")
		}
	}
}

func TestRulesSimpleCaptureAndPromotion(t *testing.T) {
	b := game.NewBoard(10)
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			b.Set(game.Pos(r, c), game.NewPiece(game.None, game.Empty))
		}
	}
	b.RecalculateCounts()

	whitePiecePos := game.Pos(2, 1)
	blackPiecePos := game.Pos(1, 2)
	b.Set(whitePiecePos, game.NewPiece(game.White, game.Man))
	b.Set(blackPiecePos, game.NewPiece(game.Black, game.Man))
	b.RecalculateCounts()

	rules := game.NewRulesEngine()
	moves := rules.GetLegalMoves(b, game.White)

	if len(moves) != 1 {
		t.Fatalf("esperava exatamente 1 movimento de captura, obtido %d", len(moves))
	}

	m := moves[0]
	if !m.IsCapture {
		t.Fatalf("movimento deveria ser de captura")
	}
	if m.To != game.Pos(0, 3) {
		t.Fatalf("destino de captura deveria ser (0, 3)")
	}

	b.ApplyMove(m)
	if b.Get(game.Pos(0, 3)).Type != game.King {
		t.Fatalf("peca branca deveria ter sido promovida a Dama na linha 0")
	}
	if b.BlackCount != 0 {
		t.Fatalf("peca preta capturada deveria ter sido removida do tabuleiro")
	}
}

func TestRulesFlyingKingMovesAndCaptures(t *testing.T) {
	b := game.NewBoard(10)
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			b.Set(game.Pos(r, c), game.NewPiece(game.None, game.Empty))
		}
	}
	b.RecalculateCounts()

	kingPos := game.Pos(5, 5)
	enemyPos := game.Pos(3, 3)
	b.Set(kingPos, game.NewPiece(game.White, game.King))
	b.Set(enemyPos, game.NewPiece(game.Black, game.Man))
	b.RecalculateCounts()

	rules := game.NewRulesEngine()
	moves := rules.GetLegalMoves(b, game.White)

	if len(moves) == 0 {
		t.Fatalf("dama deveria ter movimentos de captura")
	}

	for _, m := range moves {
		if !m.IsCapture {
			t.Fatalf("captura de dama deve ser obrigatoria quando disponivel")
		}
	}
}
