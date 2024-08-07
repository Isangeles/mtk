/*
 * messagewindow.go
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

// MessageWindow struct represents UI message window.
type MessageWindow struct {
	drawArea     pixel.Rect
	size         pixel.Vec
	color        color.Color
	colorDisable color.Color
	textbox      *Textbox
	acceptButton *Button
	cancelButton *Button
	opened       bool
	focused      bool
	accepted     bool
	dismissed    bool
	disabled     bool
	onAccept     func(msg *MessageWindow)
	onCancel     func(msg *MessageWindow)
}

// NewMessageWindow creates new message window instance.
func NewMessageWindow(params Params) *MessageWindow {
	mw := new(MessageWindow)
	mw.opened = true
	// Background.
	mw.size = params.Size.MessageWindowSize()
	mw.color = params.MainColor
	mw.colorDisable = colornames.Darkgrey
	// Buttons.
	buttonParams := Params{
		Size:      SizeSmall,
		FontSize:  SizeSmall,
		Shape:     ShapeRectangle,
		MainColor: params.SecColor,
	}
	mw.acceptButton = NewButton(buttonParams)
	mw.acceptButton.SetOnClickFunc(mw.onAcceptButtonClicked)
	// Textbox.
	boxSize := pixel.V(mw.Size().X, mw.Size().Y-mw.acceptButton.Size().Y)
	boxParams := Params{
		SizeRaw:     boxSize,
		FontSize:    params.FontSize,
		MainColor:   params.MainColor,
		AccentColor: params.SecColor,
	}
	textbox := NewTextbox(boxParams)
	mw.textbox = textbox
	mw.textbox.SetText(params.Info)
	return mw
}

// NewDialogWindow creates new dialog window with message.
func NewDialogWindow(params Params) *MessageWindow {
	// Basic message window.
	mw := NewMessageWindow(params)
	// Buttons.
	buttonParams := Params{
		Size:      SizeSmall,
		Shape:     ShapeRectangle,
		MainColor: params.SecColor,
	}
	mw.cancelButton = NewButton(buttonParams)
	mw.cancelButton.SetOnClickFunc(mw.onCancelButtonClicked)
	return mw
}

// Draw draws window.
func (mw *MessageWindow) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	mw.drawArea = MatrixToDrawArea(matrix, mw.Size())
	// Background.
	color := mw.color
	if mw.Disabled() {
		color = mw.colorDisable
	}
	DrawRect(t, mw.DrawArea(), color)
	// Buttons.
	acceptButtonPos := MoveBR(mw.Size(), mw.acceptButton.Size())
	mw.acceptButton.Draw(t, matrix.Moved(acceptButtonPos))
	if mw.cancelButton != nil {
		cancelButtonPos := MoveBL(mw.Size(), mw.cancelButton.Size())
		mw.cancelButton.Draw(t, matrix.Moved(cancelButtonPos))
	}
	// Textbox.
	boxMove := MoveTC(mw.Size(), mw.textbox.Size())
	mw.textbox.Draw(t, matrix.Moved(boxMove))
}

// Update handles key press events.
func (mw *MessageWindow) Update(win *Window) {
	if mw.Disabled() {
		return
	}
	if mw.Focused() {
		if win.JustPressed(pixelgl.KeyEscape) {
			mw.cancel()
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			mw.accept()
		}
	}

	if mw.DrawArea().Contains(win.MousePosition()) {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mw.Focus(true)
		}
	}
	mw.textbox.Update(win)
	mw.acceptButton.Update(win)
	if mw.cancelButton != nil {
		mw.cancelButton.Update(win)
	}
}

// Show toggles window visibility.
func (mw *MessageWindow) Show(show bool) {
	mw.opened = show
}

// Active toggles message active state.
func (mw *MessageWindow) Active(active bool) {
	mw.disabled = !active
}

// Focus sets or removes focus from window.
func (mw *MessageWindow) Focus(focus bool) {
	mw.focused = focus
}

// Opened checks whether window should be open.
func (mw *MessageWindow) Opened() bool {
	return mw.opened
}

// Focused checks whether window is focused.
func (mw *MessageWindow) Focused() bool {
	return mw.focused
}

// Dismissed checks whether window was dismised.
func (mw *MessageWindow) Dismissed() bool {
	return mw.dismissed
}

// Accepted checks whether message was accepted.
func (mw *MessageWindow) Accepted() bool {
	return mw.accepted
}

// Disabled checks whether message is unactive.
func (mw *MessageWindow) Disabled() bool {
	return mw.disabled
}

// Size resturns message window size.
func (mw *MessageWindow) Size() pixel.Vec {
	return mw.size
}

// DrawArea returns size of current draw area.
func (mw *MessageWindow) DrawArea() pixel.Rect {
	return mw.drawArea
}

// SetAcceptLabel sets label for accept button.
func (mw *MessageWindow) SetAcceptLabel(l string) {
	mw.acceptButton.SetLabel(l)
}

// SetCancelLabel sets label for cancel button.
func (mw *MessageWindow) SetCancelLabel(l string) {
	if mw.cancelButton == nil {
		return
	}
	mw.cancelButton.SetLabel(l)
}

// SetOnAcceptFunc sets specified function as function triggered after
// message was accepted.
func (mw *MessageWindow) SetOnAcceptFunc(f func(msg *MessageWindow)) {
	mw.onAccept = f
}

// SetOnCancelFunc sets specified function as function triggered after
// message was canceled.
func (mw *MessageWindow) SetOnCancelFunc(f func(msg *MessageWindow)) {
	mw.onCancel = f
}

// Triggered after accept button clicked.
func (mw *MessageWindow) onAcceptButtonClicked(b *Button) {
	if mw.Focused() {
		mw.accept()
	}
}

// Triggered after cancel button clicked.
func (mw *MessageWindow) onCancelButtonClicked(b *Button) {
	if mw.Focused() {
		mw.cancel()
	}
}

// reset resets window to default state(closed, unfocused).
func (mw *MessageWindow) reset() {
	mw.opened = false
	mw.focused = false
}

// accept sets message as accepted.
func (mw *MessageWindow) accept() {
	mw.reset()
	mw.dismissed = true
	mw.accepted = true
	if mw.onAccept != nil {
		mw.onAccept(mw)
	}
}

// cancel sets message as canceled.
func (mw *MessageWindow) cancel() {
	mw.reset()
	mw.dismissed = true
	mw.accepted = false
	if mw.onCancel != nil {
		mw.onCancel(mw)
	}
}
