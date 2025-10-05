package gapbuffer

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

	// require.Equal(t, gbuf.CharAt(4), 'o')
	// require.Equal(t, gbuf.CharAt(11), rune(0))
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
	gb, _ := NewGapBuffer(20)
	gb.InsertString("hello world")

	// Core functionality
	require.Equal(t, "lo wo", gb.Substring(3, 8), "Normal case")

	// Gap spanning (the critical edge case)
	gb.MoveGapTo(6)
	require.Equal(t, "lo wo", gb.Substring(3, 8), "Gap in middle")

	// Bounds handling
	require.Equal(t, "", gb.Substring(5, 2), "Start > end")
	require.Equal(t, "world", gb.Substring(6, 100), "End beyond length")
	require.Equal(t, "", gb.Substring(50, 60), "Start beyond length")
}

func TestFind(t *testing.T) {
	gb, _ := NewGapBuffer(30)
	gb.InsertString("hello world hello")

	// Core functionality
	require.Equal(t, []int{0, 12}, gb.Find("hello"), "Multiple matches")
	require.Equal(t, []int{6}, gb.Find("world"), "Single match")
	require.Equal(t, []int{}, gb.Find("xyz"), "No match")

	// Gap spanning (critical case)
	gb.MoveGapTo(8) // Gap splits "world hello"
	require.Equal(t, []int{0, 12}, gb.Find("hello"), "Gap moved")
	require.Equal(t, []int{4}, gb.Find("o w"), "Match spans gap")

	// Edge case
	require.Equal(t, []int{}, gb.Find(""), "Empty needle")
}
