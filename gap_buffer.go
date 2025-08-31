package main

import (
	"fmt"
)

type GapBuffer struct {
	buffer   []rune
	gapStart int
	gapEnd   int
}

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

// expandBuffer first decides the amount of gap that the buffer needs, then calls resize.
func (gb *GapBuffer) expandBuffer() {
	newSize := gb.calculateGrowSize(0)
	gb.resizeBuffer(newSize)
}

// calculateGrowSize handles growth size. If current size is relatively small, we just double the buffer,
// but if it's aggressive growth like loading a file to buffer, we make it really big. If it's not aggressive, we make it 5% bigger.
func (gb *GapBuffer) calculateGrowSize(insertSize int) int {
	textSize := gb.Length()
	currentSize := len(gb.buffer)

	switch {
	case currentSize < 512:
		return currentSize * 2
	case insertSize > gb.GapSize():
		// Large insertion, grow aggressively
		targetGap := max(textSize/10, insertSize*2)
		maxGap := max(8192, textSize/10)
		return textSize + min(targetGap, maxGap)
	default:
		// 5% gap ratio
		return textSize + max(textSize/20, 128)
	}
}

// resizeBuffer after we determine the amount of gap we need for expansion, we move the existing left/right parts to the appropriate locations.
func (gb *GapBuffer) resizeBuffer(newSize int) {
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
	if start > gb.Length() || start < 0 || start >= end {
		return ""
	}
	// Clamp end to buffer length
	if end > gb.Length() { // Fixed: was >=, should be >
		end = gb.Length()
	}
	// Substring is entirely on the left side of gap
	if end <= gb.gapStart {
		return string(gb.buffer[start:end])
	}
	// Substring is entirely on the right side of gap
	if start >= gb.gapStart {
		// Skip the gap offset since users see logical text positions, not buffer positions.
		// E.g., if buffer length is 15, actual text is 12, and gap moved to position 6,
		// then buffer[6:9] would target gap space instead of text, so we add gap offset.
		gapOffset := gb.GapSize()
		return string(gb.buffer[start+gapOffset : end+gapOffset])
	}
	// Substring spans the gap - need to copy from both sides
	leftLen := gb.gapStart - start
	rightLen := end - gb.gapStart
	result := make([]rune, leftLen+rightLen)
	// Copy text from left side of buffer
	copy(result[:leftLen], gb.buffer[start:gb.gapStart])
	// Copy text from right side of buffer
	copy(result[leftLen:], gb.buffer[gb.gapEnd:gb.gapEnd+rightLen])

	return string(result)
}

// Find moves through text char by char then tries to do full text check at every stop.
// E.g
// Text: "Hello Hello"
// Search: "Hello"
// Lands on "H" then start moving as much as needle length if every char matches it means starting position is a match
func (gb *GapBuffer) Find(needle string) []int {
	// If search is empty, bail
	if len(needle) == 0 {
		return []int{}
	}

	positions := make([]int, 0)
	needleLen := len(needle)
	textLen := gb.Length()

	// If search is bigger than the actual text, bail
	if needleLen > textLen {
		return positions
	}

	for start := 0; start <= textLen-needleLen; start++ {
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
