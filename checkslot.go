/*
 * checkslot.go
 *
 * Copyright 2018-2026 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// Struct for 'chackable' slots.
type CheckSlot struct {
	checked    bool
	bgSize     pixel.Vec
	bgColor    color.Color
	checkColor color.Color
	drawArea   pixel.Rect
	label      *Text
	value      interface{}
	onCheck    func(s *CheckSlot)
}

// NewCheckSlot creates new checkable slot.
func NewCheckSlot(params Params) *CheckSlot {
	cs := new(CheckSlot)
	if params.SizeRaw.X > 0 && params.SizeRaw.Y > 0 {
		cs.bgSize = params.SizeRaw
	} else {
		cs.bgSize = params.Size.ButtonSize(params.Shape)
	}
	cs.bgColor = params.MainColor
	cs.checkColor = params.SecColor
	labelParams := Params{FontSize: params.FontSize}
	cs.label = NewText(labelParams)
	cs.label.SetText(params.Label)
	return cs
}

// Draw draws slot.
func (cs *CheckSlot) Draw(t pixel.Target, matrix pixel.Matrix) {
	cs.drawArea = MatrixToDrawArea(matrix, cs.Size())
	color := cs.bgColor
	if cs.Checked() {
		color = cs.checkColor
	}
	DrawRect(t, cs.DrawArea(), color)
	labelMove := MoveBL(cs.Size(), cs.label.Size())
	cs.label.Draw(t, matrix.Moved(labelMove))
}

// Update updates slot.
func (cs *CheckSlot) Update(win *Window) {
	// Mouse events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if cs.DrawArea().Contains(win.MousePosition()) {
			if cs.Checked() {
				cs.Check(false)
			} else {
				cs.Check(true)
				if cs.onCheck != nil {
					cs.onCheck(cs)
				}
			}
		}
	}
}

// SetSize sets specified vector as
// background size.
func (cs *CheckSlot) SetSize(s pixel.Vec) {
	cs.bgSize = s
}

// Label returns slot label.
func (cs *CheckSlot) Label() string {
	return cs.label.String()
}

// Set value sets the check slot value.
func (cs *CheckSlot) SetValue(value interface{}) {
	cs.value = value
}

// Value returns slot value.
func (cs *CheckSlot) Value() interface{} {
	return cs.value
}

// Size returns slot background size.
func (cs *CheckSlot) Size() pixel.Vec {
	return cs.bgSize
}

// DrawArea returns current slot draw area.
func (cs *CheckSlot) DrawArea() pixel.Rect {
	return cs.drawArea
}

// Check toggles slot selection.
func (cs *CheckSlot) Check(check bool) {
	cs.checked = check
}

// Checked checks whether slot is checked.
func (cs *CheckSlot) Checked() bool {
	return cs.checked
}

// SetOnCheckFunc sets specified function as function triggered
// after slot was selected.
func (cs *CheckSlot) SetOnCheckFunc(f func(s *CheckSlot)) {
	cs.onCheck = f
}
