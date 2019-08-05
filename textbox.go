/*
 * textbox.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301, USA.
 *
 *
 */

package mtk

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Struct for textboxes.
type Textbox struct {
	bgSize      pixel.Vec
	color       color.Color
	textarea    *Text
	drawArea    pixel.Rect // updated at every draw
	upButton    *Button
	downButton  *Button
	textContent []string   // every line of text content
	visibleText []string
	startID     int
	buttons     bool
}

// NewTextbox creates new textbox with specified
// parameters.
func NewTextbox(params Params) *Textbox {
	t := new(Textbox)
	// Background.
	t.bgSize = params.SizeRaw
	t.color = params.MainColor
	// Text.
	textParams := Params{
		SizeRaw: pixel.V(t.bgSize.X, 0),
		FontSize: params.FontSize,
	}
	t.textarea = NewText(textParams)
	t.textarea.Align(AlignLeft)
	// Buttons.
	buttonParams := Params{
		Size:      SizeMini,
		FontSize:  SizeMedium,
		Shape:     ShapeSquare,
		MainColor: params.AccentColor,
	}
	t.upButton = NewButton(buttonParams)
	t.upButton.SetLabel("^")
	t.upButton.SetOnClickFunc(t.onButtonUpClicked)
	t.downButton = NewButton(buttonParams)
	t.downButton.SetLabel(".")
	t.downButton.SetOnClickFunc(t.onButtonDownClicked)
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Background.
	tb.drawArea = MatrixToDrawArea(matrix, tb.Size())
	DrawRectangle(t, tb.DrawArea(), pixel.RGBA{0.1, 0.1, 0.1, 0.5})
	// Text content.
	tb.textarea.Draw(t, Matrix().Moved(pixel.V(tb.DrawArea().Min.X,
		tb.DrawArea().Max.Y-tb.textarea.BoundsOf("AA").H())))
	// Buttons.
	upButtonPos := MoveTR(tb.Size(), tb.upButton.Size())
	downButtonPos := MoveBR(tb.Size(), tb.downButton.Size())
	tb.upButton.Draw(t, matrix.Moved(upButtonPos))
	tb.downButton.Draw(t, matrix.Moved(downButtonPos))
}

// Update handles key events.
func (tb *Textbox) Update(win *Window) {
	// Key events.
	if win.JustPressed(pixelgl.KeyDown) {
		if tb.startID < len(tb.textContent)-1 {
			tb.startID++
			tb.updateTextVisibility()
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		if tb.startID > 0 {
			tb.startID--
			tb.updateTextVisibility()
		}
	}
	// Elements.
	tb.upButton.Update(win)
	tb.downButton.Update(win)
}

// SetSize sets background size.
func (tb *Textbox) SetSize(s pixel.Vec) {
	tb.bgSize = s
}

// Size returns size of textbox background.
func (tb *Textbox) Size() pixel.Vec {
	return tb.bgSize
}

// TextSize returns size of text content.
func (tb *Textbox) TextSize() pixel.Vec {
	return ConvVec(tb.textarea.Size())
}

// DrawArea returns current draw area of text box
// background.
func (t *Textbox) DrawArea() pixel.Rect {
	return t.drawArea
}

// SetMaxTextWidth sets maximal width of single
// line in text area.
func (tb *Textbox) SetMaxTextWidth(width float64) {
	tb.textarea.SetMaxWidth(width)
}

// SetText clears textbox and inserts specified
// lines of text.
func (tb *Textbox) SetText(text ...string) {
	tb.Clear()
	tb.textContent = text
	tb.startID = len(tb.textContent)-1
	tb.updateTextVisibility()
}

// AddText adds specified text to box.
func (tb *Textbox) AddText(text string) {
	tb.textContent = append(tb.textContent, text)
	tb.startID = len(tb.textContent)-1
	tb.updateTextVisibility()
}

// Clear clears textbox.
func (tb *Textbox) Clear() {
	tb.textContent = []string{}
	tb.updateTextVisibility()
}

// String returns textbox content.
func (tb *Textbox) String() string {
	content := ""
	for _, line := range tb.textContent {
		content = fmt.Sprintf("%s%s", content, line)
	}
	return content
}

// ScrollBottom scrolls textbox to last lines
// of text content.
func (tb *Textbox) ScrollBottom() {
	tb.startID = len(tb.textContent)-1
}

// updateTextVisibility updates conte nt of visible
// text area.
func (tb *Textbox) updateTextVisibility() {
	/*
	tb.textarea.Clear()
	for i := 0; i < len(tb.textContent); i++ {
		if i < tb.startID {
			continue
		}
		text := tb.textarea.Content()
		tb.textarea.SetText(text + tb.textContent[i])
	}
	for tb.textarea.Size().Y > tb.Size().Y {
		lines := strings.Split(tb.textarea.Content(), "\n")
		text := ""
		for _, l := range lines[1:] {
			text = fmt.Sprintf("%s\n%s", text, l)
		}
		tb.textarea.SetText(text)
	}
        */
	var (
		visibleText       []string
		visibleTextHeight float64
	)
	boxWidth := tb.Size().X
	for i := len(tb.textContent)-1; i >= 0; i-- {
		if i > tb.startID {
			continue
		}
		if visibleTextHeight >= tb.Size().Y {
			break
		}
		line := tb.textContent[i]
		breakLines := tb.breakLine(line, boxWidth)
		for j := len(breakLines)-1; j >= 0; j-- { // reverse order
			bl := breakLines[j]
			visibleText = append(visibleText, bl)
		}
		visibleTextHeight += tb.textarea.BoundsOf(line).H() * float64(len(breakLines))
	}
	tb.textarea.Clear()
	for i := len(visibleText)-1; i >= 0; i-- {
		txt := visibleText[i]
		fmt.Fprintf(tb.textarea, txt)
	}
}

// breakLine breaks specified line into few lines with specified
// maximal width.
func (t *Textbox) breakLine(line string, width float64) []string {
	lines := make([]string, 0)
	lineWidth := t.textarea.BoundsOf(line).W()
	if width > 0 && lineWidth > width {
		breakPoint := t.breakPoint(line, width)
		breakLines := SplitSubN(line, breakPoint)
		for i, l := range breakLines {
			if !strings.HasSuffix(l, "\n") {
				breakLines[i] += "\n"
			}
		}
		lines = append(lines, breakLines...)
	} else {
		lines = append(lines, line)
	}
	return lines
}

// breakPoint return break position for specified line and width.
func (t *Textbox) breakPoint(line string, width float64) int {
	/*
	checkLine := ""
	for i, c := range line {
		checkLine += string(c)
		if t.textarea.BoundsOf(checkLine).W() >= width {
			return i
		}
	}
	return len(line)-1
        */
	checkLine := ""
	breakPoint := -1
	for _, c := range line {
		if c == '\n' {
			breakPoint = -1
		}
		checkLine += string(c)
		breakPoint++
		if t.textarea.BoundsOf(checkLine).W() >= width {
			return breakPoint
		}
	}
	return len(line)-1
}

// Triggered after button up clicked.
func (tb *Textbox) onButtonUpClicked(b *Button) {
	if tb.startID <= 0 {
		return
	}
	tb.startID--
	tb.updateTextVisibility()
}

// Triggered after button down clicked.
func (tb *Textbox) onButtonDownClicked(b *Button) {
	if tb.startID >= len(tb.textContent) {
		return
	}
	tb.startID++
	tb.updateTextVisibility()
}
