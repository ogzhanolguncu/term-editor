// This handles cursor positioning.
package main

import "fmt"

// CURSOR DATA STRUCTURES:
// [ ] - Cursor struct - position int, optional selectionStart/selectionEnd int
// [ ] - CursorManager struct - holds one Cursor and reference to TextBuffer
// [ ] - NewCursorManager(buffer *TextBuffer, pos int) CursorManager - constructor with buffer reference

// CURSOR INTERFACE METHODS:
// [ ] - GetPosition() int - return current text position
// [ ] - SetPosition(pos int) error - move cursor with bounds checking
// [ ] - GetCursors() []Cursor - return slice with single cursor (interface compliance)
// [ ] - MoveCursor(index int, pos int) error - move cursor (ignore index, always 0)
// [ ] - ApplyTextChange(pos int, delta int) - adjust cursor after text insertion/deletion

// CURSOR MOVEMENT OPERATIONS:
// [ ] - MoveLeft() bool - move left one position, handle line wrapping
// [ ] - MoveRight() bool - move right one position, handle line wrapping
// [ ] - MoveUp() bool - move to same column on previous line
// [ ] - MoveDown() bool - move to same column on next line
// [ ] - MoveToLineStart() bool - move to beginning of current line
// [ ] - MoveToLineEnd() bool - move to end of current line

// TEXT-AWARE NAVIGATION:
// [ ] - GetLineColumn() (int, int) - return current line and column numbers
// [ ] - MoveToLine(lineNum int) error - jump to start of specific line
// [ ] - IsAtStart() bool - check if cursor at position 0
// [ ] - IsAtEnd() bool - check if cursor at end of buffer

type Cursor struct {
	position int
}

type CursorManager struct {
	cursor *Cursor
	buffer *TextBuffer
}

//  CORE POSITIONING:
// [ ] NewCursorManager(buffer *TextBuffer, pos int) (*CursorManager, error)
// [ ] GetPosition() int
// [ ] SetPosition(pos int) error
// [ ] ApplyTextChange(changePos int, delta int)
// [ ] GetLineColumn() (int, int)

// BASIC MOVEMENT:
// [ ] MoveLeft() bool
// [ ] MoveRight() bool
// [ ] IsAtStart() bool
// [ ] IsAtEnd() bool

// LINE MOVEMENT:
// [ ] MoveUp() bool
// [ ] MoveDown() bool
// [ ] MoveToLineStart() bool
// [ ] MoveToLineEnd() bool
// [ ] MoveToLine(lineNum int) error

// PHASE 4 - INTERFACE COMPLIANCE (for future multi-cursor):
// [ ] GetCursors() []Cursor
// [ ] MoveCursor(index int, pos int) error

func NewCursorManager(buffer *TextBuffer) *CursorManager {
	return &CursorManager{
		cursor: &Cursor{position: 0},
		buffer: buffer,
	}
}

func (cm *CursorManager) GetPosition() int {
	return cm.cursor.position
}

func (cm *CursorManager) SetPosition(pos int) error {
	if pos < 0 || pos > cm.buffer.Length() {
		return fmt.Errorf("position out of bounds")
	}
	cm.cursor.position = pos
	return nil
}

// Text: "Hello World"
//        0123456789A (positions)
//
// Delete 3 characters starting at position 4
// - changePos = 4
// - delta = -3 (negative because we're deleting)
// - Deleted range: positions 4, 5, 6 ("o W")
// - Result: "Hell World"
//
//
// The logic checks where the cursor is:
//
// Case 1: Cursor after deleted range
//
// Cursor at position 8 ("r" in "World")
// - cursor.position > changePos? (8 > 4) ✓
// - cursor.position <= changePos + (-delta)? (8 <= 4 + 3 = 7) ✗
// - So cursor was AFTER deleted range
// - Shift it back: position 8 + (-3) = 5
//
//
// Case 2: Cursor inside deleted range
//
// Cursor at position 5 ("W" that got deleted)
// - cursor.position > changePos? (5 > 4) ✓
// - cursor.position <= changePos + (-delta)? (5 <= 7) ✓
// - So cursor was INSIDE deleted range
// - Clamp to deletion start: position = 4
//
//
// The math:
//
// changePos + (-delta) = end of deleted range
//
// (-delta) converts negative delta to positive (deletion size)
//
// If cursor ≤ end of deleted range, it was inside the deletion
//
// If cursor > end of deleted range, it was after the deletion

func (cm *CursorManager) ApplyTextChange(changePos int, delta int) {
	if delta > 0 {
		// Text inserted
		if cm.cursor.position >= changePos {
			cm.cursor.position += delta
		}
	} else if delta < 0 {
		// Text deleted
		if cm.cursor.position > changePos {
			if cm.cursor.position <= changePos+(-delta) {
				// Cursor was inside deleted range
				cm.cursor.position = changePos
			} else {
				// Cursor was after deleted range
				cm.cursor.position += delta // delta is negative
			}
		}
	}
}

// Text: "Hello\nWorld\nTest"
//       01234 5 67890 1 2345
//            ^cursor at pos 7
//
// Line = CharToLine(7) = 1
// LineStart = LineToChar(1) = 6
// Column = 7 - 6 = 1
// Returns (1, 1)
// func (cm *CursorManager) GetLineColumn() (int,int) {
