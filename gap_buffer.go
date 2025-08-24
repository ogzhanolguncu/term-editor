package main

type GapBuffer struct {
	buffer   []rune
	gapStart int
	gapEnd   int
}

func NewGapBuffer(initialSize int) *GapBuffer {
	return &GapBuffer{
		buffer:   make([]rune, initialSize),
		gapStart: 0,
		gapEnd:   initialSize,
	}
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

// MoveCursorTo add boundary checks
func (gb *GapBuffer) MoveCursorTo(pos int) {
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

func (gb *GapBuffer) Delete() {}
