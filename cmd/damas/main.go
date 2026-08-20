package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"damas-go/pkg/ai"
	"damas-go/pkg/game"
	"damas-go/pkg/terminal"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printBanner()

		size := selectBoardSize(reader)
		playerColor := selectPlayerColor(reader)
		botType, difficulty := selectBotAndDifficulty(reader)

		bot := ai.CreateBot(botType, difficulty)
		board := game.NewBoard(size)
		rules := game.NewRulesEngine()
		parser := game.NewNotationParser(size)
		renderer := terminal.NewRenderer(size)

		var history []string
		lastAiEval := ""
		statusMessage := "Jogo iniciado! Boa sorte."

		for {
			clearScreen()
			fmt.Print(renderer.RenderGame(board, history, bot.Name(), lastAiEval, statusMessage))
			statusMessage = ""

			isOver, winner := rules.IsGameOver(board)
			if isOver {
				printGameOver(winner, playerColor)
				break
			}

			if board.Turn == playerColor {
				legalMoves := rules.GetLegalMoves(board, playerColor)
				if len(legalMoves) == 0 {
					statusMessage = "Voce nao possui movimentos legais."
					continue
				}

				colorName := "Brancas ●"
				if playerColor == game.Black {
					colorName = "Pretas ○"
				}

				fmt.Print(terminal.Colorize(terminal.FgBrightCyan+terminal.Bold, fmt.Sprintf("\nSua vez (%s)! Digite sua jogada (ex: E3 para F4 ou C3 D4) [ou 'sair']: ", colorName)))
				input, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				input = strings.TrimSpace(input)
				if strings.ToLower(input) == "sair" || strings.ToLower(input) == "exit" || strings.ToLower(input) == "q" {
					fmt.Println("\nPartida encerrada pelo jogador.")
					return
				}

				move, err := parser.ParseInput(input, legalMoves)
				if err != nil {
					statusMessage = fmt.Sprintf("Erro: %s", err.Error())
					continue
				}

				moveFormatted := parser.FormatMove(move)
				history = append(history, fmt.Sprintf("Jogador: %s", moveFormatted))
				board.ApplyMove(move)
			} else {
				aiColorName := "Brancas ●"
				if playerColor == game.White {
					aiColorName = "Pretas ○"
				}

				fmt.Print(terminal.Colorize(terminal.FgBrightRed+terminal.Bold, fmt.Sprintf("\nIA pensando (%s)... calculando melhor jogada...", aiColorName)))
				time.Sleep(300 * time.Millisecond)

				aiMove, aiEval := bot.SelectMove(board)
				lastAiEval = aiEval

				if aiMove.From == aiMove.To && len(aiMove.Path) == 0 {
					statusMessage = "IA sem movimentos."
					continue
				}

				moveFormatted := parser.FormatMove(aiMove)
				history = append(history, fmt.Sprintf("IA (%s): %s", bot.Name(), moveFormatted))
				board.ApplyMove(aiMove)
			}
		}

		fmt.Print("\nDeseja jogar novamente? (S/N): ")
		rematch, _ := reader.ReadString('\n')
		rematch = strings.TrimSpace(strings.ToUpper(rematch))
		if rematch != "S" && rematch != "SIM" && rematch != "Y" && rematch != "YES" {
			fmt.Println("\nObrigado por jogar Damas Go! Ate a proxima.")
			break
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner() {
	table := terminal.NewTable()
	table.BorderStyle = terminal.UnicodeDoubleBorders
	table.SetTitle(" JOGO DE DAMAS EM GO ")
	table.AddRow(terminal.Colorize(terminal.FgBrightYellow+terminal.Bold, "  Inteligencia Artificial com MDP, Busca A* e Hill Climbing  "))
	table.AddRow("  Regras Brasileiras / Notacao Algebrica de Xadrez (ex: E3 para F4)  ")
	fmt.Println(table.Render())
}

func selectBoardSize(r *bufio.Reader) int {
	menu := terminal.NewTable()
	menu.BorderStyle = terminal.UnicodeBorders
	menu.SetTitle(" DIMENSAO DO TABULEIRO ")
	menu.SetHeaders("OPCAO", "TAMANHO", "DESCRICAO")
	menu.AddRow("1", "10x10", "Matriz Espacial Ampliada (20 pecas cada) [Recomendado]")
	menu.AddRow("2", "8x8", "Matriz Espacial Classica (12 pecas cada)")
	fmt.Println(menu.Render())

	for {
		fmt.Print("Escolha a dimensao do tabuleiro [1-2, padrao: 1]: ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" || input == "1" {
			return 10
		}
		if input == "2" {
			return 8
		}
		fmt.Println("Opcao invalida. Tente novamente.")
	}
}

func selectPlayerColor(r *bufio.Reader) game.Color {
	menu := terminal.NewTable()
	menu.BorderStyle = terminal.UnicodeBorders
	menu.SetTitle(" ESCOLHA SUA COR ")
	menu.SetHeaders("OPCAO", "COR", "CONDICAO DE INICIO")
	menu.AddRow("1", "Brancas (●)", "Voce joga primeiro [Padrao]")
	menu.AddRow("2", "Pretas (○)", "A IA joga primeiro")
	fmt.Println(menu.Render())

	for {
		fmt.Print("Escolha a sua cor [1-2, padrao: 1]: ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" || input == "1" {
			return game.White
		}
		if input == "2" {
			return game.Black
		}
		fmt.Println("Opcao invalida. Tente novamente.")
	}
}

func selectBotAndDifficulty(r *bufio.Reader) (ai.BotType, int) {
	menuAI := terminal.NewTable()
	menuAI.BorderStyle = terminal.UnicodeBorders
	menuAI.SetTitle(" MOTOR DE INTELIGENCIA ARTIFICIAL ")
	menuAI.SetHeaders("OPCAO", "MOTOR DE IA", "CARACTERISTICA")
	menuAI.AddRow("1", "Modo Hibrido Mestre", "Combina A* (Taticas), MDP (Estrategia) e Hill Climbing [Recomendado]")
	menuAI.AddRow("2", "Processo de Decisao de Markov (MDP)", "Modelagem probabilistica de transicoes e Bellman Value Iteration")
	menuAI.AddRow("3", "Busca A* (A-Star)", "Busca em arvore tica e filas de prioridade com Min-Heap")
	menuAI.AddRow("4", "Hill Climbing com Reinicializacao", "Otimizacao heuristica local com Random Restarts")
	fmt.Println(menuAI.Render())

	var chosenAI ai.BotType
	for {
		fmt.Print("Escolha o Motor de IA do adversario [1-4, padrao: 1]: ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" || input == "1" {
			chosenAI = ai.BotTypeHybrid
			break
		}
		val, err := strconv.Atoi(input)
		if err == nil && val >= 1 && val <= 4 {
			switch val {
			case 2:
				chosenAI = ai.BotTypeMDP
			case 3:
				chosenAI = ai.BotTypeAStar
			case 4:
				chosenAI = ai.BotTypeHillClimbing
			default:
				chosenAI = ai.BotTypeHybrid
			}
			break
		}
		fmt.Println("Opcao invalida. Tente novamente.")
	}

	menuDiff := terminal.NewTable()
	menuDiff.BorderStyle = terminal.UnicodeBorders
	menuDiff.SetTitle(" NIVEL DE DIFICULDADE ")
	menuDiff.SetHeaders("OPCAO", "NIVEL", "PROFUNDIDADE")
	menuDiff.AddRow("1", "Facil", "Exploracao rapida e baixa profundidade")
	menuDiff.AddRow("2", "Medio", "Equilibrio tatico e estrategico [Padrao]")
	menuDiff.AddRow("3", "Dificil", "Calculo profundo de variantes e multiplas reinicializacoes")
	fmt.Println(menuDiff.Render())

	difficulty := 2
	for {
		fmt.Print("Escolha a dificuldade [1-3, padrao: 2]: ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" || input == "2" {
			difficulty = 2
			break
		}
		val, err := strconv.Atoi(input)
		if err == nil && val >= 1 && val <= 3 {
			difficulty = val
			break
		}
		fmt.Println("Opcao invalida. Tente novamente.")
	}

	return chosenAI, difficulty
}

func printGameOver(winner game.Color, playerColor game.Color) {
	banner := terminal.NewTable()
	banner.BorderStyle = terminal.UnicodeDoubleBorders
	banner.SetTitle(" FIM DE JOGO ")

	if winner == playerColor {
		banner.AddRow(terminal.Colorize(terminal.FgBrightGreen+terminal.Bold, "  PARABENS! VOCE VENCEU A PARTIDA!  "))
	} else if winner == playerColor.Opponent() {
		banner.AddRow(terminal.Colorize(terminal.FgBrightRed+terminal.Bold, "  VITORIA DA INTELIGENCIA ARTIFICIAL!  "))
	} else {
		banner.AddRow(terminal.Colorize(terminal.FgBrightYellow+terminal.Bold, "  PARTIDA EMPATADA!  "))
	}

	fmt.Println("\n" + banner.Render())
}
