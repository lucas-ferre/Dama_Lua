package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"damas-go/pkg/ai"
	"damas-go/pkg/game"
)

type Button struct {
	X       float32
	Y       float32
	W       float32
	H       float32
	Label   string
	Active  bool
	OnClick func()
}

func (b *Button) Contains(px, py int) bool {
	x := float32(px)
	y := float32(py)
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

func (b *Button) Draw(dst *ebiten.Image, mouseX, mouseY int) {
	hover := b.Contains(mouseX, mouseY)
	clr := ColorBtnNormal
	if b.Active {
		clr = ColorBtnActive
	} else if hover {
		clr = ColorBtnHover
	}

	DrawRoundedRect(dst, b.X, b.Y, b.W, b.H, 6, clr)

	textX := int(b.X + (b.W-float32(len(b.Label)*7))/2.0)
	textY := int(b.Y + b.H/2.0 + 4)
	DrawText(dst, b.Label, textX, textY, ColorTextLight)
}

type MenuView struct {
	SelectedSize       int
	SelectedColor      game.Color
	SelectedBotType    ai.BotType
	SelectedDifficulty int
	OnStartGame        func(size int, playerColor game.Color, botType ai.BotType, difficulty int)
	buttons            []*Button
	startButton        *Button
}

func NewMenuView(onStart func(size int, playerColor game.Color, botType ai.BotType, difficulty int)) *MenuView {
	mv := &MenuView{
		SelectedSize:       10,
		SelectedColor:      game.White,
		SelectedBotType:    ai.BotTypeHybrid,
		SelectedDifficulty: 2,
		OnStartGame:        onStart,
	}
	mv.initButtons()
	return mv
}

func (mv *MenuView) initButtons() {
	mv.buttons = nil

	centerX := float32(480)

	btnSize10 := &Button{X: centerX - 220, Y: 130, W: 210, H: 36, Label: "10x10 (Ampliado)", Active: mv.SelectedSize == 10}
	btnSize8 := &Button{X: centerX + 10, Y: 130, W: 210, H: 36, Label: "8x8 (Classico)", Active: mv.SelectedSize == 8}
	btnSize10.OnClick = func() { mv.SelectedSize = 10; mv.updateActiveStates() }
	btnSize8.OnClick = func() { mv.SelectedSize = 8; mv.updateActiveStates() }
	mv.buttons = append(mv.buttons, btnSize10, btnSize8)

	btnWhite := &Button{X: centerX - 220, Y: 210, W: 210, H: 36, Label: "Brancas (Joga Primeiro)", Active: mv.SelectedColor == game.White}
	btnBlack := &Button{X: centerX + 10, Y: 210, W: 210, H: 36, Label: "Pretas (IA Joga Primeiro)", Active: mv.SelectedColor == game.Black}
	btnWhite.OnClick = func() { mv.SelectedColor = game.White; mv.updateActiveStates() }
	btnBlack.OnClick = func() { mv.SelectedColor = game.Black; mv.updateActiveStates() }
	mv.buttons = append(mv.buttons, btnWhite, btnBlack)

	btnHybrid := &Button{X: centerX - 220, Y: 290, W: 210, H: 36, Label: "Hibrido (MDP + A* + HC)", Active: mv.SelectedBotType == ai.BotTypeHybrid}
	btnMDP := &Button{X: centerX + 10, Y: 290, W: 210, H: 36, Label: "Markov (MDP)", Active: mv.SelectedBotType == ai.BotTypeMDP}
	btnAStar := &Button{X: centerX - 220, Y: 335, W: 210, H: 36, Label: "Busca A* (Tatica)", Active: mv.SelectedBotType == ai.BotTypeAStar}
	btnHC := &Button{X: centerX + 10, Y: 335, W: 210, H: 36, Label: "Hill Climbing (Restarts)", Active: mv.SelectedBotType == ai.BotTypeHillClimbing}

	btnHybrid.OnClick = func() { mv.SelectedBotType = ai.BotTypeHybrid; mv.updateActiveStates() }
	btnMDP.OnClick = func() { mv.SelectedBotType = ai.BotTypeMDP; mv.updateActiveStates() }
	btnAStar.OnClick = func() { mv.SelectedBotType = ai.BotTypeAStar; mv.updateActiveStates() }
	btnHC.OnClick = func() { mv.SelectedBotType = ai.BotTypeHillClimbing; mv.updateActiveStates() }
	mv.buttons = append(mv.buttons, btnHybrid, btnMDP, btnAStar, btnHC)

	btnEasy := &Button{X: centerX - 220, Y: 415, W: 135, H: 36, Label: "Facil", Active: mv.SelectedDifficulty == 1}
	btnMedium := &Button{X: centerX - 70, Y: 415, W: 140, H: 36, Label: "Medio", Active: mv.SelectedDifficulty == 2}
	btnHard := &Button{X: centerX + 85, Y: 415, W: 135, H: 36, Label: "Dificil", Active: mv.SelectedDifficulty == 3}

	btnEasy.OnClick = func() { mv.SelectedDifficulty = 1; mv.updateActiveStates() }
	btnMedium.OnClick = func() { mv.SelectedDifficulty = 2; mv.updateActiveStates() }
	btnHard.OnClick = func() { mv.SelectedDifficulty = 3; mv.updateActiveStates() }
	mv.buttons = append(mv.buttons, btnEasy, btnMedium, btnHard)

	mv.startButton = &Button{
		X:      centerX - 160,
		Y:      500,
		W:      320,
		H:      50,
		Label:  "INICIAR PARTIDA",
		Active: false,
		OnClick: func() {
			if mv.OnStartGame != nil {
				mv.OnStartGame(mv.SelectedSize, mv.SelectedColor, mv.SelectedBotType, mv.SelectedDifficulty)
			}
		},
	}
}

func (mv *MenuView) updateActiveStates() {
	mv.initButtons()
}

func (mv *MenuView) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, b := range mv.buttons {
			if b.Contains(mx, my) && b.OnClick != nil {
				b.OnClick()
				return
			}
		}
		if mv.startButton.Contains(mx, my) && mv.startButton.OnClick != nil {
			mv.startButton.OnClick()
			return
		}
	}
}

