package sheets

import (
	"strconv"
	"strings"
)

type cellRange struct {
	start cellKey
	end   cellKey
}

func (r cellRange) bounds() (top, bottom, left, right int) {
	return normalizeCellRange(r.start, r.end)
}

func (r cellRange) width() int {
	_, _, left, right := r.bounds()
	return right - left + 1
}

func (r cellRange) height() int {
	top, bottom, _, _ := r.bounds()
	return bottom - top + 1
}

func (r cellRange) isSingleCell() bool {
	return r.start == r.end
}

func parseCellRef(ref string) (cellKey, bool) {
	if ref == "" {
		return cellKey{}, false
	}

	ref = strings.ToUpper(strings.TrimSpace(ref))
	split := 0
	for split < len(ref) && isLetter(ref[split]) {
		split++
	}
	if split == 0 || split == len(ref) {
		return cellKey{}, false
	}

	columnPart := ref[:split]
	rowPart := ref[split:]
	for i := range rowPart {
		if !isDigit(rowPart[i]) {
			return cellKey{}, false
		}
	}

	row, err := strconv.Atoi(rowPart)
	if err != nil || row < 1 || row > maxRows {
		return cellKey{}, false
	}

	col := 0
	for i := range columnPart {
		col = col*26 + int(columnPart[i]-'A'+1)
	}
	col--
	if col < 0 || col >= totalCols {
		return cellKey{}, false
	}

	return cellKey{row: row - 1, col: col}, true
}

func parseCellRangeRef(ref string) (cellRange, bool) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return cellRange{}, false
	}

	parts := strings.Split(ref, ":")
	if len(parts) == 1 {
		cell, ok := parseCellRef(parts[0])
		if !ok {
			return cellRange{}, false
		}
		return cellRange{start: cell, end: cell}, true
	}
	if len(parts) != 2 {
		return cellRange{}, false
	}

	start, ok := parseCellRef(parts[0])
	if !ok {
		return cellRange{}, false
	}
	end, ok := parseCellRef(parts[1])
	if !ok {
		return cellRange{}, false
	}

	return cellRange{start: start, end: end}, true
}

func parseColumnRef(ref string) (int, bool) {
	cell, ok := parseCellRef(strings.TrimSpace(ref) + "1")
	if !ok {
		return 0, false
	}
	return cell.col, true
}

func isCellRefCommandPrefix(input string) bool {
	if input == "" {
		return true
	}

	input = strings.ToUpper(strings.TrimSpace(input))
	split := 0
	for split < len(input) && isLetter(input[split]) {
		split++
	}
	if split == 0 {
		return false
	}

	for i := split; i < len(input); i++ {
		if !isDigit(input[i]) {
			return false
		}
	}

	return len(input)-split <= len(strconv.Itoa(maxRows))
}

func rowLabelWidthForCount(rowCount int) int {
	return max(4, len(strconv.Itoa(max(1, rowCount))))
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func cellRef(row, col int) string {
	return columnLabel(col) + strconv.Itoa(row+1)
}

func clamp(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
