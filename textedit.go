/*
 * textedit.go
 *
 * Copyright 2018-2024 Dariusz Sikora <dev@isangeles.pl>
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
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
)

// Struct for text edit fields.
type Textedit struct {
	size       pixel.Vec
	drawArea   pixel.Rect
	color      color.Color
	colorFocus color.Color
	input      *text.Text
	text       string
	focused    bool
	disabled   bool
}

// NewTextedit creates new textedit based on
// specified parameters.
func NewTextedit(params Params) *Textedit {
	t := new(Textedit)
	// Background.
	t.color = params.MainColor
	t.colorFocus = colornames.Crimson
	// Text input.
	t.size = params.SizeRaw
	font := MainFont(params.FontSize)
	atlas := Atlas(&font)
	t.input = text.New(pixel.V(0, 0), atlas)
	return t
}

// Draw draws text edit.
func (te *Textedit) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Draw area.
	te.drawArea = MatrixToDrawArea(matrix, te.Size())
	color := te.color
	if te.Focused() {
		color = te.colorFocus
	}
	DrawRect(t, te.DrawArea(), color)
	// Text input.
	inputMove := pixel.V(-te.Size().X/2, 0)
	te.input.Draw(t, matrix.Moved(inputMove))
}

// Update updates text edit.
func (te *Textedit) Update(win *Window) {
	if te.Disabled() {
		return
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if te.DrawArea().Contains(win.MousePosition()) {
			te.Focus(true)
		} else {
			te.Focus(false)
		}
	}
	if te.focused {
		if win.JustPressed(pixelgl.KeyBackspace) {
			if len(te.text) > 0 {
				te.SetText(te.text[:len(te.text)-1])
			}
		}
		te.SetText(te.text + win.Typed())
	}
	te.input.Clear()
	te.input.WriteString(te.text)
}

// Focus sets or removes focus from text edit.
func (te *Textedit) Focus(focus bool) {
	te.focused = focus
}

// Focused checks whether text edit is focused.
func (te *Textedit) Focused() bool {
	return te.focused
}

// Active toggles field activity.
func (te *Textedit) Active(active bool) {
	te.disabled = !active
}

// Disabled checks whether field is disabled.
func (te *Textedit) Disabled() bool {
	return te.disabled
}

// Clear clears text edit input.
func (te *Textedit) Clear() {
	te.text = ""
}

// Text return current value of text edit.
func (te *Textedit) Text() string {
	return te.text
}

// SetText sets specified text as current
// value of text edit field.
func (te *Textedit) SetText(text string) {
	te.text = text
}

// SetSize sets text edit size.
func (te *Textedit) SetSize(size pixel.Vec) {
	te.size = size
}

// Size returns text edit size.
func (te *Textedit) Size() pixel.Vec {
	return te.size
}

// DrawArea returns current draw area rectangle.
func (te *Textedit) DrawArea() pixel.Rect {
	return te.drawArea
}