func (mv *MenuView) Draw(dst *ebiten.Image) {
	dst.Fill(ColorBgDark)

	mx, my := ebiten.CursorPosition()
	centerX := float32(480)

	DrawRoundedRect(dst, centerX-260, 40, 520, 530, 12, ColorPanelBg)
	DrawRing(dst, centerX, 305, 270, 2, ColorPanelBorder)

	title := "JOGO DE DAMAS EM GO (IA)"
	DrawText(dst, title, int(centerX-float32(len(title)*7)/2.0), 75, ColorAccentYellow)

	sub := "Regras Brasileiras / MDP / Busca A* / Hill Climbing"
	DrawText(dst, sub, int(centerX-float32(len(sub)*7)/2.0), 95, ColorTextDim)

	DrawText(dst, "1. DIMENSAO DO TABULEIRO:", int(centerX-220), 122, ColorTextLight)
	DrawText(dst, "2. SUA COR (QUEM JOGA PRIMEIRO):", int(centerX-220), 202, ColorTextLight)
	DrawText(dst, "3. MOTOR DE INTELIGENCIA ARTIFICIAL:", int(centerX-220), 282, ColorTextLight)
	DrawText(dst, "4. NIVEL DE DIFICULDADE:", int(centerX-220), 407, ColorTextLight)

	for _, b := range mv.buttons {
		b.Draw(dst, mx, my)
	}

	btnHover := mv.startButton.Contains(mx, my)
	btnClr := ColorAccentGreen
	if btnHover {
		btnClr = color.RGBA{100, 230, 150, 255}
	}
	DrawRoundedRect(dst, mv.startButton.X, mv.startButton.Y, mv.startButton.W, mv.startButton.H, 8, btnClr)

	sLabel := mv.startButton.Label
	textX := int(mv.startButton.X + (mv.startButton.W-float32(len(sLabel)*7))/2.0)
	textY := int(mv.startButton.Y + mv.startButton.H/2.0 + 4)
	DrawText(dst, sLabel, textX, textY, ColorBgDark)
}
