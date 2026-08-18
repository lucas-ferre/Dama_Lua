package tests

import (
	"strings"
	"testing"

	"damas-go/pkg/terminal"
)

func TestTableRender(t *testing.T) {
	table := terminal.NewTable()
	table.SetTitle("TABELA TESTE")
	table.SetHeaders("ID", "NOME", "STATUS")
	table.AddRow("1", "Processo MDP", "Ativo")
	table.AddRow("2", "Busca A*", "Pronto")

	output := table.Render()
	if !strings.Contains(output, "TABELA TESTE") {
		t.Fatalf("titulo ausente no output renderizado")
	}
	if !strings.Contains(output, "Processo MDP") {
		t.Fatalf("linha 1 ausente no output renderizado")
	}
	if !strings.Contains(output, "Busca A*") {
		t.Fatalf("linha 2 ausente no output renderizado")
	}
}

func TestVisibleLen(t *testing.T) {
	colored := terminal.Colorize(terminal.FgRed+terminal.Bold, "Texto")
	vLen := terminal.VisibleLen(colored)
	if vLen != 5 {
		t.Fatalf("esperado comprimento visivel 5, obtido %d", vLen)
	}
}
