package game

const (
	DefaultBoardSize = 10
)

type Board struct {
	Size             int
	Grid             [][]Piece
	Turn             Color
	WhiteCount       int
	BlackCount       int
	WhiteKingCount   int
	BlackKingCount   int
	MoveCount        int
	HalfMoveClock    int
	LastMove         *Move
}

func NewBoard(size int) *Board {
	if size <= 0 {
		size = DefaultBoardSize
	}
	if size%2 != 0 {
		size++
	}

	b := &Board{
		Size: size,
		Grid: make([][]Piece, size),
		Turn: White,
	}

	for r := 0; r < size; r++ {
		b.Grid[r] = make([]Piece, size)
		for c := 0; c < size; c++ {
			b.Grid[r][c] = NewPiece(None, Empty)
		}
	}

	b.InitPieces()
	return b
}

func (b *Board) InitPieces() {
	pieceRows := (b.Size - 2) / 2
	b.WhiteCount = 0
	b.BlackCount = 0
	b.WhiteKingCount = 0
	b.BlackKingCount = 0

	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			pos := Pos(r, c)
			if !pos.IsPlayable(b.Size) {
				b.Grid[r][c] = NewPiece(None, Empty)
				continue
			}

			if r < pieceRows {
				b.Grid[r][c] = NewPiece(Black, Man)
				b.BlackCount++
			} else if r >= b.Size-pieceRows {
				b.Grid[r][c] = NewPiece(White, Man)
				b.WhiteCount++
			} else {
				b.Grid[r][c] = NewPiece(None, Empty)
			}
		}
	}
}

func (b *Board) Get(pos Position) Piece {
	if !pos.IsValid(b.Size) {
		return NewPiece(None, Empty)
	}
	return b.Grid[pos.Row][pos.Col]
}

func (b *Board) Set(pos Position, p Piece) {
	if !pos.IsValid(b.Size) {
		return
	}
	b.Grid[pos.Row][pos.Col] = p
}

func (b *Board) Clone() *Board {
	cb := &Board{
		Size:           b.Size,
		Grid:           make([][]Piece, b.Size),
		Turn:           b.Turn,
		WhiteCount:     b.WhiteCount,
		BlackCount:     b.BlackCount,
		WhiteKingCount: b.WhiteKingCount,
		BlackKingCount: b.BlackKingCount,
		MoveCount:      b.MoveCount,
		HalfMoveClock:  b.HalfMoveClock,
	}

	if b.LastMove != nil {
		lm := *b.LastMove
		cb.LastMove = &lm
	}

	for r := 0; r < b.Size; r++ {
		cb.Grid[r] = make([]Piece, b.Size)
		copy(cb.Grid[r], b.Grid[r])
	}

	return cb
}

func (b *Board) ApplyMove(m Move) {
	piece := b.Get(m.From)
	b.Set(m.From, NewPiece(None, Empty))

	if m.IsCapture {
		for _, capPos := range m.Captures {
			capturedPiece := b.Get(capPos)
			if capturedPiece.IsWhite() {
				b.WhiteCount--
				if capturedPiece.IsKing() {
					b.WhiteKingCount--
				}
			} else if capturedPiece.IsBlack() {
				b.BlackCount--
				if capturedPiece.IsKing() {
					b.BlackKingCount--
				}
			}
			b.Set(capPos, NewPiece(None, Empty))
		}
		b.HalfMoveClock = 0
	} else {
		b.HalfMoveClock++
	}

	shouldPromote := false
	if piece.IsMan() {
		if piece.IsWhite() && m.To.Row == 0 {
			shouldPromote = true
		} else if piece.IsBlack() && m.To.Row == b.Size-1 {
			shouldPromote = true
		}
	}

	if shouldPromote {
		piece.Type = King
		if piece.IsWhite() {
			b.WhiteKingCount++
		} else {
			b.BlackKingCount++
		}
		m.Promotion = true
	}

	b.Set(m.To, piece)

	last := m
	b.LastMove = &last
	b.Turn = b.Turn.Opponent()
	b.MoveCount++
}

func (b *Board) RecalculateCounts() {
	b.WhiteCount = 0
	b.BlackCount = 0
	b.WhiteKingCount = 0
	b.BlackKingCount = 0

	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			p := b.Grid[r][c]
			if p.IsWhite() {
				b.WhiteCount++
				if p.IsKing() {
					b.WhiteKingCount++
				}
			} else if p.IsBlack() {
				b.BlackCount++
				if p.IsKing() {
					b.BlackKingCount++
				}
			}
		}
	}
}
