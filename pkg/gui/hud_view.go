package gui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"damas-go/pkg/game"
)

type HUDView struct {
	X           float32
	Y           float32
	W           float32
	H           float32
	btnNewGame  *Button
	btnRestart  *Button
	OnNewGame   func()
	OnRestart   func()
}

func NewHUDView(x, y, w, h float32, onNewGame, onRestart func()) *HUDView {
	hud := &HUDView{
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		OnNewGame:  onNewGame,
		OnRestart:  onRestart,
		btnNewGame: &Button{X: x + 15, Y: y + h - 50, W: 135, H: 36, Label: "Menu Inicial"},
		btnRestart: &Button{X: x + 165, Y: y + h - 50, W: 135, H: 36, Label: "Reiniciar"},
	}

	hud.btnNewGame.OnClick = onNewGame
	hud.btnRestart.OnClick = onRestart

	return hud
}

func (hud *HUDView) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if hud.btnNewGame.Contains(mx, my) && hud.btnNewGame.OnClick != nil {
			hud.btnNewGame.OnClick()
		}
		if hud.btnRestart.Contains(mx, my) && hud.btnRestart.OnClick != nil {
			hud.btnRestart.OnClick()
		}
	}
}

func (hud *HUDView) Draw(
	dst *ebiten.Image,
	b *game.Board,
	playerColor game.Color,
	aiName string,
	aiEval string,
	history []string,
	isAIThinking bool,
	isGameOver bool,
	winner game.Color,
) {
	mx, my := ebiten.CursorPosition()

	DrawRoundedRect(dst, hud.X, hud.Y, hud.W, hud.H, 10, ColorPanelBg)
	DrawRing(dst, hud.X+hud.W/2, hud.Y+hud.H/2, hud.W/2, 2, ColorPanelBorder)

	DrawText(dst, "PAINEL DE CONTROLE", int(hud.X+20), int(hud.Y+30), ColorAccentYellow)
	DrawRect(dst, hud.X+20, hud.Y+38, hud.W-40, 1, ColorPanelBorder)

	turnText := "Sua Vez (Voce)"
	turnClr := ColorAccentGreen
	if b.Turn != playerColor {
		turnText = "IA Pensando..."
		turnClr = ColorAccentRed
		if !isAIThinking {
			turnText = "Vez da IA"
		}
	}
	if isGameOver {
		turnText = "Partida Encerrada"
		turnClr = ColorAccentYellow
	}

	DrawText(dst, "TURNO ATUAL:", int(hud.X+20), int(hud.Y+65), ColorTextDim)
	DrawText(dst, turnText, int(hud.X+130), int(hud.Y+65), turnClr)

	pName := "Brancas (●)"
	aiColorName := "Pretas (○)"
	pPieces := b.WhiteCount
	pKings := b.WhiteKingCount
	aiPieces := b.BlackCount
	aiKings := b.BlackKingCount

	if playerColor == game.Black {
		pName = "Pretas (○)"
		aiColorName = "Brancas (●)"
		pPieces = b.BlackCount
		pKings = b.BlackKingCount
		aiPieces = b.WhiteCount
		aiKings = b.WhiteKingCount
	}

	DrawText(dst, "VOCE ("+pName+"):", int(hud.X+20), int(hud.Y+95), ColorTextDim)
	DrawText(dst, fmt.Sprintf("%d pecas (Damas: %d)", pPieces, pKings), int(hud.X+20), int(hud.Y+112), ColorTextLight)

	DrawText(dst, "IA ("+aiColorName+"):", int(hud.X+20), int(hud.Y+140), ColorTextDim)
	DrawText(dst, fmt.Sprintf("%d pecas (Damas: %d)", aiPieces, aiKings), int(hud.X+20), int(hud.Y+157), ColorTextLight)

	DrawText(dst, "MOTOR DE IA:", int(hud.X+20), int(hud.Y+185), ColorTextDim)
	DrawText(dst, aiName, int(hud.X+20), int(hud.Y+202), ColorAccentYellow)

	DrawText(dst, "ANALISE / ESTADO DA IA:", int(hud.X+20), int(hud.Y+230), ColorTextDim)
	if aiEval == "" {
		aiEval = "(aguardando proxima jogada)"
	}
	DrawText(dst, aiEval, int(hud.X+20), int(hud.Y+247), ColorAccentGreen)

	DrawRect(dst, hud.X+20, hud.Y+270, hud.W-40, 1, ColorPanelBorder)

	DrawText(dst, "HISTORICO (4 ULTIMAS):", int(hud.X+20), int(hud.Y+295), ColorTextDim)

	start := 0
	if len(history) > 4 {
		start = len(history) - 4
	}
	histSlice := history[start:]

	if len(histSlice) == 0 {
		DrawText(dst, "(nenhuma jogada realizada)", int(hud.X+20), int(hud.Y+320), ColorTextDim)
	} else {
		for i, h := range histSlice {
			idx := start + i + 1
			line := fmt.Sprintf("#%d %s", idx, h)
			DrawText(dst, line, int(hud.X+20), int(hud.Y+320+float32(i*22)), ColorTextLight)
		}
	}

	hud.btnNewGame.Draw(dst, mx, my)
	hud.btnRestart.Draw(dst, mx, my)

	if isGameOver {
		hud.drawGameOverOverlay(dst, winner, playerColor)
	}
}

func (hud *HUDView) drawGameOverOverlay(dst *ebiten.Image, winner, playerColor game.Color) {
	cardX := hud.X + 15
	cardY := hud.Y + 410
	cardW := hud.W - 30
	cardH := float32(65)

	cardClr := ColorPanelBg
	textClr := ColorAccentYellow
	title := "FIM DE JOGO: EMPATE!"

	if winner == playerColor {
		cardClr = color.RGBA{25, 75, 45, 255}
		textClr = ColorAccentGreen
		title = "PARABENS! VOCE VENCEU!"
	} else if winner == playerColor.Opponent() {
		cardClr = color.RGBA{85, 30, 30, 255}
		textClr = ColorAccentRed
		title = "VITORIA DA IA!"
	}

	DrawRoundedRect(dst, cardX, cardY, cardW, cardH, 8, cardClr)
	DrawText(dst, title, int(cardX+float32(cardW-float32(len(title)*7))/2.0), int(cardY+28), textClr)
	sub := "Clique em Menu para nova partida"
	DrawText(dst, sub, int(cardX+float32(cardW-float32(len(sub)*7))/2.0), int(cardY+48), ColorTextDim)
}
