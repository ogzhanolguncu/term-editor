package screen

import (
	"github.com/gdamore/tcell/v2"
)

type Palette struct {
	lineNumStyle          tcell.Style
	gutterStyle           tcell.Style
	currentLineStyle      tcell.Style
	statusBarStyle        tcell.Style
	statusBarMessageStyle tcell.Style
	normalTextStyle       tcell.Style
}

func NewPalette() *Palette {
	s := tcell.StyleDefault

	// Vesper base + mint accents
	editorBg := tcell.NewRGBColor(16, 16, 16)      // #101010
	currentLineBg := tcell.NewRGBColor(22, 22, 22) // #161616
	textColor := tcell.NewRGBColor(255, 255, 255)  // #FFFFFF
	lineNumColor := tcell.NewRGBColor(80, 80, 80)  // #505050
	mintGreen := tcell.NewRGBColor(153, 255, 228)  // #99FFE4
	darkMint := tcell.NewRGBColor(72, 134, 119)    // #488677

	return &Palette{
		lineNumStyle:          s.Foreground(lineNumColor).Background(editorBg),
		gutterStyle:           s.Foreground(mintGreen).Background(editorBg).Dim(true),
		currentLineStyle:      s.Background(currentLineBg).Foreground(textColor),
		statusBarStyle:        s.Background(darkMint).Foreground(textColor),
		statusBarMessageStyle: s.Background(mintGreen).Foreground(editorBg).Dim(true),
		normalTextStyle:       s.Foreground(textColor).Background(editorBg),
	}
}

func (p *Palette) StyleForLineNum() tcell.Style {
	return p.lineNumStyle
}

func (p *Palette) StyleForGutter() tcell.Style {
	return p.gutterStyle
}

func (p *Palette) StyleForCurrentLine() tcell.Style {
	return p.currentLineStyle
}

func (p *Palette) StyleForStatusBar() tcell.Style {
	return p.statusBarStyle
}

func (p *Palette) StyleForStatusMessage() tcell.Style {
	return p.statusBarMessageStyle
}

func (p *Palette) StyleForNormalText() tcell.Style {
	return p.normalTextStyle
}
