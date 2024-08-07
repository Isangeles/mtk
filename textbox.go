/*
 * textbox.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// Struct for textboxes.
type Textbox struct {
	bgSize      pixel.Vec
	color       color.Color
	textarea    *Text
	drawArea    pixel.Rect // updated at every draw
	upButton    *Button
	downButton  *Button
	textContent []string // every line of text content
	visibleText []string
	startID     int
	buttons     bool
	focused     bool
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
		SizeRaw:  pixel.V(t.bgSize.X, 0),
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
	t.upButton.SetOnClickFunc(t.onButtonUpClicked)
	t.downButton = NewButton(buttonParams)
	t.downButton.SetOnClickFunc(t.onButtonDownClicked)
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Background.
	tb.drawArea = MatrixToDrawArea(matrix, tb.Size())
	DrawRect(t, tb.DrawArea(), pixel.RGBA{0.1, 0.1, 0.1, 0.5})
	// Text content.
	textareaPos := pixel.V(tb.DrawArea().Min.X, tb.DrawArea().Max.Y-ConvSize(tb.textarea.Size().Y))
	tb.textarea.Draw(t, Matrix().Moved(textareaPos))
	// Buttons.
	upButtonPos := MoveTR(tb.Size(), tb.upButton.Size())
	downButtonPos := MoveBR(tb.Size(), tb.downButton.Size())
	tb.upButton.Draw(t, matrix.Moved(upButtonPos))
	tb.downButton.Draw(t, matrix.Moved(downButtonPos))
}

// Update handles key events.
func (tb *Textbox) Update(win *Window) {
	// Key events.
	if tb.Focused() {
		if win.JustPressed(pixelgl.KeyDown) {
			if tb.startID < len(tb.textContent)-1 {
				tb.startID++
			}
		}
		if win.JustPressed(pixelgl.KeyUp) {
			if tb.startID > 0 {
				tb.startID--
			}
		}
	}
	// Elements.
	tb.upButton.Update(win)
	tb.downButton.Update(win)
	tb.updateTextVisibility()
}

// SetSize sets background size.
func (tb *Textbox) SetSize(s pixel.Vec) {
	tb.bgSize = s
}

// Size returns size of textbox background.
func (tb *Textbox) Size() pixel.Vec {
	return tb.bgSize
}

// Focus sets/removes focus from textbox.
func (tb *Textbox) Focus(focus bool) {
	tb.focused = focus
}

// Focused checks if textbox is focused.
func (tb *Textbox) Focused() bool {
	return tb.focused
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


// SetUpButtonBackground sets specified sprite as scroll up
// button background.
func (tb *Textbox) SetUpButtonBackground(s *pixel.Sprite) {
	tb.upButton.SetBackground(s)
	tb.upButton.SetColor(nil)
}

// SetDownButtonBackground sets specified sprite as scroll
// down button background.
func (tb *Textbox) SetDownButtonBackground(s *pixel.Sprite) {
	tb.downButton.SetBackground(s)
	tb.downButton.SetColor(nil)
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
	tb.startID = len(tb.textContent) - 1
}

// AddText adds specified text to box.
func (tb *Textbox) AddText(text string) {
	tb.textContent = append(tb.textContent, text)
}

// Clear clears textbox.
func (tb *Textbox) Clear() {
	tb.textContent = []string{}
}

// String returns textbox content.
func (tb *Textbox) String() string {
	content := ""
	for _, line := range tb.textContent {
		content = fmt.Sprintf("%s%s", content, line)
	}
	return content
}

// AtBottom checks if textbox is scrolled to the bottom.
func (tb *Textbox) AtBottom() bool {
	return (tb.startID == 0 && len(tb.textContent) < 1) || tb.startID == len(tb.textContent)-1
}

// ScrollBottom scrolls textbox to last lines
// of text content.
func (tb *Textbox) ScrollBottom() {
	if len(tb.textContent) < 1 {
		return
	}
	tb.startID = len(tb.textContent) - 1
}

// updateTextVisibility updates conte nt of visible
// text area.
func (tb *Textbox) updateTextVisibility() {
	var (
		visibleText       []string
		visibleTextHeight float64
	)
	boxWidth := tb.Size().X
	for i := len(tb.textContent) - 1; i >= 0; i-- {
		if i > tb.startID {
			continue
		}
		if visibleTextHeight >= tb.Size().Y {
			break
		}
		line := tb.textContent[i]
		breakLines := tb.breakLine(line, boxWidth)
		for j := len(breakLines) - 1; j >= 0; j-- { // reverse order
			bl := breakLines[j]
			visibleText = append(visibleText, bl)
			visibleTextHeight += tb.textarea.BoundsOf(bl).H()
			if visibleTextHeight >= tb.Size().Y {
				break
			}
		}
	}
	tb.textarea.Clear()
	for i := len(visibleText) - 1; i >= 0; i-- {
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
	checkLine := ""
	breakPoint := -1
	for _, c := range line {
		if c == '\n' {
			checkLine = ""
			breakPoint = -1
		}
		checkLine += string(c)
		breakPoint++
		if t.textarea.BoundsOf(checkLine).W() >= width {
			return breakPoint
		}
	}
	return len(line)
}

// Triggered after button up clicked.
func (tb *Textbox) onButtonUpClicked(b *Button) {
	if tb.startID <= 0 {
		return
	}
	tb.startID--
}

// Triggered after button down clicked.
func (tb *Textbox) onButtonDownClicked(b *Button) {
	if tb.startID >= len(tb.textContent) {
		return
	}
	tb.startID++
}
