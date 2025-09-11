// TextBuffer is a pass-through wrapper around your gap buffer, but it adds line tracking intelligence.
package main

import (
	"slices"
	"sort"
)

// TEXTBUFFER COORDINATOR (OWNS GAP BUFFER + LINE TRACKING):
// [x] - TextBuffer struct - gap *GapBuffer, lineStarts []int fields
// [x] - NewTextBuffer(initialSize int) (*TextBuffer, error) - Initialize with gap buffer and line starts at [0]
// [x] - String() string - Return full text content
// [x] - Length() int - Return total character count
// [x] - Insert(pos int, ch rune)  - Insert single character, update line tracking
// [ ] - InsertString(pos int, text string) error - Insert text, handle multiple newlines efficiently
// [ ] - Delete(pos int) error - Delete single character, merge lines if deleting \n
// [ ] - DeleteRange(start, end int) error - Delete range, handle multiple newline removal
// [x] - CharAt(pos int) rune - Get character at position (delegate to gap buffer)
// [ ] - Substring(start, end int) string - Get text range (delegate to gap buffer)
// [ ] - Find(needle string) []int - Search text (delegate to gap buffer)
// [ ] - LoadFromString(content string) error - Initialize buffer from existing text content
//
// LINE-AWARE OPERATIONS:
// [ ] - LineCount() int - Return number of lines
// [ ] - LineToChar(lineNum int) int - Convert line number to starting char position
// [ ] - CharToLine(pos int) int - Convert char position to line number (binary search)
// [ ] - GetLine(lineNum int) string - Return content of specific line
// [ ] - LineLength(lineNum int) int - Return character count of line (excluding \n)
//
// INTERNAL LINE TRACKING (PRIVATE METHODS):
// [ ] - insertLineAt(pos int) - Add new line start to lineStarts array
// [ ] - deleteLineAt(pos int) - Remove line start from lineStarts array
// [ ] - shiftLinesAfter(pos int, delta int) - Adjust line positions after text change
// [ ] - rebuildLines() - Full line scan (for debugging/validation)
// [ ] - findLineInsertPosition(pos int) int - Binary search to find where to insert new line start
//
//
//Start implementing in this order:

// TextBuffer struct definition
// NewTextBuffer() constructor
// Basic pass-through methods (String(), Length(), CharAt())
// Insert() with line tracking - this is your core logic
// LineCount(), LineToChar(), CharToLine() - test your line tracking
// Rest of the methods

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
	tb.gBuf.MoveGapTo(pos)
	tb.gBuf.Insert(ch)

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
