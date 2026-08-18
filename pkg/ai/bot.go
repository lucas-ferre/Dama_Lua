package ai

import (
	"fmt"

	"damas-go/pkg/ai/astar"
	"damas-go/pkg/ai/hillclimbing"
	"damas-go/pkg/ai/mdp"
	"damas-go/pkg/game"
)

type BotType int

const (
	BotTypeMDP BotType = iota + 1
	BotTypeAStar
	BotTypeHillClimbing
	BotTypeHybrid
)

func (bt BotType) String() string {
	switch bt {
	case BotTypeMDP:
		return "Processo de Decisao de Markov (MDP)"
	case BotTypeAStar:
		return "Busca A* Tática"
	case BotTypeHillClimbing:
		return "Hill Climbing (Random Restarts)"
	case BotTypeHybrid:
		return "Modo Híbrido Mestre (MDP + A* + HC)"
	default:
		return "Desconhecido"
	}
}

type Bot interface {
	Name() string
	Type() BotType
	SelectMove(b *game.Board) (game.Move, string)
}

type MDPBot struct {
	solver *mdp.MDPSolver
}

func NewMDPBot(depth int, gamma float64) *MDPBot {
	return &MDPBot{
		solver: mdp.NewMDPSolver(depth, gamma),
	}
}

func (b *MDPBot) Name() string {
	return "MDP-Markov"
}

func (b *MDPBot) Type() BotType {
	return BotTypeMDP
}

func (b *MDPBot) SelectMove(board *game.Board) (game.Move, string) {
	move, stats := b.solver.FindBestMove(board, board.Turn)
	return move, b.solver.FormatStats(stats)
}

type AStarBot struct {
	solver *astar.AStarSolver
}

func NewAStarBot(maxNodes, maxDepth int) *AStarBot {
	return &AStarBot{
		solver: astar.NewAStarSolver(maxNodes, maxDepth),
	}
}

func (b *AStarBot) Name() string {
	return "A*-Tactical"
}

func (b *AStarBot) Type() BotType {
	return BotTypeAStar
}

func (b *AStarBot) SelectMove(board *game.Board) (game.Move, string) {
	move, stats := b.solver.FindBestMove(board, board.Turn)
	return move, b.solver.FormatStats(stats)
}

type HillClimbingBot struct {
	solver *hillclimbing.HillClimber
}

func NewHillClimbingBot(restarts, steps int) *HillClimbingBot {
	return &HillClimbingBot{
		solver: hillclimbing.NewHillClimber(restarts, steps),
	}
}

func (b *HillClimbingBot) Name() string {
	return "HillClimber"
}

func (b *HillClimbingBot) Type() BotType {
	return BotTypeHillClimbing
}

func (b *HillClimbingBot) SelectMove(board *game.Board) (game.Move, string) {
	move, stats := b.solver.FindBestMove(board, board.Turn)
	return move, b.solver.FormatStats(stats)
}

type HybridBot struct {
	mdpBot *MDPBot
	astBot *AStarBot
	hcBot  *HillClimbingBot
	rules  *game.RulesEngine
}

func NewHybridBot() *HybridBot {
	return &HybridBot{
		mdpBot: NewMDPBot(3, 0.90),
		astBot: NewAStarBot(700, 4),
		hcBot:  NewHillClimbingBot(25, 20),
		rules:  game.NewRulesEngine(),
	}
}

func (b *HybridBot) Name() string {
	return "Mestre-Híbrido"
}

func (b *HybridBot) Type() BotType {
	return BotTypeHybrid
}

func (b *HybridBot) SelectMove(board *game.Board) (game.Move, string) {
	legalMoves := b.rules.GetLegalMoves(board, board.Turn)
	if len(legalMoves) == 0 {
		return game.Move{}, "Sem movimentos"
	}
	if len(legalMoves) == 1 {
		return legalMoves[0], "Jogada unica obrigatoria"
	}

	for _, lm := range legalMoves {
		if lm.IsCapture {
			m, stats := b.astBot.solver.FindBestMove(board, board.Turn)
			return m, fmt.Sprintf("Híbrido: A* Tático | %s", b.astBot.solver.FormatStats(stats))
		}
	}

	mMDP, sMDP := b.mdpBot.solver.FindBestMove(board, board.Turn)
	mHC, sHC := b.hcBot.solver.FindBestMove(board, board.Turn)

	if mMDP.Equals(mHC) {
		return mMDP, fmt.Sprintf("Híbrido: Consenso MDP & HC (Utilidade: %.1f)", sMDP.ExpectedUtility)
	}

	return mMDP, fmt.Sprintf("Híbrido: Decisao MDP (Util: %.1f, HC: %.1f)", sMDP.ExpectedUtility, sHC.BestEvaluation)
}

func CreateBot(botType BotType, difficulty int) Bot {
	depth := 3
	nodes := 600
	restarts := 20
	steps := 15

	if difficulty == 1 {
		depth = 2
		nodes = 250
		restarts = 10
		steps = 8
	} else if difficulty == 3 {
		depth = 4
		nodes = 1200
		restarts = 40
		steps = 25
	}

	switch botType {
	case BotTypeMDP:
		return NewMDPBot(depth, 0.90)
	case BotTypeAStar:
		return NewAStarBot(nodes, depth)
	case BotTypeHillClimbing:
		return NewHillClimbingBot(restarts, steps)
	default:
		return NewHybridBot()
	}
}
