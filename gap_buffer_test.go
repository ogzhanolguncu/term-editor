package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWords(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)

	require.True(t, gbuf.gapEnd == 30)
	require.True(t, gbuf.gapStart == 0)
	require.Equal(t, gbuf.Length(), 0)
	require.Equal(t, gbuf.String(), "")
}

func TestInsert(t *testing.T) {
	gbuf, err := NewGapBuffer(6)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')

	require.Equal(t, gbuf.String(), "Hel")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 6)
}

func TestExpand(t *testing.T) {
	gbuf, err := NewGapBuffer(3)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')

	require.Equal(t, gbuf.String(), "Hell")
	require.Equal(t, gbuf.Length(), 4)
	require.Equal(t, gbuf.gapStart, 4)
	require.Equal(t, gbuf.gapEnd, 6)
	require.Equal(t, len(gbuf.buffer), 6)
}

func TestCursorMove(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(2)
	require.Equal(t, gbuf.String(), "Hello")
	require.Equal(t, gbuf.gapStart, 2)
	require.Equal(t, gbuf.gapEnd, 7)

	gbuf.Insert('X')
	require.Equal(t, gbuf.String(), "HeXllo")

	gbuf.MoveGapTo(6)

	gbuf.Insert('T')
	require.Equal(t, gbuf.String(), "HeXlloT")
	require.Equal(t, gbuf.gapStart, 7)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(10)
	gbuf.MoveGapTo(8)

	gbuf.Insert('X')
	gbuf.Insert('Y')
	gbuf.Insert('Z')
	require.Equal(t, gbuf.gapStart, 10)
	require.Equal(t, gbuf.gapEnd, 10)
}

func TestBackspace(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.Backspace()

	require.Equal(t, gbuf.gapStart, 4)
	require.Equal(t, gbuf.gapEnd, 10)

	require.Equal(t, gbuf.String(), "Hell")
}

func TestDelete(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(0)

	require.Equal(t, gbuf.gapStart, 0)
	require.Equal(t, gbuf.gapEnd, 5)

	gbuf.Delete()

	require.Equal(t, gbuf.String(), "ello")
	require.Equal(t, gbuf.gapStart, 0)
	require.Equal(t, gbuf.gapEnd, 6)

	gbuf.MoveGapTo(3)

	gbuf.Delete()

	require.Equal(t, gbuf.String(), "ell")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.Delete()
	require.Equal(t, gbuf.String(), "ell")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 10)
}

func TestGapSize(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.GapSize(), 5)
}

func TestGapPos(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	gbuf.MoveGapTo(3)
	require.Equal(t, gbuf.GapPos(), 3)
}

func TestCharAt(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.CharAt(4), 'o')
	require.Equal(t, gbuf.CharAt(11), rune(0))
}

func TestInsertString(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')
	gbuf.InsertString(" World!")

	require.Equal(t, gbuf.String(), "Hello World!")
}

func TestDeleteRange(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")

	// Delete "o Wo" (positions 4-7, exclusive end)
	gbuf.DeleteRange(4, 8)
	// Should result in "Hellrld!"
	require.Equal(t, "Hellrld!", gbuf.String())
	require.Equal(t, 8, gbuf.Length()) // Original 12 chars - 4 deleted chars
}

func TestDeleteRangeInReverse(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")

	// Delete "o Wo" (positions 4-7, exclusive end)
	gbuf.DeleteRange(8, 4)
	// Should result in "Hellrld!"
	require.Equal(t, "Hellrld!", gbuf.String())
	require.Equal(t, 8, gbuf.Length()) // Original 12 chars - 4 deleted chars
}

func TestDeleteRangeAtBeginning(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")
	// Delete "Hell" (positions 0-3, exclusive end)
	gbuf.DeleteRange(0, 4)
	// Should result in "o World!"
	require.Equal(t, "o World!", gbuf.String())
	require.Equal(t, 8, gbuf.Length()) // Original 12 chars - 4 deleted chars
}

func TestDeleteRangeAtEnd(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")
	// Delete "rld!" (positions 8-11, exclusive end)
	gbuf.DeleteRange(8, 12)
	// Should result in "Hello Wo"
	require.Equal(t, "Hello Wo", gbuf.String())
	require.Equal(t, 8, gbuf.Length()) // Original 12 chars - 4 deleted chars
}

