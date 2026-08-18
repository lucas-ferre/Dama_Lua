package tests

import (
	"testing"

	"damas-go/pkg/ai"
	"damas-go/pkg/game"
)

func TestMDPSolverDecision(t *testing.T) {
	b := game.NewBoard(10)
	bot := ai.CreateBot(ai.BotTypeMDP, 2)

	move, info := bot.SelectMove(b)
	if move.From == move.To && len(move.Path) == 0 {
		t.Fatalf("MDP deveria retornar um movimento valido")
	}
	if info == "" {
		t.Fatalf("MDP deveria retornar informacoes de avaliacao")
	}
}

func TestAStarSolverDecision(t *testing.T) {
	b := game.NewBoard(10)
	bot := ai.CreateBot(ai.BotTypeAStar, 2)

	move, info := bot.SelectMove(b)
	if move.From == move.To && len(move.Path) == 0 {
		t.Fatalf("A* deveria retornar um movimento valido")
	}
	if info == "" {
		t.Fatalf("A* deveria retornar informacoes de avaliacao")
	}
}

func TestHillClimberDecision(t *testing.T) {
	b := game.NewBoard(10)
	bot := ai.CreateBot(ai.BotTypeHillClimbing, 2)

	move, info := bot.SelectMove(b)
	if move.From == move.To && len(move.Path) == 0 {
		t.Fatalf("Hill Climbing deveria retornar um movimento valido")
	}
	if info == "" {
		t.Fatalf("Hill Climbing deveria retornar informacoes de avaliacao")
	}
}

func TestHybridBotDecision(t *testing.T) {
	b := game.NewBoard(10)
	bot := ai.CreateBot(ai.BotTypeHybrid, 2)

	move, info := bot.SelectMove(b)
	if move.From == move.To && len(move.Path) == 0 {
		t.Fatalf("Bot Hibrido deveria retornar um movimento valido")
	}
	if info == "" {
		t.Fatalf("Bot Hibrido deveria retornar informacoes de avaliacao")
	}
}

func TestAICapturePriority(t *testing.T) {
	b := game.NewBoard(10)
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			b.Set(game.Pos(r, c), game.NewPiece(game.None, game.Empty))
		}
	}
	b.Turn = game.Black

	blackPiecePos := game.Pos(4, 3)
	whitePiecePos := game.Pos(5, 4)
	b.Set(blackPiecePos, game.NewPiece(game.Black, game.Man))
	b.Set(whitePiecePos, game.NewPiece(game.White, game.Man))
	b.RecalculateCounts()

	bot := ai.CreateBot(ai.BotTypeAStar, 2)
	move, _ := bot.SelectMove(b)

	if !move.IsCapture {
		t.Fatalf("IA deveria obrigatoriamente executar a captura disponivel")
	}
	if move.To != game.Pos(6, 5) {
		t.Fatalf("destino da captura da IA incorreto, esperado (6, 5), obtido (%d, %d)", move.To.Row, move.To.Col)
	}
}
