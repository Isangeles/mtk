/*
 * infowindow.go
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
	"image/color"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

// InfoWindow struct for small text boxes that
// follows mouse cursor.
type InfoWindow struct {
	*Text
	draw     *imdraw.IMDraw
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
	iw.draw = imdraw.New(nil)
	iw.bgColor = params.MainColor
	return iw
}

// Draw draws info window.
func (iw *InfoWindow) Draw(t pixel.Target) {
	iw.drawBackground(t)
	iw.Text.Draw(t, Matrix().Moved(iw.drawArea.Center()))
}

// Update updates info window.
func (iw *InfoWindow) Update(win *Window) {
	iw.drawArea = pixel.R(win.MousePosition().X, win.MousePosition().Y,
		win.MousePosition().X+iw.Size().X,
		win.MousePosition().Y+iw.Size().Y*1.5)
}

// drawBackground draws info background.
func (iw *InfoWindow) drawBackground(t pixel.Target) {
	iw.draw.Clear()
	iw.draw.Color = iw.bgColor
	iw.draw.Push(iw.DrawArea().Min)
	iw.draw.Push(iw.DrawArea().Max)
	iw.draw.Rectangle(0)
	iw.draw.Draw(t)
}
