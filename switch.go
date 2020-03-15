/*
 * switch.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Switch struct represents graphical switch for values.
type Switch struct {
	bgDraw      *imdraw.IMDraw
	bgSpr       *pixel.Sprite
	prevButton  *Button
	nextButton  *Button
	valueText   *Text
	valueSprite *pixel.Sprite
	label       *Text
	info        *InfoWindow
	drawArea    pixel.Rect // updated on each draw
	size        pixel.Vec
	color       color.Color
	index       int
	focused     bool
	disabled    bool
	hovered     bool
	values      []SwitchValue
	onChange    func(s *Switch, old, new *SwitchValue)
}

// Tuple for switch values, contains value to
// display(view) and real value.
type SwitchValue struct {
	View  interface{}
	Value interface{}
}

// NewSwitch creates new switch with specified size and color.
func NewSwitch(params Params) *Switch {
	s := new(Switch)
	// Background.
	s.bgDraw = imdraw.New(nil)
	s.bgSpr = params.Background
	s.size = params.Size.SwitchSize()
	s.color = params.MainColor
	s.bgSpr = params.Background
	// Buttons.
	buttonColor := params.SecColor
	if buttonColor == nil {
		buttonColor = colornames.Red
	}
	buttonParams := Params{
		Size:      params.Size - 2,
		Shape:     ShapeSquare,
		MainColor: buttonColor,
	}
	s.prevButton = NewButton(buttonParams)
	s.prevButton.SetLabel("<")
	s.prevButton.SetOnClickFunc(s.onPrevButtonClicked)
	s.nextButton = NewButton(buttonParams)
	s.nextButton.SetLabel(">")
	s.nextButton.SetOnClickFunc(s.onNextButtonClicked)
	// Label & info.
	labelParams := Params{
		SizeRaw:  pixel.V(s.Size().X, 0),
		FontSize: params.Size - 1,
	}
	s.label = NewText(labelParams)
	s.label.Align(AlignCenter)
	infoParams := Params{
		FontSize:  SizeSmall,
		MainColor: pixel.RGBA{0.1, 0.1, 0.1, 0.5},
	}
	s.info = NewInfoWindow(infoParams)
	// Values.
	s.index = 0
	s.valueText = NewText(labelParams)
	s.updateValueView()
	return s
}

// Draw draws switch.
func (s *Switch) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Size())
	// Background.
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, s.DrawArea(), s.color)
	}
	// Value view.
	valueDA := s.valueText.DrawArea()
	if s.valueSprite == nil {
		s.valueText.Draw(t, matrix)
	} else {
		s.valueSprite.Draw(t, matrix)
		valueDA = MatrixToDrawArea(matrix, s.valueSprite.Frame().Size())
	}
	// Label & info window.
	labelPos := MoveBC(s.Size(), s.label.Size())
	s.label.Draw(t, matrix.Moved(labelPos))
	if s.info != nil && s.hovered {
		s.info.Draw(t)
	}
	// Buttons.
	prevButtonPos := LeftOf(valueDA, s.prevButton.Size(), 10)
	nextButtonPos := RightOf(valueDA, s.nextButton.Size(), 10)
	s.prevButton.Draw(t, Matrix().Moved(prevButtonPos))
	s.nextButton.Draw(t, Matrix().Moved(nextButtonPos))
}

// Update updates switch and all elements.
func (s *Switch) Update(win *Window) {
	if s.Disabled() {
		return
	}
	// Mouse events.
	if s.DrawArea().Contains(win.MousePosition()) {
		s.hovered = true
		if s.info != nil {
			s.info.Update(win)
		}
	} else {
		s.hovered = false
	}
	// Elements update.
	s.prevButton.Update(win)
	s.nextButton.Update(win)
}

// SetBackground sets specified sprite as switch
// background, also removes background color.
func (s *Switch) SetBackground(spr *pixel.Sprite) {
	s.bgSpr = spr
	s.color = nil
}

// SetColor sets specified color as switch background
// color.
func (s *Switch) SetColor(c color.Color) {
	s.color = c
}

// SetLabel sets specified text as label.
func (s *Switch) SetLabel(t string) {
	s.label.SetText(t)
}

// SetInfo sets specified text as info.
func (s *Switch) SetInfo(t string) {
	s.info.SetText(t)
}

// SetNextButtonBackground sets specified sprite as next
// button background.
func (s *Switch) SetNextButtonBackground(spr *pixel.Sprite) {
	s.nextButton.SetBackground(spr)
}

// SetPrevButtonBackground sets specified sprite as previous
// button background.
func (s *Switch) SetPrevButtonBackground(spr *pixel.Sprite) {
	s.prevButton.SetBackground(spr)
}

// SetValues sets specified list with values as switch values.
func (s *Switch) SetValues(values ...SwitchValue) {
	s.values = values
	s.updateValueView()
}

// SwtIntValues adds to witch all integers from
// specified min/max range.
func (s *Switch) SetIntValues(min, max int) {
	values := make([]SwitchValue, max-min)
	for i := min; i < max; i++ {
		values[i] = SwitchValue{i, i}
	}
	s.SetValues(values...)
}

// Focus toggles focus on element.
func (s *Switch) Focus(focus bool) {
	s.focused = focus
}

// Focused checks whether switch is focused.
func (s *Switch) Focused() bool {
	return s.focused
}

// Active toggles switch activity.
func (s *Switch) Active(active bool) {
	s.prevButton.Active(active)
	s.nextButton.Active(active)
	s.disabled = !active
}

// Disabled checks whether switch is active.
func (s *Switch) Disabled() bool {
	return s.disabled
}

// Size returns switch background size.
func (s *Switch) Size() pixel.Vec {
	if s.bgSpr == nil {
		return s.size
	}
	return s.bgSpr.Frame().Size()
}

// DrawArea returns current switch background position and size.
func (s *Switch) DrawArea() pixel.Rect {
	return s.drawArea
}

// Value returns current switch value.
func (s *Switch) Value() *SwitchValue {
	if s.index >= len(s.values) || s.index < 0 {
		return nil
	}
	return &s.values[s.index]
}

// Reset sets value with first index as current
// value.
func (s *Switch) Reset() {
	s.SetIndex(0)
}

// Find checks if switch constains specified value and returns
// index of this value or -1 if switch does not contains
// such value.
func (s *Switch) Find(value interface{}) int {
	for i, v := range s.values {
		if value == v.Value {
			return i
		}
	}
	return -1
}

// Find searches switch values for value with specified index
// and returns this value or nil if switch does not contains
// value with such index.
func (s *Switch) FindValue(index int) *SwitchValue {
	if index >= len(s.values) || index < 0 {
		return nil
	}
	return &s.values[index]
}

// SetIndex sets value with specified index as current value
// of this switch. If specified value is bigger than maximal
// possible index, then index of first value is set, if specified
// index is smaller than minimal, then index of last value is set.
func (s *Switch) SetIndex(index int) {
	if index > len(s.values)-1 {
		s.index = 0
	} else if index < 0 {
		s.index = len(s.values) - 1
	} else {
		s.index = index
	}
	s.updateValueView()
}

// Sets specified function as function triggered on on switch value change.
func (s *Switch) SetOnChangeFunc(f func(s *Switch, old, new *SwitchValue)) {
	s.onChange = f
}

// updateValueView updates value view with current switch value.
func (s *Switch) updateValueView() {
	if s.values == nil || len(s.values) < 1 {
		return
	}
	if pic, ok := s.Value().View.(pixel.Picture); ok {
		s.valueSprite = pixel.NewSprite(pic, pic.Bounds())
	} else {
		label, _ := s.Value().View.(string)
		s.valueText.SetText(label)
	}
}

// Triggered after next button clicked.
func (s *Switch) onNextButtonClicked(b *Button) {
	oldIndex := s.index
	s.SetIndex(s.index + 1)
	if s.onChange != nil {
		oldValue := s.FindValue(oldIndex)
		s.onChange(s, oldValue, s.Value())
	}
}

// Triggered after prev button clicked.
func (s *Switch) onPrevButtonClicked(b *Button) {
	oldIndex := s.index
	s.SetIndex(s.index - 1)
	if s.onChange != nil {
		s.onChange(s, s.FindValue(oldIndex), s.Value())
	}
}
