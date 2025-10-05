package screen

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/ogzhanolguncu/go_editor/editor"
)

type Screen struct {
	screen  tcell.Screen
	editor  *editor.Editor
	palette *Palette

	width   int
	height  int
	yOffset int
	xOffset int
}

const (
	statusBarHeight = 1
	gutterPadding   = 2 // Space between gutter and text (separator + margin)
	tabSize         = 4
)

func NewScreen(editor *editor.Editor) (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	width, height := screen.Size()

	return &Screen{
		screen:  screen,
		editor:  editor,
		palette: NewPalette(),

		width:  width,
		height: height,

		yOffset: 0,
		xOffset: 0,
	}, nil
}

func (s *Screen) Run() {
	for {
		s.render()
		ev := s.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if !s.handleKey(ev) {
				return
			}

		case *tcell.EventResize:
			s.screen.Sync()
		}
	}
}

func (s *Screen) Close() {
	s.screen.Fini()
}

func (s *Screen) render() {
	s.screen.Clear()
	s.width, s.height = s.screen.Size()
	s.fillBg()

	// Content height has to account for status bar thats why we skip the last line
	contentHeight := s.height - statusBarHeight
	cursorLine, cursorCol := s.editor.GetLineColumn()

	if cursorLine < s.yOffset {
		s.yOffset = cursorLine
	}
	if cursorLine >= s.yOffset+contentHeight {
		s.yOffset = cursorLine - contentHeight + 1
	}

	// Calculate gutter width based on max line number
	gutterWidth := len(fmt.Sprintf("%d", s.yOffset+contentHeight)) + 1 // +1 for space after number
	availableWidth := s.width - gutterWidth - gutterPadding

	if cursorCol < s.xOffset {
		s.xOffset = cursorCol
	}
	if cursorCol >= s.xOffset+availableWidth {
		s.xOffset = cursorCol - availableWidth + 1
	}

	textStartCol := gutterWidth + gutterPadding
	s.renderLines(gutterWidth, cursorLine, textStartCol)
	s.renderStatusBar()

	screenRow := cursorLine - s.yOffset
	screenCol := (cursorCol - s.xOffset) + textStartCol

	s.screen.ShowCursor(screenCol, screenRow)
	if s.editor.GetMode() == editor.ModeNormal {
		s.screen.SetCursorStyle(tcell.CursorStyleSteadyBlock)
	} else {
		s.screen.SetCursorStyle(tcell.CursorStyleSteadyBar)
	}
	s.screen.Show()
}

func (s *Screen) fillBg() {
	bgStyle := s.palette.StyleForNormalText()
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}
}

func (s *Screen) renderLines(gutterWidth, cursorLine, textStartCol int) {
	availableWidth := s.width - gutterWidth - gutterPadding
	availableHeight := s.height - statusBarHeight
	lines := s.editor.GetVisibleContent(s.yOffset, availableHeight)

	for row, lineContent := range lines {
		lineNum := s.yOffset + row + 1

		// Draw line number
		gutterText := fmt.Sprintf("%*d ", gutterWidth-1, lineNum)
		for col, ch := range gutterText {
			s.screen.SetContent(col, row, ch, nil, s.palette.StyleForLineNum())
		}

		// Draw empty vertical separator
		s.screen.SetContent(gutterWidth, row, ' ', nil, s.palette.StyleForGutter())

		// Draw text content
		style := s.palette.StyleForNormalText()
		if s.yOffset+row == cursorLine {
			style = s.palette.StyleForCurrentLine()
		}

		visibleContent := s.getVisibleSlice(lineContent, s.xOffset, availableWidth)

		s.drawLine(textStartCol, row, visibleContent, style)
	}
}

func (s *Screen) getVisibleSlice(line string, offset, width int) string {
	runes := []rune(line)
	end := min(offset+width, len(runes))

	return string(runes[offset:end])
}

func (ui *Screen) renderStatusBar() {
	mode := ui.editor.GetMode()
	modeStr := "NORMAL"
	if mode == editor.ModeInsert {
		modeStr = "INSERT"
	}

	statusLine := fmt.Sprintf(" %s | %s", modeStr, ui.editor.GetStatusLine())
	ui.drawLine(0, ui.height-statusBarHeight, statusLine, ui.palette.StyleForStatusBar(mode))

	if msg := ui.editor.GetMessage(); msg != "" {
		ui.drawLine(0, ui.height-statusBarHeight, msg, ui.palette.StyleForStatusMessage())
	}
}

func (s *Screen) drawLine(x, y int, text string, style tcell.Style) {
	col := x
	for _, ch := range text {
		if col >= s.width {
			break
		}
		s.screen.SetContent(col, y, ch, nil, style)
		col++
	}

	// Fill rest of line with spaces (for background color)
	for col < s.width {
		s.screen.SetContent(col, y, ' ', nil, style)
		col++
	}
}
