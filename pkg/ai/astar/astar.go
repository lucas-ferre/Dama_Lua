package astar

import (
	"container/heap"
	"fmt"
	"math"

	"damas-go/pkg/ai/evaluation"
	"damas-go/pkg/game"
)

type Node struct {
	Board     *game.Board
	FirstMove game.Move
	GScore    float64
	HScore    float64
	FScore    float64
	Depth     int
	Index     int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].FScore < pq[j].FScore }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Node)
	item.Index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

type AStarSolver struct {
	evaluator *evaluation.Evaluator
	rules     *game.RulesEngine
	maxNodes  int
	maxDepth  int
}

func NewAStarSolver(maxNodes, maxDepth int) *AStarSolver {
	if maxNodes <= 0 {
		maxNodes = 600
	}
	if maxDepth <= 0 {
		maxDepth = 4
	}
	return &AStarSolver{
		evaluator: evaluation.NewEvaluator(),
		rules:     game.NewRulesEngine(),
		maxNodes:  maxNodes,
		maxDepth:  maxDepth,
	}
}

type SearchStats struct {
	BestMove      game.Move
	NodesExpanded int
	MinFScore     float64
	TargetReached bool
}

func (s *AStarSolver) FindBestMove(b *game.Board, color game.Color) (game.Move, SearchStats) {
	moves := s.rules.GetLegalMoves(b, color)
	stats := SearchStats{}

	if len(moves) == 0 {
		return game.Move{}, stats
	}

	if len(moves) == 1 {
		stats.BestMove = moves[0]
		stats.NodesExpanded = 1
		stats.MinFScore = 0
		return moves[0], stats
	}

	targetScore := 1000.0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for _, m := range moves {
		nb := b.Clone()
		nb.ApplyMove(m)

		g := 10.0
		if !m.IsCapture {
			g += 5.0
		} else {
			g -= float64(m.CaptureCount()) * 15.0
		}

		eval := s.evaluator.Evaluate(nb, color)
		h := math.Max(0, targetScore-eval)
		f := g + h

		node := &Node{
			Board:     nb,
			FirstMove: m,
			GScore:    g,
			HScore:    h,
			FScore:    f,
			Depth:     1,
		}
		heap.Push(&pq, node)
	}

	bestMove := moves[0]
	bestNodeF := math.MaxFloat64
	nodesCount := 0

	for pq.Len() > 0 && nodesCount < s.maxNodes {
		curr := heap.Pop(&pq).(*Node)
		nodesCount++

		if curr.FScore < bestNodeF {
			bestNodeF = curr.FScore
			bestMove = curr.FirstMove
		}

		if curr.HScore <= 0 || curr.Depth >= s.maxDepth {
			continue
		}

		nextMoves := s.rules.GetLegalMoves(curr.Board, curr.Board.Turn)
		for _, nm := range nextMoves {
			nb := curr.Board.Clone()
			nb.ApplyMove(nm)

			stepCost := 10.0
			if curr.Board.Turn != color {
				stepCost += 5.0
			}
			newG := curr.GScore + stepCost

			eval := s.evaluator.Evaluate(nb, color)
			newH := math.Max(0, targetScore-eval)
			newF := newG + newH

			child := &Node{
				Board:     nb,
				FirstMove: curr.FirstMove,
				GScore:    newG,
				HScore:    newH,
				FScore:    newF,
				Depth:     curr.Depth + 1,
			}
			heap.Push(&pq, child)
		}
	}

	stats.BestMove = bestMove
	stats.NodesExpanded = nodesCount
	stats.MinFScore = bestNodeF
	stats.TargetReached = (bestNodeF < 500.0)

	return bestMove, stats
}

func (s *AStarSolver) FormatStats(stats SearchStats) string {
	return fmt.Sprintf("A*: Custo f(n) = %.1f | Nós Expandidos = %d", stats.MinFScore, stats.NodesExpanded)
}
