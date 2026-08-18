package terminal

import "fmt"

const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Underline   = "\033[4m"
	
	FgBlack     = "\033[30m"
	FgRed       = "\033[31m"
	FgGreen     = "\033[32m"
	FgYellow    = "\033[33m"
	FgBlue      = "\033[34m"
	FgMagenta   = "\033[35m"
	FgCyan      = "\033[36m"
	FgWhite     = "\033[37m"
	
	FgBrightBlack   = "\033[90m"
	FgBrightRed     = "\033[91m"
	FgBrightGreen   = "\033[92m"
	FgBrightYellow  = "\033[93m"
	FgBrightBlue    = "\033[94m"
	FgBrightMagenta = "\033[95m"
	FgBrightCyan    = "\033[96m"
	FgBrightWhite   = "\033[97m"
	
	BgBlack     = "\033[40m"
	BgRed       = "\033[41m"
	BgGreen     = "\033[42m"
	BgYellow    = "\033[43m"
	BgBlue      = "\033[44m"
	BgMagenta   = "\033[45m"
	BgCyan      = "\033[46m"
	BgWhite     = "\033[47m"
	
	BgDarkSquare  = "\033[48;5;236m"
	BgLightSquare = "\033[48;5;244m"
	BgHighlight   = "\033[48;5;28m"
	BgLastMove    = "\033[48;5;58m"
)

func Colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}
