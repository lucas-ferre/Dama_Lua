package gui

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"damas-go/pkg/ai"
	"damas-go/pkg/game"
)

type AppState int

const (
	StateMenu AppState = iota
	StatePlaying
)

type aiResult struct {
	move game.Move
	eval string
}

type App struct {
	state                  AppState
	menu                   *MenuView
	boardView              *BoardView
	hudView                *HUDView
	board                  *game.Board
	rules                  *game.RulesEngine
	bot                    ai.Bot
	parser                 *game.NotationParser
	playerColor            game.Color
	selectedPos            *game.Position
	validMovesFromSelected []game.Move
	allLegalMoves          []game.Move
	history                []string
	lastAiEval             string
	isAIThinking           bool
	isGameOver             bool
	winner                 game.Color
	aiChan                 chan aiResult
	lastSize               int
	lastBotType            ai.BotType
	lastDifficulty         int
}

func NewApp() *App {
	app := &App{
		state:  StateMenu,
		rules:  game.NewRulesEngine(),
		aiChan: make(chan aiResult, 1),
	}

	app.menu = NewMenuView(app.StartNewGame)
	return app
}

func (app *App) StartNewGame(size int, playerColor game.Color, botType ai.BotType, difficulty int) {
	app.lastSize = size
	app.playerColor = playerColor
	app.lastBotType = botType
	app.lastDifficulty = difficulty

	app.board = game.NewBoard(size)
	app.bot = ai.CreateBot(botType, difficulty)
	app.parser = game.NewNotationParser(size)
	app.boardView = NewBoardView(40, 40, 560, size, playerColor)
	app.hudView = NewHUDView(640, 30, 300, 580, app.BackToMenu, app.RestartGame)

	app.selectedPos = nil
	app.validMovesFromSelected = nil
	app.history = nil
	app.lastAiEval = ""
	app.isAIThinking = false
	app.isGameOver = false
	app.winner = game.None

	app.allLegalMoves = app.rules.GetLegalMoves(app.board, app.board.Turn)
	app.state = StatePlaying
}

func (app *App) BackToMenu() {
	app.state = StateMenu
}

func (app *App) RestartGame() {
	app.StartNewGame(app.lastSize, app.playerColor, app.lastBotType, app.lastDifficulty)
}

func (app *App) Update() error {
	if app.state == StateMenu {
		app.menu.Update()
		return nil
	}

	app.hudView.Update()

	select {
	case res := <-app.aiChan:
		app.isAIThinking = false
		app.lastAiEval = res.eval

		if res.move.From != res.move.To || len(res.move.Path) > 0 {
			formatted := app.parser.FormatMove(res.move)
			app.history = append(app.history, fmt.Sprintf("IA: %s", formatted))
			app.board.ApplyMove(res.move)
		}

		app.checkGameState()
	default:
	}

	if !app.isGameOver && app.board.Turn != app.playerColor && !app.isAIThinking {
		app.isAIThinking = true
		boardClone := app.board.Clone()
		botInstance := app.bot

		go func() {
			time.Sleep(250 * time.Millisecond)
			move, eval := botInstance.SelectMove(boardClone)
			app.aiChan <- aiResult{move: move, eval: eval}
		}()
	}

	if !app.isGameOver && app.board.Turn == app.playerColor && !app.isAIThinking {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			clickedPos := app.boardView.PosAtScreen(mx, my)

			if clickedPos != nil {
				app.handlePlayerBoardClick(*clickedPos)
			}
		}
	}

	return nil
}

func (app *App) handlePlayerBoardClick(pos game.Position) {
	piece := app.board.Get(pos)

	if piece.Color == app.playerColor {
		var movesFromPos []game.Move
		for _, m := range app.allLegalMoves {
			if m.From == pos {
				movesFromPos = append(movesFromPos, m)
			}
		}

		if len(movesFromPos) > 0 {
			app.selectedPos = &pos
			app.validMovesFromSelected = movesFromPos
		}
		return
	}

	if app.selectedPos != nil {
		var matchedMove *game.Move
		for _, m := range app.validMovesFromSelected {
			if m.To == pos {
				matchedMove = &m
				break
			}
		}

		if matchedMove != nil {
			formatted := app.parser.FormatMove(*matchedMove)
			app.history = append(app.history, fmt.Sprintf("Voce: %s", formatted))

			app.board.ApplyMove(*matchedMove)
			app.selectedPos = nil
			app.validMovesFromSelected = nil

			app.checkGameState()
		}
	}
}

func (app *App) checkGameState() {
	isOver, winner := app.rules.IsGameOver(app.board)
	app.isGameOver = isOver
	app.winner = winner

	if !app.isGameOver {
		app.allLegalMoves = app.rules.GetLegalMoves(app.board, app.board.Turn)
		if len(app.allLegalMoves) == 0 {
			app.isGameOver = true
			app.winner = app.board.Turn.Opponent()
		}
	}
}

func (app *App) Draw(screen *ebiten.Image) {
	if app.state == StateMenu {
		app.menu.Draw(screen)
		return
	}

	screen.Fill(ColorBgDark)

	var mandatoryPieces []game.Position
	if app.board.Turn == app.playerColor {
		hasCaptures := false
		for _, m := range app.allLegalMoves {
			if m.IsCapture {
				hasCaptures = true
				break
			}
		}
		if hasCaptures {
			seen := make(map[game.Position]bool)
			for _, m := range app.allLegalMoves {
				if m.IsCapture && !seen[m.From] {
					seen[m.From] = true
					mandatoryPieces = append(mandatoryPieces, m.From)
				}
			}
		}
	}

	app.boardView.Draw(
		screen,
		app.board,
		app.selectedPos,
		app.validMovesFromSelected,
		mandatoryPieces,
		app.board.LastMove,
	)

	app.hudView.Draw(
		screen,
		app.board,
		app.playerColor,
		app.bot.Name(),
		app.lastAiEval,
		app.history,
		app.isAIThinking,
		app.isGameOver,
		app.winner,
	)
}

func (app *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 960, 640
}
