package sheets

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var readBuildInfo = debug.ReadBuildInfo
var fixedCellWidth bool

const helpText = `JSON viewer in your terminal.

USAGE
  sheets <file.json>
      Launch sheets interactively with a JSON file.

OPTIONS
  -h, --help
      Show this help message.

  -v, --version
      Show the current version.

  -w N, --width N
      Set the column cell width (default: 40).

EXAMPLES
  sheets data.json
  sheets config.json
`

func maybeHandleTopLevelOption(args []string, stdout io.Writer) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {
	case "-h", "--help", "-help":
		_, err := io.WriteString(stdout, helpText)
		return true, err
	case "-v", "--version", "-version":
		_, err := fmt.Fprintf(stdout, "sheets %s\n", buildVersion())
		return true, err
	default:
		return false, nil
	}
}

func buildVersion() string {
	info, ok := readBuildInfo()
	if !ok {
		return "dev"
	}

	version := strings.TrimSpace(info.Main.Version)
	if version == "" || version == "(devel)" {
		return "dev"
	}

	return version
}

func run(args []string, stdout io.Writer, cellWidth int) error {
	if handled, err := maybeHandleTopLevelOption(args, stdout); handled || err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("usage: sheets <file.json>")
	}

	path := args[0]

	m := newModel()
	if cellWidth > 0 {
		m.cellWidth = cellWidth
	}

	if err := m.loadJSONFile(path); err != nil {
		return err
	}

	options := []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
	program := tea.NewProgram(m, options...)
	_, err := program.Run()
	return err
}

func Main(args []string, stdout, stderr io.Writer) int {
	cellWidth, args := parseCellWidthFlag(args)
	err := run(args, stdout, cellWidth)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func parseCellWidthFlag(args []string) (int, []string) {
	remaining := make([]string, 0, len(args))
	width := 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "-w" || arg == "--width") && i+1 < len(args) {
			if n := parsePositiveInt(args[i+1]); n > 0 {
				width = n
				i++
				continue
			}
		}
		if strings.HasPrefix(arg, "--width=") {
			if n := parsePositiveInt(strings.TrimPrefix(arg, "--width=")); n > 0 {
				width = n
				continue
			}
		}
		if strings.HasPrefix(arg, "-w=") {
			if n := parsePositiveInt(strings.TrimPrefix(arg, "-w=")); n > 0 {
				width = n
				continue
			}
		}
		remaining = append(remaining, arg)
	}
	// set global variable flag that we'll use fixden width
	fixedCellWidth = true
	return width, remaining
}

func parsePositiveInt(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
