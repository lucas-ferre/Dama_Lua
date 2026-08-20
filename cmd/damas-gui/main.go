package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"damas-go/pkg/gui"
)

func main() {
	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Damas Go - Jogo de Damas com IA")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := gui.NewApp()
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
