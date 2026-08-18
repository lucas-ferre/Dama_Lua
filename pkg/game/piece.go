package game

type Color int

const (
	None Color = iota
	White
	Black
)

func (c Color) Opponent() Color {
	if c == White {
		return Black
	}
	if c == Black {
		return White
	}
	return None
}

func (c Color) String() string {
	switch c {
	case White:
		return "Brancas"
	case Black:
		return "Pretas"
	default:
		return "Nenhuma"
	}
}

type PieceType int

const (
	Empty PieceType = iota
	Man
	King
)

func (pt PieceType) String() string {
	switch pt {
	case Man:
		return "Pedra"
	case King:
		return "Dama"
	default:
		return "Vazia"
	}
}

type Piece struct {
	Color Color
	Type  PieceType
}

func NewPiece(color Color, pieceType PieceType) Piece {
	return Piece{
		Color: color,
		Type:  pieceType,
	}
}

func (p Piece) IsEmpty() bool {
	return p.Color == None || p.Type == Empty
}

func (p Piece) IsWhite() bool {
	return p.Color == White
}

func (p Piece) IsBlack() bool {
	return p.Color == Black
}

func (p Piece) IsKing() bool {
	return p.Type == King
}

func (p Piece) IsMan() bool {
	return p.Type == Man
}

func (p Piece) Symbol() string {
	if p.IsEmpty() {
		return " "
	}
	if p.Color == White {
		if p.Type == King {
			return "★"
		}
		return "●"
	}
	if p.Type == King {
		return "☆"
	}
	return "○"
}
