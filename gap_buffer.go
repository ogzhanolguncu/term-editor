package main

import "fmt"

type GapBuffer struct {
	buffer   []rune
	gapStart int
	gapEnd   int
}

// TODO:
// Insert string - InsertString(s string) for pasting or inserting multiple characters efficiently
// Get character at position - CharAt(pos int) rune for syntax highlighting, search, etc.
// Delete range - DeleteRange(start, end int) for selecting and deleting blocks of text
// Get substring - Substring(start, end int) string for copying selected text
// Find/search - Find(needle string) []int to locate text patterns
// Cursor position tracking - Add a CursorPos() int method to know where you are
// Better error handling - Maybe add those error-returning variants we discussed
// Gap info - GapSize() int, GapPosition() int for debugging or stats
// Shrink buffer - When gap gets too large, compact it

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

func (gb *GapBuffer) ToString() string {
	leftPart := string(gb.buffer[0:gb.gapStart])
	rightPart := string(gb.buffer[gb.gapEnd:])
	return leftPart + rightPart
}

func (gb *GapBuffer) Length() int {
	return len(gb.buffer) - (gb.gapEnd - gb.gapStart)
}

func (gb *GapBuffer) Insert(ch rune) {
	if gb.gapEnd-gb.gapStart == 0 {
		gb.expandBuffer()
	}
	gb.buffer[gb.gapStart] = ch
	gb.gapStart++
}

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

// MoveGapTo add boundary checks
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

// Backspace we just have to move gapStart to left
func (gb *GapBuffer) Backspace() {
	if gb.gapStart == 0 {
		return
	}
	gb.gapStart--
}

func (gb *GapBuffer) Delete() {
	if gb.gapStart == gb.Length() {
		return
	}
	gb.gapEnd++
}
