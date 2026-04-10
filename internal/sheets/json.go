package sheets

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// loadJSONFile loads a flat JSON object ({key: value, ...}) from path.
// The TUI displays it as a 2-column table: KEYS | VALUES.
func (m *model) loadJSONFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("json: %w", err)
	}

	keys := sortedAnyKeys(obj)
	m.jsonKeys = keys

	m.cells = make(map[cellKey]string)
	m.rowCount = defaultRows
	m.syncRowLabelWidth()
	m.jsonShape = jsonShapeSingleObject
	m.undoStack = nil
	m.redoStack = nil
	m.promptKind = noPrompt
	m.editingValue = ""
	m.editingCursor = 0
	m.deletePending = false
	m.yankPending = false
	m.yankCount = 0
	m.zPending = false
	m.gotoPending = false
	m.gotoBuffer = ""
	m.commandPending = false
	m.commandBuffer = ""
	m.commandCursor = 0
	m.commandMessage = ""
	m.commandError = false
	m.countBuffer = ""
	m.registerPending = false
	m.activeRegister = 0
	m.searchDirection = 0
	m.markPending = false
	m.markJumpPending = false
	m.markJumpExact = false
	m.selectRows = false
	m.hasCopyBuffer = false
	m.selectedRow = 0
	m.selectedCol = 1 // start on VALUES column
	m.selectRow = 0
	m.selectCol = 1 // start first cell of VALUES column
	m.rowOffset = 0
	m.colOffset = 0
	m.jumpBack = nil
	m.jumpForward = nil
	m.lastChange = nil
	m.insertKeys = nil
	m.recordingInsert = false
	m.replayingChange = false
	m.dirtyFile = false
	m.currentFilePath = path

	for row, k := range keys {
		m.setCellValue(row, 0, k)
		m.setCellValue(row, 1, anyToString(obj[k]))
	}

	if len(keys) > m.rowCount {
		m.rowCount = len(keys)
		m.syncRowLabelWidth()
	}

	return nil
}

// writeJSONFile saves the 2-column TUI data back as a flat JSON object.
func (m model) writeJSONFile(path string) error {
	// Collect rows: col0 = key, col1 = value.
	orderedKeys := make([]string, 0)
	obj := make(map[string]any)

	for row := 0; ; row++ {
		k := m.cellValue(row, 0)
		if k == "" {
			// Check whether any further rows exist; if not, stop.
			hasMore := false
			for r := row + 1; r < m.rowCount; r++ {
				if m.cellValue(r, 0) != "" {
					hasMore = true
					break
				}
			}
			if !hasMore {
				break
			}
			continue
		}
		v := m.cellValue(row, 1)
		obj[k] = stringToJSONValue(v)
		orderedKeys = append(orderedKeys, k)
	}

	data, err := marshalOrdered(orderedKeys, obj)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

// marshalOrdered writes a JSON object preserving key order.
func marshalOrdered(keys []string, obj map[string]any) ([]byte, error) {
	// Build JSON manually to preserve insertion order.
	buf := []byte("{\n")
	for i, k := range keys {
		v, ok := obj[k]
		if !ok {
			continue
		}
		vBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		kBytes, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf = append(buf, "  "...)
		buf = append(buf, kBytes...)
		buf = append(buf, ": "...)
		buf = append(buf, vBytes...)
		if i < len(keys)-1 {
			buf = append(buf, ',')
		}
		buf = append(buf, '\n')
	}
	buf = append(buf, '}')
	return buf, nil
}

// stringToJSONValue restores native JSON types from a cell string.
func stringToJSONValue(s string) any {
	if s == "" {
		return nil
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err == nil {
		switch v.(type) {
		case bool, float64, nil:
			return v
		case map[string]any, []any:
			return v
		}
	}
	return s
}

// anyToString converts a JSON value to a display string.
func anyToString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		b, _ := json.Marshal(val)
		return string(b)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

// sortedAnyKeys returns sorted keys of a map[string]any.
func sortedAnyKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
