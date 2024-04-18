/*
 * button.go
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
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/gopxl/beep"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	buttonPushColor  = colornames.Grey
	buttonHoverColor = colornames.Crimson
)

// Button struct for UI button.
type Button struct {
	bgSpr      *pixel.Sprite
	bgDraw     *imdraw.IMDraw
	label      *Text
	info       *InfoWindow
	size       pixel.Vec
	shape      Shape
	color      color.Color
	colorPush  color.Color
	colorHover color.Color
	pressed    bool
	focused    bool
	hovered    bool
	disabled   bool
	drawArea   pixel.Rect // updated on each draw
	onClick    func(b *Button)
	clickSound *beep.Buffer
}

// NewButton creates new button with specified parameters.
func NewButton(params Params) *Button {
	b := new(Button)
	// Parameters.
	b.shape = params.Shape
	b.size = params.Size.ButtonSize(b.shape)
	b.color = params.MainColor
	b.colorPush = buttonPushColor
	b.colorHover = buttonHoverColor
	// Background.
	b.bgDraw = imdraw.New(nil)
	b.bgSpr = params.Background
	// Label.
	labelParams := Params{
		SizeRaw:  pixel.V(b.Size().X, 0),
		FontSize: params.FontSize,
	}
	b.label = NewText(labelParams)
	// Info window.
	infoParams := Params{
		FontSize:  SizeSmall,
		MainColor: pixel.RGBA{0.1, 0.1, 0.1, 0.5},
	}
	b.info = NewInfoWindow(infoParams)
	// Global click sound.
	b.SetClickSound(buttonClickSound)
	return b
}

// Draw draws button.
func (b *Button) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	b.drawArea = MatrixToDrawArea(matrix, b.Size())
	// Drawing background.
	bgColor := b.color
	if b.pressed || b.Disabled() {
		bgColor = b.colorPush
	} else if b.hovered {
		bgColor = b.colorHover
	}
	if b.bgSpr != nil {
		if bgColor == nil {
			b.bgSpr.Draw(t, matrix)
		} else {
			b.bgSpr.DrawColorMask(t, matrix, bgColor)
		}
	} else {
		DrawRectangle(t, b.DrawArea(), bgColor)
	}
	// Drawing label.
	if b.label != nil {
		labelPos := pixel.V(0, -b.label.Size().Y/2)
		b.label.Draw(t, matrix.Moved(labelPos))
	}
	// Info window.
	if b.info != nil && b.hovered {
		b.info.Draw(t)
	}
}

// Update updates button.
func (b *Button) Update(win *Window) {
	if b.Disabled() {
		return
	}
	// Mouse events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if b.DrawArea().Contains(win.MousePosition()) {
			b.pressed = true
		}
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) {
		if b.pressed && b.DrawArea().Contains(win.MousePosition()) {
			if b.onClick != nil {
				b.onClick(b)
			}
			if audio != nil && b.clickSound != nil {
				audio.Play(b.clickSound)
			}
		}
		b.pressed = false
	}
	// On-hover.
	b.hovered = b.DrawArea().Contains(win.MousePosition())
	if b.hovered || b.Focused() {
		if b.info != nil {
			b.info.Update(win)
		}
	}
	// On-focus events.
	if b.Focused() {
		if win.JustPressed(pixelgl.KeyEnter) {
			b.onClick(b)
		}
	}
}

// SetBackground sets specified sprite as button
// background, also removes background color.
func (b *Button) SetBackground(s *pixel.Sprite) {
	b.bgSpr = s
	b.color = nil
}

// SetColor sets specified color as current
// button backgroun color.
func (b *Button) SetColor(c color.Color) {
	b.color = c
}

// SetLabel sets specified text as button label.
func (b *Button) SetLabel(t string) {
	b.label.SetText(t)
}

// SetInfo sets specified text as content of
// button info window.
func (b *Button) SetInfo(t string) {
	b.info.SetText(t)
}

// Focus sets/removes focus from button
func (b *Button) Focus(focus bool) {
	b.focused = focus
}

// Focused checks whether buttons is focused.
func (b *Button) Focused() bool {
	return b.focused
}

// Active toggles button active state.
func (b *Button) Active(active bool) {
	b.disabled = !active
}

// Disabled checks whether button is disabled.
func (b *Button) Disabled() bool {
	return b.disabled
}

// SetOnClickFunc sets specified function as on-click
// callback function.
func (b *Button) SetOnClickFunc(callback func(b *Button)) {
	b.onClick = callback
}

// SetClickSound sets specified audio buffer as
// on-click audio effect.
func (b *Button) SetClickSound(s *beep.Buffer) {
	b.clickSound = s
}

// DrawArea returns current button background position and size.
func (b *Button) DrawArea() pixel.Rect {
	return b.drawArea
}

// Size returns button background size.
func (b *Button) Size() pixel.Vec {
	if b.bgSpr == nil {
		return b.size
	}
	return b.bgSpr.Frame().Size()
}
