/*
 * infowindow.go
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
)

// InfoWindow struct for small text boxes that
// follows mouse cursor.
type InfoWindow struct {
	*Text
	bgColor  color.Color
	drawArea pixel.Rect
}

// NewInfoWindow creates new information window.
func NewInfoWindow(params Params) *InfoWindow {
	iw := new(InfoWindow)
	textParams := Params{
		FontSize: params.FontSize,
	}
	iw.Text = NewText(textParams)
	iw.bgColor = colornames.Black
	if params.MainColor != nil {
		iw.bgColor = params.MainColor
	}
	return iw
}

// Draw draws info window.
func (iw *InfoWindow) Draw(t pixel.Target) {
	DrawRectangle(t, iw.DrawArea(), iw.bgColor)
	textPos := pixel.V(iw.drawArea.Center().X, iw.drawArea.Center().Y-iw.Size().Y/2)
	iw.Text.Draw(t, Matrix().Moved(textPos))
}

// Update updates info window.
func (iw *InfoWindow) Update(win *Window) {
	iw.drawArea = pixel.R(win.MousePosition().X, win.MousePosition().Y,
		win.MousePosition().X+iw.Size().X,
		win.MousePosition().Y+iw.Size().Y*1.5)
}
