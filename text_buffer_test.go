package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewlineBugExposed(t *testing.T) {
	tb, err := NewTextBuffer(30)
	require.NoError(t, err)

	tb.Insert(0, 'a')
	tb.Insert(1, '\n')
	tb.Insert(2, 'b')
	tb.Insert(3, '\n')
	tb.Insert(4, 'c')

	require.Equal(t, []int{0, 2, 4}, tb.lineStarts)

	tb.Insert(1, '\n')

	require.Equal(t, "a\n\nb\nc", tb.String())
	require.Equal(t, []int{0, 2, 3, 5}, tb.lineStarts)
}
