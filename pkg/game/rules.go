package game

type RulesEngine struct{}

func NewRulesEngine() *RulesEngine {
	return &RulesEngine{}
}

func (re *RulesEngine) GetLegalMoves(b *Board, color Color) []Move {
	var allCaptures []Move
	maxCaptures := 0

	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			pos := Pos(r, c)
			p := b.Get(pos)
			if p.IsEmpty() || p.Color != color {
				continue
			}

			captures := re.getPieceCaptures(b, pos, p)
			for _, capMove := range captures {
				count := capMove.CaptureCount()
				if count > maxCaptures {
					maxCaptures = count
					allCaptures = []Move{capMove}
				} else if count == maxCaptures && count > 0 {
					allCaptures = append(allCaptures, capMove)
				}
			}
		}
	}

	if len(allCaptures) > 0 {
		return allCaptures
	}

	var simpleMoves []Move
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			pos := Pos(r, c)
			p := b.Get(pos)
			if p.IsEmpty() || p.Color != color {
				continue
			}

			moves := re.getPieceSimpleMoves(b, pos, p)
			simpleMoves = append(simpleMoves, moves...)
		}
	}

	return simpleMoves
}

func (re *RulesEngine) getPieceSimpleMoves(b *Board, pos Position, p Piece) []Move {
	var moves []Move

	if p.IsMan() {
		dir := -1
		if p.Color == Black {
			dir = 1
		}

		targets := []Position{
			Pos(pos.Row+dir, pos.Col-1),
			Pos(pos.Row+dir, pos.Col+1),
		}

		for _, target := range targets {
			if target.IsValid(b.Size) && b.Get(target).IsEmpty() {
				moves = append(moves, NewSimpleMove(pos, target))
			}
		}
	} else if p.IsKing() {
		directions := [][2]int{
			{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
		}

		for _, d := range directions {
			step := 1
			for {
				target := Pos(pos.Row+d[0]*step, pos.Col+d[1]*step)
				if !target.IsValid(b.Size) {
					break
				}
				if !b.Get(target).IsEmpty() {
					break
				}
				moves = append(moves, NewSimpleMove(pos, target))
				step++
			}
		}
	}

	return moves
}

func (re *RulesEngine) getPieceCaptures(b *Board, pos Position, p Piece) []Move {
	var results []Move
	currentPath := []Position{pos}
	var currentCaptures []Position
	visitedPositions := make(map[Position]bool)

	re.exploreCaptures(b, pos, p, currentPath, currentCaptures, visitedPositions, &results)
	return results
}

func (re *RulesEngine) exploreCaptures(
	b *Board,
	currentPos Position,
	p Piece,
	path []Position,
	captures []Position,
	capturedSoFar map[Position]bool,
	results *[]Move,
) {
	foundSubCapture := false
	directions := [][2]int{
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	if p.IsMan() {
		for _, d := range directions {
			jumpOver := Pos(currentPos.Row+d[0], currentPos.Col+d[1])
			landPos := Pos(currentPos.Row+d[0]*2, currentPos.Col+d[1]*2)

			if !landPos.IsValid(b.Size) {
				continue
			}

			overPiece := b.Get(jumpOver)
			if overPiece.IsEmpty() || overPiece.Color != p.Color.Opponent() {
				continue
			}

			if capturedSoFar[jumpOver] {
				continue
			}

			landPiece := b.Get(landPos)
			if !landPiece.IsEmpty() && landPos != path[0] {
				continue
			}

			foundSubCapture = true
			newPath := make([]Position, len(path), len(path)+1)
			copy(newPath, path)
			newPath = append(newPath, landPos)

			newCaptures := make([]Position, len(captures), len(captures)+1)
			copy(newCaptures, captures)
			newCaptures = append(newCaptures, jumpOver)

			newCapturedSoFar := make(map[Position]bool)
			for k, v := range capturedSoFar {
				newCapturedSoFar[k] = v
			}
			newCapturedSoFar[jumpOver] = true

			re.exploreCaptures(b, landPos, p, newPath, newCaptures, newCapturedSoFar, results)
		}
	} else if p.IsKing() {
		for _, d := range directions {
			step := 1
			var enemyPos *Position

			for {
				checkPos := Pos(currentPos.Row+d[0]*step, currentPos.Col+d[1]*step)
				if !checkPos.IsValid(b.Size) {
					break
				}

				checkPiece := b.Get(checkPos)
				if checkPiece.IsEmpty() {
					step++
					continue
				}

				if checkPiece.Color == p.Color {
					break
				}

				if capturedSoFar[checkPos] {
					break
				}

				if checkPiece.Color == p.Color.Opponent() {
					ep := checkPos
					enemyPos = &ep
					break
				}
				step++
			}

			if enemyPos != nil {
				landStep := 1
				for {
					landPos := Pos(enemyPos.Row+d[0]*landStep, enemyPos.Col+d[1]*landStep)
					if !landPos.IsValid(b.Size) {
						break
					}

					landPiece := b.Get(landPos)
					if !landPiece.IsEmpty() && landPos != path[0] {
						break
					}

					foundSubCapture = true
					newPath := make([]Position, len(path), len(path)+1)
					copy(newPath, path)
					newPath = append(newPath, landPos)

					newCaptures := make([]Position, len(captures), len(captures)+1)
					copy(newCaptures, captures)
					newCaptures = append(newCaptures, *enemyPos)

					newCapturedSoFar := make(map[Position]bool)
					for k, v := range capturedSoFar {
						newCapturedSoFar[k] = v
					}
					newCapturedSoFar[*enemyPos] = true

					re.exploreCaptures(b, landPos, p, newPath, newCaptures, newCapturedSoFar, results)
					landStep++
				}
			}
		}
	}

	if !foundSubCapture && len(captures) > 0 {
		move := NewCaptureMove(path[0], currentPos, path, captures)
		*results = append(*results, move)
	}
}

func (re *RulesEngine) IsGameOver(b *Board) (bool, Color) {
	if b.WhiteCount == 0 {
		return true, Black
	}
	if b.BlackCount == 0 {
		return true, White
	}

	moves := re.GetLegalMoves(b, b.Turn)
	if len(moves) == 0 {
		return true, b.Turn.Opponent()
	}

	if b.HalfMoveClock >= 40 {
		return true, None
	}

	return false, None
}
