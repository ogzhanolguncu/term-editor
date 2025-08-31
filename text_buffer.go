// TextBuffer is a pass-through wrapper around your gap buffer, but it adds line tracking intelligence.
package main

// TEXTBUFFER COORDINATOR (OWNS GAP BUFFER + LINE TRACKING):
// [ ] - TextBuffer struct - gap *GapBuffer, lineStarts []int fields
// [ ] - NewTextBuffer(initialSize int) (*TextBuffer, error) - Initialize with gap buffer and line starts at [0]
// [ ] - String() string - Return full text content
// [ ] - Length() int - Return total character count
// [ ] - Insert(pos int, ch rune) error - Insert single character, update line tracking
// [ ] - InsertString(pos int, text string) error - Insert text, handle multiple newlines efficiently
// [ ] - Delete(pos int) error - Delete single character, merge lines if deleting \n
// [ ] - DeleteRange(start, end int) error - Delete range, handle multiple newline removal
// [ ] - CharAt(pos int) rune - Get character at position (delegate to gap buffer)
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
//
// The private helper methods will emerge naturally as you implement the public ones. You'll know exactly what insertLineAt() and shiftLinesAfter() need to do once you're writing Insert().