func TestDeleteRangeEntireBuffer(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")
	// Delete everything (positions 0-11, exclusive end)
	gbuf.DeleteRange(0, 12)
	// Should result in empty string
	require.Equal(t, "", gbuf.String())
	require.Equal(t, 0, gbuf.Length())
}

func TestDeleteRangeEmptyRange(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")
	// Delete nothing (same start and end position)
	gbuf.DeleteRange(5, 5)
	// Should remain unchanged
	require.Equal(t, "Hello World!", gbuf.String())
	require.Equal(t, 12, gbuf.Length())
}

func TestDeleteRangeOutOfBounds(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	gbuf.InsertString("Hello World!")
	gbuf.DeleteRange(5, 20)
	require.Equal(t, "Hello", gbuf.String())
	require.Equal(t, 5, gbuf.Length())
}

func TestDeleteRangeOnEmptyBuffer(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)
	// Try to delete from empty buffer - should be a noop
	gbuf.DeleteRange(0, 1)
	// Should remain empty
	require.Equal(t, "", gbuf.String())
	require.Equal(t, 0, gbuf.Length())
}

func TestSubstring(t *testing.T) {
	// Test normal case
	gb, _ := NewGapBuffer(20)
	gb.InsertString("hello world")
	result := gb.Substring(3, 8)
	expected := "lo wo"
	require.Equal(t, expected, result, "Normal case")

	// Test with gap in middle of substring range
	gb.MoveGapTo(6) // Gap between "hello " and "world"
	result = gb.Substring(3, 8)
	require.Equal(t, expected, result, "Gap in middle")

	// Test empty string
	emptyGB, _ := NewGapBuffer(10)
	result = emptyGB.Substring(0, 0)
	require.Equal(t, "", result, "Empty string")

	// Test start equals end
	result = gb.Substring(2, 2)
	require.Equal(t, "", result, "Start equals end")

	// Test start greater than end
	result = gb.Substring(5, 2)
	require.Equal(t, "", result, "Start > end")

	// Test negative start (should clamp to 0)
	result = gb.Substring(-3, 4)
	expected = "hell"
	require.Equal(t, expected, result, "Negative start")

	// Test end beyond text length (should clamp to length)
	result = gb.Substring(6, 100)
	expected = "world"
	require.Equal(t, expected, result, "End beyond length")

	// Test both bounds out of range
	result = gb.Substring(-5, 100)
	expected = "hello world"
	require.Equal(t, expected, result, "Both bounds out of range")

	// Test start beyond text length
	result = gb.Substring(50, 60)
	require.Equal(t, "", result, "Start beyond length")

	// Test single character
	result = gb.Substring(0, 1)
	expected = "h"
	require.Equal(t, expected, result, "Single char")

	// Test entire string
	result = gb.Substring(0, gb.Length())
	expected = "hello world"
	require.Equal(t, expected, result, "Entire string")
}

func TestFind(t *testing.T) {
	gb, _ := NewGapBuffer(30)
	gb.InsertString("hello world hello")

	// Normal case - multiple matches
	result := gb.Find("hello")
	expected := []int{0, 12}
	require.Equal(t, expected, result, "Multiple matches")
	//
	// // Single match
	// result = gb.Find("world")
	// expected = []int{6}
	// require.Equal(t, expected, result, "Single match")
	//
	// // No match
	// result = gb.Find("xyz")
	// expected = []int{}
	// require.Equal(t, expected, result, "No match")
	//
	// // Empty needle
	// result = gb.Find("")
	// expected = []int{}
	// require.Equal(t, expected, result, "Empty needle")
	//
	// // Needle longer than text
	// result = gb.Find("this is way too long")
	// expected = []int{}
	// require.Equal(t, expected, result, "Long needle")
	//
	// // Test with gap in different position
	// gb.MoveGapTo(8) // Gap between "world" and " hello"
	// result = gb.Find("hello")
	// expected = []int{0, 12}
	// require.Equal(t, expected, result, "Gap moved")
}
