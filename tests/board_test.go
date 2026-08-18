package tests

import (
	"testing"

	"damas-go/pkg/game"
)

func TestBoardInitialization10x10(t *testing.T) {
	b := game.NewBoard(10)
	if b.Size != 10 {
		t.Fatalf("esperado tamanho 10, obtido %d", b.Size)
	}
	if b.WhiteCount != 20 {
		t.Fatalf("esperado 20 pecas brancas, obtido %d", b.WhiteCount)
	}
	if b.BlackCount != 20 {
		t.Fatalf("esperado 20 pecas pretas, obtido %d", b.BlackCount)
	}
	if b.Turn != game.White {
		t.Fatalf("esperado inicio com as brancas")
	}
}

func TestBoardInitialization8x8(t *testing.T) {
	b := game.NewBoard(8)
	if b.Size != 8 {
		t.Fatalf("esperado tamanho 8, obtido %d", b.Size)
	}
	if b.WhiteCount != 12 {
		t.Fatalf("esperado 12 pecas brancas, obtido %d", b.WhiteCount)
	}
	if b.BlackCount != 12 {
		t.Fatalf("esperado 12 pecas pretas, obtido %d", b.BlackCount)
	}
}

func TestBoardCloneAndApplyMove(t *testing.T) {
	b := game.NewBoard(10)
	from := game.Pos(6, 1)
	to := game.Pos(5, 0)
	move := game.NewSimpleMove(from, to)

	clone := b.Clone()
	clone.ApplyMove(move)

	if !b.Get(from).IsWhite() {
		t.Fatalf("tabuleiro original foi modificado na posicao de origem")
	}
	if !clone.Get(from).IsEmpty() {
		t.Fatalf("clone deveria estar vazio na posicao de origem")
	}
	if !clone.Get(to).IsWhite() {
		t.Fatalf("clone deveria conter a peca branca na posicao de destino")
	}
	if clone.Turn != game.Black {
		t.Fatalf("turno deveria ter passado para as Pretas")
	}
}
