package main

import (
	"fmt"
)

type GapBuffer struct {
	buffer   []rune
	gapStart int
	gapEnd   int
}

// TODO:
// [x] - Insert string - InsertString(s string) for pasting or inserting multiple characters efficiently
// [x] - Get character at position - CharAt(pos int) rune for syntax highlighting, search, etc.
// [x] - Delete range - DeleteRange(start, end int) for selecting and deleting blocks of text
// [x] - Get substring - Substring(start, end int) string for copying selected text
// [x] - Find/search - Find(needle string) []int to locate text patterns
// [x] - Gap info - GapSize() int, GapPosition() int for debugging or stats
// [ ] - Shrink buffer - When gap gets too large, compact it
// [ ] - Smart expandBuffer - Use adaptive gap sizing instead of always doubling
// [ ] - UTF-8 safety - Ensure gap movement doesn't corrupt multi-byte sequences (low risk with []rune)
// [ ] - Handle grapheme clusters - Current []rune approach splits composed characters like 👨‍👩‍👧‍👦
//
// SEPARATE LINE HANDLING (DON'T PUT IN GAP BUFFER):
// [ ] - LineIndex struct - Separate component that tracks line boundaries
// [ ] - LineIndex.Update() - Watches gap buffer changes and updates line positions
// [ ] - LineIndex.CharToLine() - Convert character position to line number
// [ ] - LineIndex.LineToChar() - Convert line number to character position
// [ ] - LineIndex.LineCount() - Total number of lines in buffer
//
// MULTI-CURSOR SUPPORT (SEPARATE FROM GAP BUFFER):
// [ ] - CursorManager struct - Tracks multiple cursor positions independently
// [ ] - CursorManager.Update() - Adjusts all cursors when gap buffer changes
// [ ] - CursorManager.Insert() - Apply same edit at all cursor positions
// [ ] - CursorManager.AddCursor() - Add new cursor at position
// [ ] - Selection support - Each cursor can have associated selection range

func NewGapBuffer(initialSize int) (*GapBuffer, error) {
	if initialSize <= 0 {
		return nil, fmt.Errorf("initialSize must be positive, got %d", initialSize)
	}
	return &GapBuffer{
		buffer:   make([]rune, initialSize),
		gapStart: 0,
		gapEnd:   initialSize,
	}, nil
}

// String returns the actual text, skipping the gap.
func (gb *GapBuffer) String() string {
	leftPart := string(gb.buffer[0:gb.gapStart])
	rightPart := string(gb.buffer[gb.gapEnd:])
	return leftPart + rightPart
}

func (gb *GapBuffer) Length() int {
	return len(gb.buffer) - (gb.gapEnd - gb.gapStart)
}

// Insert adds a char at the current position, expanding buffer if gap is full.
func (gb *GapBuffer) Insert(ch rune) {
	if gb.gapEnd-gb.gapStart == 0 {
		gb.expandBuffer()
	}
	gb.buffer[gb.gapStart] = ch
	gb.gapStart++
}

// expandBuffer doubles the buffer size while preserving the gap structure.
// Creates a new gap in the middle by copying text before/after the old gap
// to the start/end of the new buffer.
func (gb *GapBuffer) expandBuffer() {
	newSize := len(gb.buffer) * 2

	// Copy left part, create new gap, copy right part
	leftPart := gb.buffer[0:gb.gapStart]
	rightPart := gb.buffer[gb.gapEnd:]

	newBuffer := make([]rune, newSize)
	copy(newBuffer, leftPart)                           // Copy left part to start
	copy(newBuffer[newSize-len(rightPart):], rightPart) // Copy right part to end

	gb.buffer = newBuffer
	gb.gapStart = len(leftPart)          // Gap starts after left part
	gb.gapEnd = newSize - len(rightPart) // Gap ends before right part
}

