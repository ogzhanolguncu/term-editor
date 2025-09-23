// TextBuffer is a pass-through wrapper around your gap buffer, but it adds line tracking intelligence.
package main

import (
	"slices"
	"sort"
	"strings"
)

// TEXTBUFFER COORDINATOR (OWNS GAP BUFFER + LINE TRACKING):
// [ ] - Delete(pos int) error - Delete single character, merge lines if deleting \n
// [ ] - DeleteRange(start, end int) error - Delete range, handle multiple newline removal

// [ ] - LoadFromString(content string) error - Initialize buffer from existing text content
//
// INTERNAL LINE TRACKING (PRIVATE METHODS):
// [ ] - insertLineAt(pos int) - Add new line start to lineStarts array
// [ ] - deleteLineAt(pos int) - Remove line start from lineStarts array
// [ ] - shiftLinesAfter(pos int, delta int) - Adjust line positions after text change
// [ ] - rebuildLines() - Full line scan (for debugging/validation)
// [ ] - findLineInsertPosition(pos int) int - Binary search to find where to insert new line start
//

type TextBuffer struct {
	gBuf       *GapBuffer
	lineStarts []int
}

func NewTextBuffer(initialSize int) (*TextBuffer, error) {
	gapBuffer, err := NewGapBuffer(initialSize)
	if err != nil {
		return nil, err
	}

	return &TextBuffer{
		gBuf:       gapBuffer,
		lineStarts: []int{0},
	}, nil
}

func (tb *TextBuffer) String() string {
	return tb.gBuf.String()
}

func (tb *TextBuffer) Length() int {
	return tb.gBuf.Length()
}

func (tb *TextBuffer) CharAt(pos int) rune {
	return tb.gBuf.CharAt(pos)
}

func (tb *TextBuffer) Insert(pos int, ch rune) {
	tb.gBuf.InsertAt(pos, ch)

	if ch == '\n' {
		insertPos := sort.SearchInts(tb.lineStarts, pos+1)
		tb.lineStarts = slices.Insert(tb.lineStarts, insertPos, pos+1)

		for i := insertPos + 1; i < len(tb.lineStarts); i++ {
			tb.lineStarts[i]++
		}
	} else {
		for i, lineStart := range tb.lineStarts {
			if pos < lineStart {
				tb.lineStarts[i]++
			}
		}
	}
}

func (tb *TextBuffer) InsertString(pos int, text string) {
	if pos >= tb.Length() {
		// Fast path, inserting at end, no existing lines to shift
		tb.gBuf.InsertStringAt(pos, text)
		for i, ch := range text {
			if ch == '\n' {
				tb.lineStarts = append(tb.lineStarts, pos+i+1)
			}
		}
	} else {
		// Slow path, inserting in middle, need to shift existing lines
		tb.gBuf.InsertStringAt(pos, text)
		textLen := len(text)

		// Shift all existing line starts that come after insertion point
		for i := range tb.lineStarts {
			if tb.lineStarts[i] > pos {
				tb.lineStarts[i] += textLen
			}
		}

		// Add new line starts from inserted text
		for i, ch := range text {
			if ch == '\n' {
				newLinePos := pos + i + 1
				insertPos := sort.SearchInts(tb.lineStarts, newLinePos)
				tb.lineStarts = slices.Insert(tb.lineStarts, insertPos, newLinePos)
			}
		}
	}
}

func (tb *TextBuffer) LineCount() int {
	return len(tb.lineStarts)
}

func (tb *TextBuffer) LineToChar(lineNum int) int {
	if lineNum <= 0 {
		return 0
	}

	if lineNum >= len(tb.lineStarts) {
		return tb.lineStarts[len(tb.lineStarts)-1]
	}
	return tb.lineStarts[lineNum]
}

func (tb *TextBuffer) CharToLine(pos int) int {
	if pos <= 0 {
		return 0
	}
	if pos >= tb.Length() {
		return len(tb.lineStarts) - 1
	}

	// Find the rightmost line start <= pos
	line := sort.Search(len(tb.lineStarts), func(i int) bool {
		return tb.lineStarts[i] > pos
	}) - 1

	return line
}

func (tb *TextBuffer) Find(needle string) []int {
	return tb.gBuf.Find(needle)
}

func (tb *TextBuffer) Substring(start, end int) string {
	return tb.gBuf.Substring(start, end)
}

// Line returns the string content of specific line number. Newlines are included
func (tb *TextBuffer) Line(lineNum int) string {
	if lineNum < 0 || lineNum >= len(tb.lineStarts) {
		return ""
	}

	start := tb.lineStarts[lineNum]
	if lineNum == len(tb.lineStarts)-1 {
		return tb.Substring(start, tb.Length())
	}

	end := tb.lineStarts[lineNum+1]
	return tb.Substring(start, end)
}

// LineLength returns the length of the given line. Excludes newlines
func (tb *TextBuffer) LineLength(lineNum int) int {
	return len(strings.TrimSpace(tb.Line(lineNum)))
}

func (tb *TextBuffer) Delete(pos int) {
	if pos < 0 || pos >= tb.Length() {
		return
	}

	isNewLine := tb.CharAt(pos) == '\n'
	tb.gBuf.MoveGapTo(pos)
	tb.gBuf.Delete()

	if isNewLine {
		for i, start := range tb.lineStarts {
			if start == pos+1 {
				tb.lineStarts = append(tb.lineStarts[:i], tb.lineStarts[i+1:]...)
				break
			}
		}
	}

	for i := range len(tb.lineStarts) {
		if tb.lineStarts[i] > pos {
			tb.lineStarts[i]--
		}
	}
}
