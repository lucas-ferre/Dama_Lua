package mdp

import (
	"fmt"
	"math"

	"damas-go/pkg/ai/evaluation"
	"damas-go/pkg/game"
)

type MDPSolver struct {
	evaluator *evaluation.Evaluator
	rules     *game.RulesEngine
	gamma     float64
	maxDepth  int
	temp      float64
}

func NewMDPSolver(maxDepth int, gamma float64) *MDPSolver {
	if maxDepth <= 0 {
		maxDepth = 3
	}
	if gamma <= 0.0 || gamma > 1.0 {
		gamma = 0.90
	}
	return &MDPSolver{
		evaluator: evaluation.NewEvaluator(),
		rules:     game.NewRulesEngine(),
		gamma:     gamma,
		maxDepth:  maxDepth,
		temp:      50.0,
	}
}

type DecisionStats struct {
	BestMove        game.Move
	ExpectedUtility float64
	StatesEvaluated int
	ActionValues    map[string]float64
}

func (s *MDPSolver) FindBestMove(b *game.Board, color game.Color) (game.Move, DecisionStats) {
	moves := s.rules.GetLegalMoves(b, color)
	stats := DecisionStats{
		ActionValues: make(map[string]float64),
	}

	if len(moves) == 0 {
		return game.Move{}, stats
	}

	if len(moves) == 1 {
		stats.BestMove = moves[0]
		stats.ExpectedUtility = s.evaluator.Evaluate(b, color)
		stats.ActionValues[moves[0].String(b.Size)] = stats.ExpectedUtility
		return moves[0], stats
	}

	bestUtility := -math.MaxFloat64
	var bestMove game.Move
	statesCount := 0

	for _, m := range moves {
		nextBoard := b.Clone()
		nextBoard.ApplyMove(m)

		reward := s.calculateImmediateReward(b, nextBoard, m, color)
		futureVal, count := s.evaluateState(nextBoard, color, 1)
		statesCount += count

		qVal := reward + s.gamma*futureVal
		moveStr := m.String(b.Size)
		stats.ActionValues[moveStr] = qVal

		if qVal > bestUtility {
			bestUtility = qVal
			bestMove = m
		}
	}

	stats.BestMove = bestMove
	stats.ExpectedUtility = bestUtility
	stats.StatesEvaluated = statesCount

	return bestMove, stats
}

func (s *MDPSolver) calculateImmediateReward(before, after *game.Board, m game.Move, color game.Color) float64 {
	reward := 0.0
	if m.IsCapture {
		reward += float64(m.CaptureCount()) * 120.0
	}
	if m.Promotion {
		reward += 200.0
	}

	diff := s.evaluator.Evaluate(after, color) - s.evaluator.Evaluate(before, color)
	reward += diff * 0.5

	return reward
}

func (s *MDPSolver) evaluateState(b *game.Board, color game.Color, depth int) (float64, int) {
	if depth >= s.maxDepth {
		return s.evaluator.Evaluate(b, color), 1
	}

	over, winner := s.rules.IsGameOver(b)
	if over {
		if winner == color {
			return 10000.0, 1
		} else if winner == color.Opponent() {
			return -10000.0, 1
		}
		return 0.0, 1
	}

	isMyTurn := (b.Turn == color)
	moves := s.rules.GetLegalMoves(b, b.Turn)
	if len(moves) == 0 {
		if isMyTurn {
			return -10000.0, 1
		}
		return 10000.0, 1
	}

	statesCount := 1

	if isMyTurn {
		maxVal := -math.MaxFloat64
		for _, m := range moves {
			nb := b.Clone()
			nb.ApplyMove(m)
			reward := s.calculateImmediateReward(b, nb, m, color)
			futureVal, count := s.evaluateState(nb, color, depth+1)
			statesCount += count
			val := reward + s.gamma*futureVal
			if val > maxVal {
				maxVal = val
			}
		}
		return maxVal, statesCount
	}

	utilities := make([]float64, len(moves))
	maxUtil := -math.MaxFloat64
	for i, m := range moves {
		nb := b.Clone()
		nb.ApplyMove(m)
		val := s.evaluator.Evaluate(nb, color.Opponent())
		utilities[i] = val
		if val > maxUtil {
			maxUtil = val
		}
	}

	expSum := 0.0
	expVals := make([]float64, len(moves))
	for i, u := range utilities {
		ev := math.Exp((u - maxUtil) / s.temp)
		expVals[i] = ev
		expSum += ev
	}

	expectedVal := 0.0
	for i, m := range moves {
		prob := expVals[i] / expSum
		nb := b.Clone()
		nb.ApplyMove(m)
		futureVal, count := s.evaluateState(nb, color, depth+1)
		statesCount += count
		expectedVal += prob * futureVal
	}

	return expectedVal, statesCount
}

func (s *MDPSolver) FormatStats(stats DecisionStats) string {
	return fmt.Sprintf("MDP: Utilidade Esperada = %.1f | Estados = %d", stats.ExpectedUtility, stats.StatesEvaluated)
}
