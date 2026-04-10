package sheets

import (
	"maps"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newModel() model {
	insertAccent := lipgloss.Color("#D79921")
	selectAccent := lipgloss.Color("#2F66C7")
	statusSelectAccent := lipgloss.Color("13")
	formulaGreen := lipgloss.Color("2")
	errorRed := lipgloss.Color("9")
	gridGray := lipgloss.Color("8")
	selectBackground := lipgloss.Color("#264F78")
	white := lipgloss.Color("15")

	editCursor := cursor.New()
	editCursor.Style = lipgloss.NewStyle().Foreground(insertAccent)
	editCursor.TextStyle = lipgloss.NewStyle()
	editCursor.Blur()

	headerGray := lipgloss.Color("8")
	statusGray := lipgloss.Color("0")
	statusText := lipgloss.Color("7")
	statusAccent := insertAccent

	return model{
		mode:          normalMode,
		rowCount:      defaultRows,
		selectedRow:   0,
		selectedCol:   0,
		selectRow:     0,
		selectCol:     0,
		cellWidth:     40,
		rowLabelWidth: rowLabelWidthForCount(defaultRows),
		cells:         make(map[cellKey]string),
		registers:     make(map[rune]clipboard),
		marks:         make(map[rune]cellKey),
		editCursor:    editCursor,
		headerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#ADBCD9")),
		activeHeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5E5E5E")).
			Background(lipgloss.Color("#F4E1D0")).
			Bold(true),
		rowLabelStyle: lipgloss.NewStyle().
			Foreground(headerGray),
		activeRowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2D4F67")).
			Bold(true),
		activeRowCellStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2D4F67")),
		gridStyle: lipgloss.NewStyle().
			Foreground(gridGray),
		formulaCellStyle: lipgloss.NewStyle().
			Foreground(formulaGreen),
		formulaErrorStyle: lipgloss.NewStyle().
			Foreground(errorRed),
		activeCellStyle: lipgloss.NewStyle().
			//Reverse(true),
			Foreground(lipgloss.Color("#5E5E5E")).
			Background(lipgloss.Color("#F4E1D0")).
			Bold(true),
		activeFormulaStyle: lipgloss.NewStyle().
			Reverse(true).
			Foreground(formulaGreen),
		activeFormulaErrorStyle: lipgloss.NewStyle().
			Reverse(true).
			Foreground(errorRed),
		selectCellStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(white).
			Bold(true),
		selectFormulaStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(formulaGreen).
			Bold(true),
		selectFormulaErrorStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(errorRed).
			Bold(true),
		selectActiveCellStyle: lipgloss.NewStyle().
			Background(selectAccent).
			Foreground(white).
			Bold(true).
			Underline(true),
		selectActiveFormulaStyle: lipgloss.NewStyle().
			Background(selectAccent).
			Foreground(formulaGreen).
			Bold(true).
			Underline(true),
		selectActiveFormulaErrorStyle: lipgloss.NewStyle().
			Background(selectAccent).
			Foreground(errorRed).
			Bold(true).
			Underline(true),
		selectHeaderStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(white).
			Bold(true),
		selectActiveHeaderStyle: lipgloss.NewStyle().
			Background(selectAccent).
			Foreground(white).
			Bold(true),
		selectRowStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(white).
			Bold(true),
		selectBorderStyle: lipgloss.NewStyle().
			Background(selectBackground).
			Foreground(selectAccent),
		statusBarStyle: lipgloss.NewStyle().
			Background(statusGray).
			Foreground(statusText),
		statusTextStyle: lipgloss.NewStyle().
			Background(statusGray).
			Foreground(statusText),
		statusAccentStyle: lipgloss.NewStyle().
			Background(statusGray).
			Foreground(statusAccent),
		statusNormalStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("33")).
			Foreground(white).
			Padding(0, 1),
		statusInsertStyle: lipgloss.NewStyle().
			Background(insertAccent).
			Foreground(white).
			Padding(0, 1),
		statusSelectStyle: lipgloss.NewStyle().
			Background(statusSelectAccent).
			Foreground(white).
			Padding(0, 1),
		commandLineStyle: lipgloss.NewStyle().
			Foreground(statusText),
		commandErrorStyle: lipgloss.NewStyle().
			Foreground(errorRed),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !fixedCellWidth {
			m.cellWidth = m.computeCellWidth()
		}
		m.ensureVisible()
		return m, nil

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			if row, col, ok := m.cellFromMouse(msg.X, msg.Y); ok {
				if m.mode == insertMode {
					m.commitCurrentInput()
					m.mode = normalMode
					m.insertKeys = nil
					m.recordingInsert = false
					m.editCursor.Blur()
				}
				if m.mode == selectMode {
					m.mode = normalMode
				}
				m.clearNormalPrefixes()
				m.goToCell(row, col)
			}
		}
		return m, nil

	case tea.KeyMsg:
		if isQuitKey(msg) {
			return m, tea.Quit
		}
		if !m.commandPending {
			m.commandMessage = ""
			m.commandError = false
		}
		if m.mode != insertMode && m.commandPending {
			if cmd, handled := m.handlePendingCommand(msg); handled {
				return m, cmd
			}
		}
		if m.mode != insertMode && m.registerPending {
			if m.handlePendingRegister(msg) {
				return m, nil
			}
		}
		if m.mode != insertMode && m.deletePending {
			if m.handlePendingDelete(msg) {
				return m, nil
			}
		}
		if m.mode != insertMode && m.yankPending {
			if m.handlePendingYank(msg) {
				return m, nil
			}
		}
		if m.mode != insertMode && m.zPending {
			if m.handlePendingZ(msg) {
				return m, nil
			}
		}
		if m.mode != insertMode && m.gotoPending {
			if m.handlePendingGoto(msg) {
				return m, nil
			}
		}
		if m.mode != insertMode && (m.markPending || m.markJumpPending) {
			if m.handlePendingMark(msg) {
				return m, nil
			}
		}

		if m.mode == insertMode && isEscapeKey(msg) {
			if m.recordingInsert && !m.replayingChange {
				m.insertKeys = append(m.insertKeys, msg)
			}
			return m.exitInsertMode()
		}
		if m.mode == selectMode && isEscapeKey(msg) {
			m.clearNormalPrefixes()
			return m.exitSelectMode(), nil
		}

		switch m.mode {
		case normalMode:
			return m.updateNormal(msg)
		case insertMode:
			return m.updateInsert(msg)
		case selectMode:
			return m.updateSelect(msg)
		case commandMode:
			return m, nil
		}
	}

	if m.mode == insertMode || m.mode == commandMode {
		var cmd tea.Cmd
		m.editCursor, cmd = m.editCursor.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) computeCellWidth() int {
	// Layout per row: [rowLabel] [space] [│] [cell] [│] [cell] [│]
	// rowLabelWidth + 1 (space) + 1 (border) + cellWidth + 1 (border) + cellWidth + 1 (border)
	// So: available = width - rowLabelWidth - 2 - (totalCols + 1) borders
	//                                                  ^^^ = 3 borders for 2 cols
	overhead := m.rowLabelWidth + 1 + (totalCols + 1)
	available := m.width - overhead
	if available < totalCols*2 {
		return 2 // minimum sane width
	}
	return available / totalCols
}

