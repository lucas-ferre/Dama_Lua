package evaluation

import (
	"math"

	"damas-go/pkg/game"
)

const (
	ManValue  = 100.0
	KingValue = 350.0
)

type Evaluator struct {
	rules *game.RulesEngine
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		rules: game.NewRulesEngine(),
	}
}

func (e *Evaluator) Evaluate(b *game.Board, color game.Color) float64 {
	over, winner := e.rules.IsGameOver(b)
	if over {
		if winner == color {
			return 100000.0
		} else if winner == color.Opponent() {
			return -100000.0
		}
		return 0.0
	}

	score := 0.0
	opponent := color.Opponent()

	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			pos := game.Pos(r, c)
			p := b.Get(pos)
			if p.IsEmpty() {
				continue
			}

			val := e.evaluatePiece(b, p, pos)
			if p.Color == color {
				score += val
			} else {
				score -= val
			}
		}
	}

	myMoves := len(e.rules.GetLegalMoves(b, color))
	oppMoves := len(e.rules.GetLegalMoves(b, opponent))
	score += float64(myMoves-oppMoves) * 5.0

	return score
}

func (e *Evaluator) evaluatePiece(b *game.Board, p game.Piece, pos game.Position) float64 {
	val := ManValue
	if p.IsKing() {
		val = KingValue
	}

	centerMin := b.Size/2 - 1
	centerMax := b.Size / 2
	distToCenter := math.Abs(float64(pos.Row)-float64(b.Size)/2.0+0.5) + math.Abs(float64(pos.Col)-float64(b.Size)/2.0+0.5)
	val += (float64(b.Size) - distToCenter) * 3.0

	if pos.Row >= centerMin && pos.Row <= centerMax && pos.Col >= centerMin && pos.Col <= centerMax {
		val += 15.0
	}

	if p.IsMan() {
		if p.Color == game.White {
			progress := float64(b.Size - 1 - pos.Row)
			val += progress * 8.0
			if pos.Row == b.Size-1 {
				val += 20.0
			}
		} else if p.Color == game.Black {
			progress := float64(pos.Row)
			val += progress * 8.0
			if pos.Row == 0 {
				val += 20.0
			}
		}
	} else if p.IsKing() {
		if pos.Row == pos.Col || pos.Row+pos.Col == b.Size-1 {
			val += 25.0
		}
	}

	return val
}
