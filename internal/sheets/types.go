package sheets

import (
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultRows = 999
	maxRows     = 50000
	totalCols   = 2 // KEYS | VALUES
)

type mode string

const (
	normalMode  mode = "NORMAL"
	insertMode  mode = "INSERT"
	selectMode  mode = "SELECT"
	commandMode mode = ":"
)

type cellKey struct {
	row int
	col int
}

type clipboard struct {
	cells     [][]string
	sourceRow int
	sourceCol int
}

type promptKind rune

const (
	noPrompt             promptKind = 0
	commandPrompt        promptKind = ':'
	searchForwardPrompt  promptKind = '/'
	searchBackwardPrompt promptKind = '?'
)

type undoState struct {
	cells       map[cellKey]string
	rowCount    int
	selectedRow int
	selectedCol int
	selectRow   int
	selectCol   int
	selectRows  bool
	rowOffset   int
	colOffset   int
}

// jsonShape is kept for write-back compatibility.
type jsonShape int

const (
	jsonShapeUnknown jsonShape = iota
	jsonShapeObjectArray
	jsonShapeSingleObject
	jsonShapeScalarArray
)

type model struct {
	width    int
	height   int
	rowCount int

	// jsonKeys holds the original JSON field-name order.
	jsonKeys  []string
	jsonShape jsonShape

	mode mode

	promptKind      promptKind
	gotoPending     bool
	gotoBuffer      string
	commandPending  bool
	commandBuffer   string
	commandCursor   int
	commandMessage  string
	deletePending   bool
	yankPending     bool
	yankCount       int
	zPending        bool
	registerPending bool
	activeRegister  rune
	countBuffer     string
	currentFilePath string
	searchQuery     string
	searchDirection int
	markPending     bool
	markJumpPending bool
	markJumpExact   bool
	commandError    bool
	dirtyFile       bool

	selectedRow int
	selectedCol int
	selectRow   int
	selectCol   int
	selectRows  bool
	rowOffset   int
	colOffset   int

	cellWidth     int
	rowLabelWidth int

	cells           map[cellKey]string
	copyBuffer      clipboard
	hasCopyBuffer   bool
	registers       map[rune]clipboard
	marks           map[rune]cellKey
	jumpBack        []cellKey
	jumpForward     []cellKey
	undoStack       []undoState
	redoStack       []undoState
	editingValue    string
	editingCursor   int
	editCursor      cursor.Model
	insertKeys      []tea.KeyMsg
	recordingInsert bool
	lastChange      []tea.KeyMsg
	replayingChange bool

	headerStyle                   lipgloss.Style
	activeHeaderStyle             lipgloss.Style
	rowLabelStyle                 lipgloss.Style
	activeRowStyle                lipgloss.Style
	gridStyle                     lipgloss.Style
	formulaCellStyle              lipgloss.Style
	formulaErrorStyle             lipgloss.Style
	activeCellStyle               lipgloss.Style
	activeRowCellStyle            lipgloss.Style
	activeFormulaStyle            lipgloss.Style
	activeFormulaErrorStyle       lipgloss.Style
	selectCellStyle               lipgloss.Style
	selectFormulaStyle            lipgloss.Style
	selectFormulaErrorStyle       lipgloss.Style
	selectActiveCellStyle         lipgloss.Style
	selectActiveFormulaStyle      lipgloss.Style
	selectActiveFormulaErrorStyle lipgloss.Style
	selectHeaderStyle             lipgloss.Style
	selectActiveHeaderStyle       lipgloss.Style
	selectRowStyle                lipgloss.Style
	selectBorderStyle             lipgloss.Style
	statusBarStyle                lipgloss.Style
	statusTextStyle               lipgloss.Style
	statusAccentStyle             lipgloss.Style
	statusNormalStyle             lipgloss.Style
	statusInsertStyle             lipgloss.Style
	statusSelectStyle             lipgloss.Style
	commandLineStyle              lipgloss.Style
	commandErrorStyle             lipgloss.Style
}