// MoveGapTo moves gap around by figuring out which direction to go first.
// Then, shifts existing one by one till we hit the desired gap.
// Idea is we swap x+1 with y+1, x is gapStart and y is gapEnd or x-1 with y-1 depending on the direction.
func (gb *GapBuffer) MoveGapTo(pos int) {
	if pos > gb.Length() {
		return
	}
	// Don't allow negative
	if pos < 0 {
		return
	}
	// No need to move we are already there
	if pos == gb.gapStart {
		return
	}

	direction := "left"
	if pos > gb.gapStart {
		direction = "right"
	}

	if direction == "right" {
		diff := pos - gb.gapStart
		// There has to be enough space to move left
		if gb.gapEnd+diff > len(gb.buffer) {
			return
		}
		for i := range diff {
			gb.buffer[gb.gapStart+i] = gb.buffer[gb.gapEnd+i]
		}

		gb.gapStart += diff
		gb.gapEnd += diff

	}

	if direction == "left" {
		diff := gb.gapStart - pos
		for i := range diff {
			gb.buffer[gb.gapEnd-1-i] = gb.buffer[gb.gapStart-1-i]
		}

		gb.gapStart -= diff
		gb.gapEnd -= diff
	}
}

// Backspace removes the char before the cursor. No-op at start of buffer.
func (gb *GapBuffer) Backspace() {
	if gb.gapStart == 0 {
		return
	}
	gb.gapStart--
}

// Delete removes the char at the cursor. No-op at end of buffer.
func (gb *GapBuffer) Delete() {
	if gb.gapStart == gb.Length() {
		return
	}
	gb.gapEnd++
}

// GapSize returns the current gap size for debugging and testing.
func (gb *GapBuffer) GapSize() int {
	return gb.gapEnd - gb.gapStart
}

// GapPos returns where the gap (cursor) is positioned in the logical text.
func (gb *GapBuffer) GapPos() int {
	return gb.gapStart
}

func (gb *GapBuffer) CharAt(pos int) rune {
	if pos < 0 || pos >= gb.Length() {
		return rune(0)
	}
	if pos < gb.gapStart {
		// Character is in left part of buffer (before gap)
		return gb.buffer[pos]
	} else {
		// Character is in right part, adjust index to skip over gap
		return gb.buffer[pos+gb.GapSize()]
	}
}

// InsertString adds a string to the buffer. Calls Insert under the hood.
func (gb *GapBuffer) InsertString(text string) {
	for _, v := range text {
		gb.Insert(v)
	}
}

func (gb *GapBuffer) DeleteRange(start, end int) {
	if gb.Length() == 0 || start == end {
		return
	}

	startingPoint := max(0, min(start, end))
	endPoint := min(gb.Length(), max(start, end))

	if startingPoint >= endPoint {
		return
	}

	diff := endPoint - startingPoint
	gb.MoveGapTo(startingPoint)
	gb.gapEnd = gb.gapEnd + diff
}

func (gb *GapBuffer) Substring(start, end int) string {
	if start > gb.Length() {
		return ""
	}
	if start < 0 {
		start = 0
	}

	if start >= end {
		return ""
	}

	// Clamp to end
	if end > gb.Length() {
		end = gb.Length()
	}

	// Text is on the left side of gap
	if end <= gb.gapStart {
		return string(gb.buffer[start:end])
	} else if start >= gb.gapStart {
		// Text is on the right side of gap
		gapAddedStart := start + gb.GapSize()
		gapAddedEnd := end + gb.GapSize()
		return string(gb.buffer[gapAddedStart:gapAddedEnd])
	} else {
		result := make([]rune, 0, end-start)
		for i := start; i < end; i++ {
			result = append(result, gb.CharAt(i))
		}
		return string(result)
	}
}

func (gb *GapBuffer) Find(needle string) []int {
	if len(needle) == 0 {
		return []int{}
	}

	positions := make([]int, 0)
	needleLen := len(needle)
	bufferLen := gb.Length()

	if needleLen > bufferLen {
		return positions
	}

	for start := 0; start <= bufferLen-needleLen; start++ {
		match := true
		for i := range needleLen {
			if gb.CharAt(start+i) != rune(needle[i]) {
				match = false
				break
			}
		}
		if match {
			positions = append(positions, start)
		}
	}

	return positions
}
