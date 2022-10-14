package table

import "strings"

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Generate generates a markdown table from the given data and alignments.
// The first slice of data is used as header row.
// If the data has length zero, an empty string is returned.
// The number of columns is the length of the header line.
// Empty string or left alignment are used if data is missing.
// Extra data is ignored.
func Generate(data [][]string, align []Align) string {
	if len(data) == 0 {
		return ""
	}
	n := len(data[0])
	align = sanitizeRow(align, n)
	data = sanitizeData(data, n)
	sizes := cellSizes(data, align, n)
	sb := &strings.Builder{}

	printRow(sb, data[0], sizes, align) // print header line
	printSeparators(sb, sizes, align)   // print separator line
	for _, row := range data[1:] {
		printRow(sb, row, sizes, align) // print data lines
	}
	return sb.String()
}

func sanitizeData(data [][]string, n int) [][]string {
	result := make([][]string, len(data))

	for i, row := range data {
		result[i] = sanitizeRow(row, n)
	}
	return result
}

func sanitizeRow[T any](row []T, n int) []T {
	if len(row) == n {
		return row
	}
	if len(row) > n {
		return row[:n]
	}
	result := make([]T, n)
	copy(result, row)
	return result
}

func cellSizes(data [][]string, align []Align, n int) []int {
	sizes := make([]int, n)

	for _, row := range data {
		for j, cell := range row {
			sizes[j] = max(len(cell), sizes[j])
		}
	}
	for i, size := range sizes {
		if size < 1 {
			if align[i] == AlignCenter {
				sizes[i] = 1
			}
		}
	}
	return sizes
}

func printRow(sb *strings.Builder, row []string, sizes []int, align []Align) {
	sb.WriteRune('|')
	for i, cell := range row {
		printCell(sb, cell, sizes[i], align[i])
	}
	sb.WriteRune('\n')
}

func printCell(sb *strings.Builder, cell string, size int, align Align) {
	padding := size - len(cell)
	sb.WriteRune(' ')
	switch align {
	case AlignLeft:
		sb.WriteString(cell)
		printRunes(sb, ' ', padding)
	case AlignCenter:
		padLeft := padding / 2
		printRunes(sb, ' ', padLeft)
		sb.WriteString(cell)
		printRunes(sb, ' ', padding-padLeft)
	case AlignRight:
		printRunes(sb, ' ', padding)
		sb.WriteString(cell)
	}
	sb.WriteString(" |")
}

func printSeparators(sb *strings.Builder, sizes []int, align []Align) {
	sb.WriteRune('|')
	for i, size := range sizes {
		switch align[i] {
		case AlignLeft:
			sb.WriteRune(':')
			printRunes(sb, '-', size+1)
		case AlignCenter:
			sb.WriteRune(':')
			printRunes(sb, '-', size)
			sb.WriteRune(':')
		case AlignRight:
			printRunes(sb, '-', size+1)
			sb.WriteRune(':')
		}
		sb.WriteRune('|')
	}
	sb.WriteRune('\n')
}

func printRunes(sb *strings.Builder, r rune, n int) {
	for i := 0; i < n; i++ {
		sb.WriteRune(r)
	}
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
