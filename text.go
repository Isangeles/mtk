/*
 * text.go
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
	"bytes"
	"fmt"
	"image/color"
	"strings"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/text"
)

// Text struct for short text like labels, names, etc.
type Text struct {
	*text.Text
	content  string
	drawArea pixel.Rect // updated on each draw
	color    color.Color
	fontSize Size
	width    float64
	align    Align
}

// NewText creates new text with specified
// parameters.
func NewText(p Params) *Text {
	t := new(Text)
	// Parameters.
	t.fontSize = p.FontSize
	t.width = p.SizeRaw.X
	// Text.
	font := MainFont(t.fontSize)
	atlas := Atlas(&font)
	t.Text = text.New(pixel.V(0, 0), atlas)
	t.color = p.MainColor
	if t.color == nil {
		t.color = colornames.White // default color white
	}
	t.align = AlignCenter
	return t
}

// SetText sets specified text as text to display.
func (t *Text) SetText(text string) {
	t.Clear()
	// If text too wide, then split to more lines.
	breakLines := t.breakLine(text, t.width)
	for j := 0; j < len(breakLines); j++ { // reverse order
		bl := breakLines[j]
		t.content = fmt.Sprintf("%s%s", t.content, bl)
	}
	t.Align(t.align)
}

// SetColor sets specified color as
// current text color.
func (tx *Text) SetColor(c color.Color) {
	tx.color = c
}

// SetMaxWidth sets maximal width of single text line.
func (tx *Text) SetMaxWidth(width float64) {
	tx.width = width
}

// Align aligns text to specified position.
func (t *Text) Align(a Align) {
	t.align = a
	switch a {
	case AlignCenter:
		mariginX := (-t.BoundsOf(t.content).Max.X) / 2
		t.Orig = pixel.V(mariginX, 0)
		t.Text.Clear()
		t.WriteString(t.content)
	case AlignRight:
		mariginX := (-t.BoundsOf(t.content).Max.X)
		t.Orig = pixel.V(mariginX, 0)
		t.Text.Clear()
		t.WriteString(t.content)
	case AlignLeft:
		t.Orig = pixel.V(0, 0)
		t.Text.Clear()
		t.WriteString(t.content)
	}
}

// Draw draws text.
func (tx *Text) Draw(t pixel.Target, matrix pixel.Matrix) {
	tx.drawArea = MatrixToDrawArea(matrix, tx.Size())
	tx.Text.DrawColorMask(t, matrix, tx.color)
}

// Size returns size of current text.
func (tx *Text) Size() pixel.Vec {
	return tx.Text.Bounds().Size()
}

// Clear clears texts,
func (t *Text) Clear() {
	t.Text.Clear()
	t.content = ""
	t.Align(t.align)
}

// DrawArea returns current draw area of text.
func (tx *Text) DrawArea() pixel.Rect {
	return tx.drawArea
}

// String returns text content.
func (tx *Text) String() string {
	return tx.content
}

// breakLine breaks specified line into few lines with specified
// maximal width.
func (t *Text) breakLine(line string, width float64) []string {
	lines := make([]string, 0)
	lineWidth := t.BoundsOf(line).W()
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
func (t *Text) breakPoint(line string, width float64) int {
	checkLine := ""
	breakPoint := -1
	for _, c := range line {
		if c == '\n' {
			breakPoint = -1
		}
		checkLine += string(c)
		breakPoint++
		if t.BoundsOf(checkLine).W() >= width {
			return breakPoint
		}
	}
	return len(line) - 1
}

// Splits string to chunks with n as max chunk width.
// Author: mozey(@stackoverflow).
func SplitSubN(s string, n int) []string {
	if n == 0 {
		return []string{s}
	}
	sub := ""
	subs := []string{}
	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}
