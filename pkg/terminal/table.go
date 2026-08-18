package terminal

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

type BorderStyle struct {
	TopLeft     string
	TopMid      string
	TopRight    string
	MidLeft     string
	MidMid      string
	MidRight    string
	BotLeft     string
	BotMid      string
	BotRight    string
	Horizontal  string
	Vertical    string
}

var (
	UnicodeBorders = BorderStyle{
		TopLeft:    "┌",
		TopMid:     "┬",
		TopRight:   "┐",
		MidLeft:    "├",
		MidMid:     "┼",
		MidRight:   "┤",
		BotLeft:    "└",
		BotMid:     "┴",
		BotRight:   "┘",
		Horizontal: "─",
		Vertical:   "│",
	}

	UnicodeDoubleBorders = BorderStyle{
		TopLeft:    "╔",
		TopMid:     "╦",
		TopRight:   "╗",
		MidLeft:    "╠",
		MidMid:     "╬",
		MidRight:   "╣",
		BotLeft:    "╚",
		BotMid:     "╩",
		BotRight:   "╝",
		Horizontal: "═",
		Vertical:   "║",
	}

	AsciiBorders = BorderStyle{
		TopLeft:    "+",
		TopMid:     "+",
		TopRight:   "+",
		MidLeft:    "+",
		MidMid:     "+",
		MidRight:   "+",
		BotLeft:    "+",
		BotMid:     "+",
		BotRight:   "+",
		Horizontal: "-",
		Vertical:   "|",
	}
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func VisibleLen(s string) int {
	clean := ansiRegex.ReplaceAllString(s, "")
	return utf8.RuneCountInString(clean)
}

type Table struct {
	Title       string
	Headers     []string
	Rows        [][]string
	Alignments  []Alignment
	BorderStyle BorderStyle
	Padding     int
}

func NewTable() *Table {
	return &Table{
		Headers:     nil,
		Rows:        nil,
		Alignments:  nil,
		BorderStyle: UnicodeBorders,
		Padding:     1,
	}
}

func (t *Table) SetTitle(title string) *Table {
	t.Title = title
	return t
}

func (t *Table) SetHeaders(headers ...string) *Table {
	t.Headers = headers
	if len(t.Alignments) < len(headers) {
		for len(t.Alignments) < len(headers) {
			t.Alignments = append(t.Alignments, AlignLeft)
		}
	}
	return t
}

func (t *Table) AddRow(cells ...string) *Table {
	t.Rows = append(t.Rows, cells)
	if len(t.Alignments) < len(cells) {
		for len(t.Alignments) < len(cells) {
			t.Alignments = append(t.Alignments, AlignLeft)
		}
	}
	return t
}

func (t *Table) SetAlignments(aligns ...Alignment) *Table {
	t.Alignments = aligns
	return t
}

func (t *Table) Render() string {
	numCols := len(t.Headers)
	for _, row := range t.Rows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}

	if numCols == 0 {
		return ""
	}

	colWidths := make([]int, numCols)
	for i, h := range t.Headers {
		vl := VisibleLen(h)
		if vl > colWidths[i] {
			colWidths[i] = vl
		}
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < numCols {
				vl := VisibleLen(cell)
				if vl > colWidths[i] {
					colWidths[i] = vl
				}
			}
		}
	}

	totalInnerWidth := 0
	for _, w := range colWidths {
		totalInnerWidth += w + (t.Padding * 2)
	}
	totalInnerWidth += (numCols - 1)

	var sb strings.Builder

	if t.Title != "" {
		sb.WriteString(t.BorderStyle.TopLeft)
		sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, totalInnerWidth))
		sb.WriteString(t.BorderStyle.TopRight + "\n")

		titlePad := totalInnerWidth - VisibleLen(t.Title)
		leftPad := titlePad / 2
		rightPad := titlePad - leftPad
		sb.WriteString(t.BorderStyle.Vertical)
		sb.WriteString(strings.Repeat(" ", leftPad))
		sb.WriteString(t.Title)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(t.BorderStyle.Vertical + "\n")

		sb.WriteString(t.BorderStyle.MidLeft)
		for i, w := range colWidths {
			sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, w+t.Padding*2))
			if i < numCols-1 {
				sb.WriteString(t.BorderStyle.TopMid)
			}
		}
		sb.WriteString(t.BorderStyle.MidRight + "\n")
	} else {
		sb.WriteString(t.BorderStyle.TopLeft)
		for i, w := range colWidths {
			sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, w+t.Padding*2))
			if i < numCols-1 {
				sb.WriteString(t.BorderStyle.TopMid)
			}
		}
		sb.WriteString(t.BorderStyle.TopRight + "\n")
	}

	if len(t.Headers) > 0 {
		sb.WriteString(t.BorderStyle.Vertical)
		for i := 0; i < numCols; i++ {
			headerText := ""
			if i < len(t.Headers) {
				headerText = t.Headers[i]
			}
			align := AlignCenter
			if i < len(t.Alignments) {
				align = t.Alignments[i]
			}
			sb.WriteString(strings.Repeat(" ", t.Padding))
			sb.WriteString(formatCell(headerText, colWidths[i], align))
			sb.WriteString(strings.Repeat(" ", t.Padding))
			sb.WriteString(t.BorderStyle.Vertical)
		}
		sb.WriteString("\n")

		sb.WriteString(t.BorderStyle.MidLeft)
		for i, w := range colWidths {
			sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, w+t.Padding*2))
			if i < numCols-1 {
				sb.WriteString(t.BorderStyle.MidMid)
			}
		}
		sb.WriteString(t.BorderStyle.MidRight + "\n")
	}

	for rIdx, row := range t.Rows {
		sb.WriteString(t.BorderStyle.Vertical)
		for i := 0; i < numCols; i++ {
			cellText := ""
			if i < len(row) {
				cellText = row[i]
			}
			align := AlignLeft
			if i < len(t.Alignments) {
				align = t.Alignments[i]
			}
			sb.WriteString(strings.Repeat(" ", t.Padding))
			sb.WriteString(formatCell(cellText, colWidths[i], align))
			sb.WriteString(strings.Repeat(" ", t.Padding))
			sb.WriteString(t.BorderStyle.Vertical)
		}
		sb.WriteString("\n")

		if rIdx < len(t.Rows)-1 {
			sb.WriteString(t.BorderStyle.MidLeft)
			for i, w := range colWidths {
				sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, w+t.Padding*2))
				if i < numCols-1 {
					sb.WriteString(t.BorderStyle.MidMid)
				}
			}
			sb.WriteString(t.BorderStyle.MidRight + "\n")
		}
	}

	sb.WriteString(t.BorderStyle.BotLeft)
	for i, w := range colWidths {
		sb.WriteString(strings.Repeat(t.BorderStyle.Horizontal, w+t.Padding*2))
		if i < numCols-1 {
			sb.WriteString(t.BorderStyle.BotMid)
		}
	}
	sb.WriteString(t.BorderStyle.BotRight + "\n")

	return sb.String()
}

func formatCell(text string, width int, align Alignment) string {
	vLen := VisibleLen(text)
	if vLen >= width {
		return text
	}

	diff := width - vLen
	switch align {
	case AlignRight:
		return strings.Repeat(" ", diff) + text
	case AlignCenter:
		left := diff / 2
		right := diff - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default:
		return text + strings.Repeat(" ", diff)
	}
}