func isQuitKey(msg tea.KeyMsg) bool {
	if msg.Type == tea.KeyCtrlC {
		return true
	}
	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == rune(3) {
		return true
	}
	return msg.String() == "ctrl+c"
}

func isEscapeKey(msg tea.KeyMsg) bool {
	if msg.Type == tea.KeyEscape {
		return true
	}
	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == rune(27) {
		return true
	}
	switch msg.String() {
	case "esc", "ctrl+[", "\x1b":
		return true
	}
	return false
}

func (m *model) pushUndoState() {
	m.undoStack = append(m.undoStack, m.snapshotUndoState())
	m.redoStack = nil
	m.dirtyFile = true
}

func (m *model) undoLastOperation() {
	if len(m.undoStack) == 0 {
		return
	}
	m.redoStack = append(m.redoStack, m.snapshotUndoState())
	last := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]
	m.restoreUndoState(last)
}

func (m *model) redoLastOperation() {
	if len(m.redoStack) == 0 {
		return
	}
	m.undoStack = append(m.undoStack, m.snapshotUndoState())
	last := m.redoStack[len(m.redoStack)-1]
	m.redoStack = m.redoStack[:len(m.redoStack)-1]
	m.restoreUndoState(last)
}

func (m model) snapshotUndoState() undoState {
	return undoState{
		cells:       cloneCells(m.cells),
		rowCount:    m.rowCount,
		selectedRow: m.selectedRow,
		selectedCol: m.selectedCol,
		selectRow:   m.selectRow,
		selectCol:   m.selectCol,
		selectRows:  m.selectRows,
		rowOffset:   m.rowOffset,
		colOffset:   m.colOffset,
	}
}

func (m *model) restoreUndoState(state undoState) {
	m.cells = cloneCells(state.cells)
	m.rowCount = max(1, state.rowCount)
	m.syncRowLabelWidth()
	m.selectedRow = state.selectedRow
	m.selectedCol = state.selectedCol
	m.selectRow = state.selectRow
	m.selectCol = state.selectCol
	m.selectRows = state.selectRows
	m.rowOffset = state.rowOffset
	m.colOffset = state.colOffset
	m.ensureVisible()
}

func (m model) cellValue(row, col int) string {
	return m.cells[cellKey{row: row, col: col}]
}

func (m *model) setCellValue(row, col int, value string) {
	key := cellKey{row: row, col: col}
	if value == "" {
		delete(m.cells, key)
		return
	}
	m.cells[key] = value
}

func (m *model) syncRowLabelWidth() {
	m.rowLabelWidth = rowLabelWidthForCount(m.rowCount)
}

func cloneCells(cells map[cellKey]string) map[cellKey]string {
	cloned := make(map[cellKey]string, len(cells))
	maps.Copy(cloned, cells)
	return cloned
}
