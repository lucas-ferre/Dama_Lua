package hillclimbing

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"damas-go/pkg/ai/evaluation"
	"damas-go/pkg/game"
)

type HillClimber struct {
	evaluator    *evaluation.Evaluator
	rules        *game.RulesEngine
	maxRestarts  int
	maxSteps     int
	randomSource *rand.Rand
}

func NewHillClimber(restarts, steps int) *HillClimber {
	if restarts <= 0 {
		restarts = 20
	}
	if steps <= 0 {
		steps = 15
	}
	return &HillClimber{
		evaluator:    evaluation.NewEvaluator(),
		rules:        game.NewRulesEngine(),
		maxRestarts:  restarts,
		maxSteps:     steps,
		randomSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type ClimbStats struct {
	BestMove        game.Move
	BestEvaluation  float64
	RestartsDone    int
	TotalIterations int
}

func (hc *HillClimber) FindBestMove(b *game.Board, color game.Color) (game.Move, ClimbStats) {
	moves := hc.rules.GetLegalMoves(b, color)
	stats := ClimbStats{}

	if len(moves) == 0 {
		return game.Move{}, stats
	}

	if len(moves) == 1 {
		stats.BestMove = moves[0]
		stats.BestEvaluation = hc.evaluator.Evaluate(b, color)
		stats.RestartsDone = 1
		stats.TotalIterations = 1
		return moves[0], stats
	}

	globalBestMove := moves[0]
	globalBestScore := -math.MaxFloat64
	totalIters := 0

	for restart := 0; restart < hc.maxRestarts; restart++ {
		startIdx := hc.randomSource.Intn(len(moves))
		currentMove := moves[startIdx]
		currentScore := hc.evaluateMovePlan(b, currentMove, color)

		for step := 0; step < hc.maxSteps; step++ {
			totalIters++
			improved := false
			neighbors := hc.getNeighborMoves(moves, currentMove)

			for _, nbMove := range neighbors {
				nbScore := hc.evaluateMovePlan(b, nbMove, color)
				if nbScore > currentScore {
					currentScore = nbScore
					currentMove = nbMove
					improved = true
					break
				}
			}

			if !improved {
				break
			}
		}

		if currentScore > globalBestScore {
			globalBestScore = currentScore
			globalBestMove = currentMove
		}
	}

	stats.BestMove = globalBestMove
	stats.BestEvaluation = globalBestScore
	stats.RestartsDone = hc.maxRestarts
	stats.TotalIterations = totalIters

	return globalBestMove, stats
}

func (hc *HillClimber) evaluateMovePlan(b *game.Board, m game.Move, color game.Color) float64 {
	nb := b.Clone()
	nb.ApplyMove(m)

	score := hc.evaluator.Evaluate(nb, color)
	if m.IsCapture {
		score += float64(m.CaptureCount()) * 80.0
	}
	if m.Promotion {
		score += 150.0
	}

	oppMoves := hc.rules.GetLegalMoves(nb, color.Opponent())
	if len(oppMoves) > 0 {
		worstOppResponse := math.MaxFloat64
		for _, om := range oppMoves {
			onb := nb.Clone()
			onb.ApplyMove(om)
			eval := hc.evaluator.Evaluate(onb, color)
			if eval < worstOppResponse {
				worstOppResponse = eval
			}
		}
		score = score*0.4 + worstOppResponse*0.6
	}

	return score
}

func (hc *HillClimber) getNeighborMoves(allMoves []MoveWrapper, current game.Move) []game.Move {
	var neighbors []game.Move
	for _, m := range allMoves {
		if !m.Equals(current) {
			neighbors = append(neighbors, m)
		}
	}

	hc.randomSource.Shuffle(len(neighbors), func(i, j int) {
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	})

	return neighbors
}

type MoveWrapper = game.Move

func (hc *HillClimber) FormatStats(stats ClimbStats) string {
	return fmt.Sprintf("Hill-Climbing: Score = %.1f | Restarts = %d | Iters = %d", stats.BestEvaluation, stats.RestartsDone, stats.TotalIterations)
}
