/*
 * slot.go
 *
 * Copyright 2019-2024 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	defSlotColor = pixel.RGBA{0.1, 0.1, 0.1, 0.5} // default color
)

// Struct for slot.
type Slot struct {
	bgSpr               *pixel.Sprite
	drawArea            pixel.Rect
	size                pixel.Vec
	color               color.Color
	fontSize            Size
	label               *Text
	countLabel          *Text
	info                *InfoWindow
	icon                *pixel.Sprite
	values              []interface{}
	mousePos            pixel.Vec
	specialKey          pixelgl.Button
	onRightClick        func(s *Slot)
	onLeftClick         func(s *Slot)
	onSpecialLeftClick  func(s *Slot)
	onSpecialRightClick func(s *Slot)
	hovered             bool
	dragged             bool
}

// NewSlot creates new slot without background.
func NewSlot(params Params) *Slot {
	s := new(Slot)
	// Background.
	s.size = params.Size.SlotSize()
	s.color = params.MainColor
	if s.color == nil {
		s.color = defSlotColor
	}
	// Labels & info.
	s.fontSize = params.FontSize
	labelParams := Params{
		FontSize: s.fontSize,
	}
	s.label = NewText(labelParams)
	s.countLabel = NewText(labelParams)
	s.countLabel.Align(AlignCenter)
	infoParams := Params{
		FontSize:  SizeSmall,
		MainColor: pixel.RGBA{0.1, 0.1, 0.1, 0.5},
	}
	s.info = NewInfoWindow(infoParams)
	return s
}

// SlotSwitch transfers all contant of slot A
// (value, icon, label, info) to slot B and
// vice versa.
func SlotSwitch(slotA, slotB *Slot) {
	slotAValues := slotA.Values()
	slotALabel := slotA.label.String()
	slotAInfo := slotA.info.String()
	slotAIcon := slotA.Icon()
	slotA.SetValues(slotB.Values())
	slotA.SetIcon(slotB.Icon())
	slotA.SetInfo(slotB.info.String())
	slotA.SetLabel(slotB.label.String())
	slotB.SetValues(slotAValues)
	slotB.SetIcon(slotAIcon)
	slotB.SetInfo(slotAInfo)
	slotB.SetLabel(slotALabel)
}

// SlotCopy copies content from slot A to
// slot B(overwrites current content).
func SlotCopy(slotA, slotB *Slot) {
	slotB.SetValues(slotA.Values())
	slotB.SetIcon(slotA.Icon())
	slotB.SetInfo(slotA.info.String())
	slotB.SetLabel(slotA.label.String())
}

// Draw draws slot.
func (s *Slot) Draw(t pixel.Target, matrix pixel.Matrix) {
	s.drawWithoutInfo(t, matrix)
	if s.hovered {
		s.info.Draw(t)
	}
}

// Update updates slot.
func (s *Slot) Update(win *Window) {
	// Mouse position.
	s.mousePos = win.MousePosition()
	// Mouse events.
	if s.DrawArea().Contains(s.mousePos) {
		switch {
		case s.specialKey != 0 && win.Pressed(s.specialKey) &&
			win.JustPressed(pixelgl.MouseButtonRight):
			if s.onSpecialRightClick != nil {
				s.onSpecialRightClick(s)
			}
		case win.JustPressed(pixelgl.MouseButtonRight):
			if s.onRightClick != nil {
				s.onRightClick(s)
			}
		case s.specialKey != 0 && win.Pressed(s.specialKey) &&
			win.JustPressed(pixelgl.MouseButtonLeft):
			if s.onSpecialLeftClick != nil {
				s.onSpecialLeftClick(s)
			}
		case win.JustPressed(pixelgl.MouseButtonLeft):
			if s.onLeftClick != nil {
				s.onLeftClick(s)
			}
		}
	}
	// On-hover.
	s.hovered = s.DrawArea().Contains(s.mousePos)
	// Count label.
	s.countLabel.SetText(fmt.Sprintf("%d", len(s.values)))
	// Elements update.
	s.info.Update(win)
}

// Values returns all slot values.
func (s *Slot) Values() []interface{} {
	return s.values
}

// Pop removes and returns first value
// from slot. Clears slot if removed value
// was last value in slot.
func (s *Slot) Pop() interface{} {
	if s.values == nil {
		return nil
	}
	lastID := len(s.values) - 1
	v := s.values[lastID]
	s.values = s.values[:lastID]
	if len(s.values) < 1 {
		s.Clear()
	}
	return v
}

// Icon returns current slot icon
// picture.
func (s *Slot) Icon() pixel.Picture {
	if s.icon == nil {
		return nil
	}
	return s.icon.Picture()
}

// Label returns slot label text.
func (s *Slot) Label() string {
	return s.label.String()
}

// Drag toggles slot drag mode(icon
// follows mouse cursor).
func (s *Slot) Drag(drag bool) {
	s.dragged = drag
}

// Dragged checks whether slot is in
// drag mode(icon follows mouse cursor).
func (s *Slot) Dragged() bool {
	return s.dragged
}

// SetColor sets specified color as slot
// color.
func (s *Slot) SetColor(c color.Color) {
	s.color = c
}

// SetIcon sets specified sprite as current
// slot icon.
func (s *Slot) SetIcon(pic pixel.Picture) {
	iconBounds := pixel.R(0, 0, s.Size().X, s.Size().Y)
	s.icon = pixel.NewSprite(pic, iconBounds)
}

// AddValue adds specified interface to slot
// values list.
func (s *Slot) AddValues(vls ...interface{}) {
	s.values = append(s.values, vls...)
}

// SetValues replaces current values with
// specified ones.
func (s *Slot) SetValues(vls []interface{}) {
	s.values = vls
}

// SetLabel sets specified text as slot label.
func (s *Slot) SetLabel(text string) {
	s.label.SetText(text)
}

// SetInfo sets specified text as content
// of slot info window.
func (s *Slot) SetInfo(text string) {
	s.info.SetText(text)
}

// Clear removes slot value, icon,
// label and text.
func (s *Slot) Clear() {
	s.SetValues(nil)
	s.SetIcon(nil)
	s.SetLabel("")
	s.SetInfo("")
	s.Drag(false)
}

// DrawArea returns current slot background
// draw area.
func (s *Slot) DrawArea() pixel.Rect {
	return s.drawArea
}

// Size returns slot size.
func (s *Slot) Size() pixel.Vec {
	if s.bgSpr == nil {
		return s.size
	}
	return s.bgSpr.Frame().Size()
}

// SetSpecialKey sets special key for slot click
// events.
func (s *Slot) SetSpecialKey(k pixelgl.Button) {
	s.specialKey = k
}

// SetOnClickFunc set speicfied function as function
// triggered after right mouse click event.
func (s *Slot) SetOnRightClickFunc(f func(s *Slot)) {
	s.onRightClick = f
}

// SetOnLeftClickFunc set speicfied function as function
// triggered after left mouse click event.
func (s *Slot) SetOnLeftClickFunc(f func(s *Slot)) {
	s.onLeftClick = f
}

// SetOnSpecialRightClickFunc set speicfied function as function
// triggered after special key pressed + right mouse click event.
func (s *Slot) SetOnSpecialRightClickFunc(f func(s *Slot)) {
	s.onSpecialRightClick = f
}

// SetOnSpecialLeftClickFunc set speicfied function as function
// triggered after special key pressed + left mouse click event.
func (s *Slot) SetOnSpecialLeftClickFunc(f func(s *Slot)) {
	s.onSpecialLeftClick = f
}

// drawWithoutInfo draws slot without the info window.
func (s *Slot) drawWithoutInfo(t pixel.Target, matrix pixel.Matrix) {
	// Draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Size())
	// Icon.
	if s.icon != nil && s.icon.Picture() != nil {
		if s.dragged {
			s.icon.Draw(t, Matrix().Moved(s.mousePos))
		} else {
			s.icon.Draw(t, matrix)
		}
	}
	// Slot.
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		DrawRect(t, s.DrawArea(), s.color)
	}
	// Labels.
	if len(s.values) > 0 {
		countPos := MoveTL(s.Size(), s.countLabel.Size())
		s.countLabel.Draw(t, matrix.Moved(countPos))
	}
	s.label.Draw(t, matrix)
}
