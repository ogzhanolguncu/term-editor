package screen

import (
	"fmt"
	"strings"

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

func (ui *Screen) Run() {
	for {
		ui.render()

		ev := ui.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if !ui.handleKey(ev) {
				return // Quit signal
			}

		case *tcell.EventResize:
			ui.screen.Sync()
		}
	}
}

func (ui *Screen) Close() {
	ui.screen.Fini()
}

func (ui *Screen) render() {
	ui.screen.Clear()
	ui.width, ui.height = ui.screen.Size()
	ui.fillBg()

	// Content height has to account for status bar thats why we skip the last line
	contentHeight := ui.height - statusBarHeight
	cursorLine, cursorCol := ui.editor.GetLineColumn()

	if cursorLine < ui.yOffset {
		ui.yOffset = cursorLine
	}
	if cursorLine >= ui.yOffset+contentHeight {
		ui.yOffset = cursorLine - contentHeight + 1
	}

	// Calculate gutter width based on max line number
	gutterWidth := len(fmt.Sprintf("%d", ui.yOffset+contentHeight)) + 1 // +1 for space after number
	availableWidth := ui.width - gutterWidth - gutterPadding

	if cursorCol < ui.xOffset {
		ui.xOffset = cursorCol
	}
	if cursorCol >= ui.xOffset+availableWidth {
		ui.xOffset = cursorCol - availableWidth + 1
	}

	textStartCol := gutterWidth + gutterPadding
	ui.renderLines(gutterWidth, cursorLine, textStartCol)
	ui.renderStatusBar()

	screenRow := cursorLine - ui.yOffset
	screenCol := (cursorCol - ui.xOffset) + textStartCol

	ui.screen.ShowCursor(screenCol, screenRow)
	ui.screen.Show()
}

func (ui *Screen) fillBg() {
	bgStyle := ui.palette.StyleForNormalText()
	for y := 0; y < ui.height; y++ {
		for x := 0; x < ui.width; x++ {
			ui.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}
}

func (ui *Screen) renderLines(gutterWidth, cursorLine, textStartCol int) {
	availableWidth := ui.width - gutterWidth - gutterPadding
	availableHeight := ui.height - statusBarHeight
	lines := ui.editor.GetVisibleContent(ui.yOffset, availableHeight)

	for row, lineContent := range lines {
		lineNum := ui.yOffset + row + 1

		// Draw line number
		gutterText := fmt.Sprintf("%*d ", gutterWidth-1, lineNum)
		for col, ch := range gutterText {
			ui.screen.SetContent(col, row, ch, nil, ui.palette.StyleForLineNum())
		}

		// Draw empty vertical separator
		ui.screen.SetContent(gutterWidth, row, ' ', nil, ui.palette.StyleForGutter())

		// Draw text content
		style := ui.palette.StyleForNormalText()
		if ui.yOffset+row == cursorLine {
			style = ui.palette.StyleForCurrentLine()
		}

		visibleContent := ui.getVisibleSlice(lineContent, ui.xOffset, availableWidth)

		ui.drawLine(textStartCol, row, visibleContent, style)
	}
}

func (ui *Screen) getVisibleSlice(line string, offset, width int) string {
	runes := []rune(line)
	end := min(offset+width, len(runes))

	return string(runes[offset:end])
}

func (ui *Screen) renderStatusBar() {
	statusLine := ui.editor.GetStatusLine()
	ui.drawLine(0, ui.height-statusBarHeight, statusLine, ui.palette.StyleForStatusBar())

	// Show message if any
	if msg := ui.editor.GetMessage(); msg != "" {
		ui.drawLine(0, ui.height-statusBarHeight, msg, ui.palette.StyleForStatusMessage())
	}
}

func (ui *Screen) drawLine(x, y int, text string, style tcell.Style) {
	col := x
	for _, ch := range text {
		if col >= ui.width {
			break
		}
		ui.screen.SetContent(col, y, ch, nil, style)
		col++
	}

	// Fill rest of line with spaces (for background color)
	for col < ui.width {
		ui.screen.SetContent(col, y, ' ', nil, style)
		col++
	}
}

func (ui *Screen) handleKey(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlQ:
		return false // Quit

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		ui.editor.Backspace()

	case tcell.KeyDelete:
		ui.editor.Delete()

	case tcell.KeyEnter:
		ui.editor.InsertChar('\n')

	case tcell.KeyLeft:
		ui.editor.MoveLeft()

	case tcell.KeyRight:
		ui.editor.MoveRight()

	case tcell.KeyUp:
		ui.editor.MoveUp()

	case tcell.KeyDown:
		ui.editor.MoveDown()

	case tcell.KeyHome:
		ui.editor.MoveToLineStart()

	case tcell.KeyEnd:
		ui.editor.MoveToLineEnd()

	case tcell.KeyCtrlA:
		ui.editor.MoveToStart()

	case tcell.KeyCtrlE:
		ui.editor.MoveToEnd()

	case tcell.KeyRune:
		ui.editor.InsertChar(ev.Rune())

	case tcell.KeyTAB:
		ui.editor.InsertString(strings.Repeat(" ", tabSize))
	}

	return true
}
