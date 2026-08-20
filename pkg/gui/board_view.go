package gui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"damas-go/pkg/game"
)

type BoardView struct {
	OffsetX     float32
	OffsetY     float32
	BoardPixels float32
	size        int
	squareSize  float32
	playerColor game.Color
}

func NewBoardView(offsetX, offsetY, boardPixels float32, size int, playerColor game.Color) *BoardView {
	return &BoardView{
		OffsetX:     offsetX,
		OffsetY:     offsetY,
		BoardPixels: boardPixels,
		size:        size,
		squareSize:  boardPixels / float32(size),
		playerColor: playerColor,
	}
}

func (bv *BoardView) PosAtScreen(mx, my int) *game.Position {
	fx := float32(mx) - bv.OffsetX
	fy := float32(my) - bv.OffsetY

	if fx < 0 || fy < 0 || fx >= bv.BoardPixels || fy >= bv.BoardPixels {
		return nil
	}

	colDisplay := int(fx / bv.squareSize)
	rowDisplay := int(fy / bv.squareSize)

	if colDisplay < 0 || colDisplay >= bv.size || rowDisplay < 0 || rowDisplay >= bv.size {
		return nil
	}

	gridRow := rowDisplay
	gridCol := colDisplay
	if bv.playerColor == game.Black {
		gridRow = bv.size - 1 - rowDisplay
		gridCol = bv.size - 1 - colDisplay
	}

	pos := game.Pos(gridRow, gridCol)
	return &pos
}

func (bv *BoardView) Draw(
	dst *ebiten.Image,
	b *game.Board,
	selectedPos *game.Position,
	validMovesFromSelected []game.Move,
	mandatoryPieces []game.Position,
	lastMove *game.Move,
) {
	DrawRect(dst, bv.OffsetX-18, bv.OffsetY-18, bv.BoardPixels+36, bv.BoardPixels+36, ColorPanelBorder)
	DrawRect(dst, bv.OffsetX-4, bv.OffsetY-4, bv.BoardPixels+8, bv.BoardPixels+8, ColorBgDark)

	for rDisp := 0; rDisp < bv.size; rDisp++ {
		for cDisp := 0; cDisp < bv.size; cDisp++ {
			gridRow := rDisp
			gridCol := cDisp
			if bv.playerColor == game.Black {
				gridRow = bv.size - 1 - rDisp
				gridCol = bv.size - 1 - cDisp
			}

			pos := game.Pos(gridRow, gridCol)
			sqX := bv.OffsetX + float32(cDisp)*bv.squareSize
			sqY := bv.OffsetY + float32(rDisp)*bv.squareSize

			isDark := pos.IsPlayable(bv.size)
			sqColor := ColorLightSquare
			if isDark {
				sqColor = ColorDarkSquare
			}

			DrawRect(dst, sqX, sqY, bv.squareSize, bv.squareSize, sqColor)

			if lastMove != nil && (lastMove.From == pos || lastMove.To == pos) {
				DrawRect(dst, sqX, sqY, bv.squareSize, bv.squareSize, ColorLastMove)
			}

			if selectedPos != nil && *selectedPos == pos {
				DrawRect(dst, sqX, sqY, bv.squareSize, bv.squareSize, ColorSelected)
			}

			for _, mp := range mandatoryPieces {
				if mp == pos {
					DrawRing(dst, sqX+bv.squareSize/2, sqY+bv.squareSize/2, bv.squareSize/2-3, 3, ColorMandatoryCap)
				}
			}

			piece := b.Get(pos)
			if !piece.IsEmpty() {
				bv.drawPiece(dst, piece, sqX, sqY)
			}

			for _, vm := range validMovesFromSelected {
				if vm.To == pos {
					cx := sqX + bv.squareSize/2
					cy := sqY + bv.squareSize/2
					dotR := bv.squareSize * 0.18
					if piece.IsEmpty() {
						DrawCircle(dst, cx, cy, dotR, ColorValidMove)
					} else {
						DrawRing(dst, cx, cy, bv.squareSize/2-4, 4, ColorValidMove)
					}
				}
			}
		}
	}

	bv.drawCoordinates(dst)
}

func (bv *BoardView) drawPiece(dst *ebiten.Image, p game.Piece, sqX, sqY float32) {
	cx := sqX + bv.squareSize/2
	cy := sqY + bv.squareSize/2
	r := bv.squareSize * 0.40

	DrawCircle(dst, cx+1.5, cy+2.0, r, color.RGBA{0, 0, 0, 120})

	pieceClr := ColorWhitePiece
	rimClr := ColorWhitePieceRim
	crownClr := ColorGoldCrown

	if p.IsBlack() {
		pieceClr = ColorBlackPiece
		rimClr = ColorBlackPieceRim
		crownClr = ColorSilverCrown
	}

	DrawCircle(dst, cx, cy, r, rimClr)
	DrawCircle(dst, cx, cy, r*0.90, pieceClr)
	DrawRing(dst, cx, cy, r*0.65, 1.5, rimClr)

	if p.IsKing() {
		DrawStar(dst, cx, cy, r*0.48, crownClr)
	}
}

func (bv *BoardView) drawCoordinates(dst *ebiten.Image) {
	for cDisp := 0; cDisp < bv.size; cDisp++ {
		colIdx := cDisp
		if bv.playerColor == game.Black {
			colIdx = bv.size - 1 - cDisp
		}
		letter := string(rune('A' + colIdx))
		x := int(bv.OffsetX + float32(cDisp)*bv.squareSize + bv.squareSize/2 - 3)
		DrawText(dst, letter, x, int(bv.OffsetY-6), ColorTextDim)
		DrawText(dst, letter, x, int(bv.OffsetY+bv.BoardPixels+14), ColorTextDim)
	}

	for rDisp := 0; rDisp < bv.size; rDisp++ {
		rank := bv.size - rDisp
		if bv.playerColor == game.Black {
			rank = rDisp + 1
		}
		rankStr := fmt.Sprintf("%d", rank)
		y := int(bv.OffsetY + float32(rDisp)*bv.squareSize + bv.squareSize/2 + 4)
		DrawText(dst, rankStr, int(bv.OffsetX-15), y, ColorTextDim)
		DrawText(dst, rankStr, int(bv.OffsetX+bv.BoardPixels+6), y, ColorTextDim)
	}
}
