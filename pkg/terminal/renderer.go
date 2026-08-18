package terminal

import (
	"fmt"
	"strings"

	"damas-go/pkg/game"
)

type Renderer struct {
	size int
}

func NewRenderer(size int) *Renderer {
	return &Renderer{size: size}
}

func (r *Renderer) RenderGame(
	b *game.Board,
	history []string,
	aiName string,
	aiEval string,
	message string,
) string {
	boardLines := r.renderBoardLines(b)
	hudLines := r.renderHUDLines(b, history, aiName, aiEval)

	maxLines := len(boardLines)
	if len(hudLines) > maxLines {
		maxLines = len(hudLines)
	}

	boardWidth := 0
	if len(boardLines) > 0 {
		boardWidth = VisibleLen(boardLines[0])
	}

	hudWidth := 0
	if len(hudLines) > 0 {
		hudWidth = VisibleLen(hudLines[0])
	}

	var sb strings.Builder
	for i := 0; i < maxLines; i++ {
		bLine := ""
		if i < len(boardLines) {
			bLine = boardLines[i]
		} else {
			bLine = strings.Repeat(" ", boardWidth)
		}

		hLine := ""
		if i < len(hudLines) {
			hLine = hudLines[i]
		} else {
			hLine = strings.Repeat(" ", hudWidth)
		}

		sb.WriteString(bLine)
		sb.WriteString("   ")
		sb.WriteString(hLine)
		sb.WriteString("\n")
	}

	if message != "" {
		totalWidth := boardWidth + 3 + hudWidth
		if totalWidth < 40 {
			totalWidth = 60
		}
		sb.WriteString(r.renderStatusBanner(message, totalWidth))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (r *Renderer) renderBoardLines(b *game.Board) []string {
	table := NewTable()
	table.BorderStyle = UnicodeDoubleBorders
	table.Padding = 0

	headers := make([]string, b.Size+2)
	headers[0] = "  "
	for c := 0; c < b.Size; c++ {
		headers[c+1] = fmt.Sprintf(" %c ", 'A'+c)
	}
	headers[b.Size+1] = "  "
	table.SetHeaders(headers...)

	for row := 0; row < b.Size; row++ {
		rank := b.Size - row
		cells := make([]string, b.Size+2)
		cells[0] = fmt.Sprintf("%2d", rank)
		cells[b.Size+1] = fmt.Sprintf("%-2d", rank)

		for col := 0; col < b.Size; col++ {
			pos := game.Pos(row, col)
			piece := b.Get(pos)
			isDark := pos.IsPlayable(b.Size)

			cellContent := r.formatPieceCell(piece, pos, isDark, b.LastMove)
			cells[col+1] = cellContent
		}

		table.AddRow(cells...)
	}

	rendered := table.Render()
	return strings.Split(strings.TrimRight(rendered, "\n"), "\n")
}

func (r *Renderer) formatPieceCell(p game.Piece, pos game.Position, isDark bool, lastMove *game.Move) string {
	bg := BgLightSquare
	if isDark {
		bg = BgDarkSquare
	}

	if lastMove != nil {
		if lastMove.From == pos || lastMove.To == pos {
			bg = BgLastMove
		}
	}

	sym := "   "
	if !p.IsEmpty() {
		switch p.Color {
		case game.White:
			if p.IsKing() {
				sym = fmt.Sprintf(" %s%s%s ", FgBrightYellow+Bold, "★", Reset+bg)
			} else {
				sym = fmt.Sprintf(" %s%s%s ", FgBrightCyan+Bold, "●", Reset+bg)
			}
		case game.Black:
			if p.IsKing() {
				sym = fmt.Sprintf(" %s%s%s ", FgBrightMagenta+Bold, "☆", Reset+bg)
			} else {
				sym = fmt.Sprintf(" %s%s%s ", FgBrightRed+Bold, "○", Reset+bg)
			}
		}
	}

	return fmt.Sprintf("%s%s%s", bg, sym, Reset)
}

func (r *Renderer) renderHUDLines(
	b *game.Board,
	history []string,
	aiName string,
	aiEval string,
) []string {
	hud := NewTable()
	hud.BorderStyle = UnicodeBorders
	hud.SetTitle(" PAINEL DE CONTROLE ")
	hud.SetHeaders("PROPRIEDADE", "VALOR")
	hud.SetAlignments(AlignLeft, AlignLeft)

	turnStr := Colorize(FgBrightCyan+Bold, "Jogador (Brancas ●)")
	if b.Turn == game.Black {
		turnStr = Colorize(FgBrightRed+Bold, "IA "+aiName+" (Pretas ○)")
	}

	hud.AddRow("Vez da Jogada", turnStr)
	hud.AddRow("Motor de IA", Colorize(FgBrightYellow, aiName))
	hud.AddRow("Dimensao", fmt.Sprintf("%dx%d (%d pecas cada)", b.Size, b.Size, ((b.Size-2)/2)*(b.Size/2)))
	hud.AddRow("Brancas (Voce)", fmt.Sprintf("Total: %d | Damas: %d", b.WhiteCount, b.WhiteKingCount))
	hud.AddRow("Pretas (IA)", fmt.Sprintf("Total: %d | Damas: %d", b.BlackCount, b.BlackKingCount))
	hud.AddRow("Total Jogadas", fmt.Sprintf("%d", b.MoveCount))

	if aiEval != "" {
		hud.AddRow("Analise da IA", Colorize(FgBrightGreen, aiEval))
	}

	start := 0
	if len(history) > 4 {
		start = len(history) - 4
	}
	histSlice := history[start:]
	if len(histSlice) == 0 {
		hud.AddRow("Historico (4 ult.)", "(nenhuma jogada)")
	} else {
		for i, h := range histSlice {
			idx := start + i + 1
			hud.AddRow(fmt.Sprintf("Jogada #%d", idx), h)
		}
	}

	rendered := hud.Render()
	return strings.Split(strings.TrimRight(rendered, "\n"), "\n")
}

func (r *Renderer) renderStatusBanner(message string, targetWidth int) string {
	banner := NewTable()
	banner.BorderStyle = UnicodeBorders
	banner.Padding = 1

	coloredMsg := Colorize(FgBrightYellow+Bold, "» "+message)
	banner.AddRow(coloredMsg)

	return strings.TrimRight(banner.Render(), "\n")
}
