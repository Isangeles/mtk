/*
 * main.go
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

// Example for creating MTK message window.
package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	winConfig := pixelgl.WindowConfig{
		Title:  "MTK message window example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK wrapper for Pixel window.
	win, err := mtk.NewWindow(winConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create info message.
	messageParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		MainColor: colornames.Grey,
		SecColor:  colornames.Red,
		Info:      "Click OK to close this message",
	}
	message := mtk.NewMessageWindow(messageParams)
	message.SetAcceptLabel("OK")
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		messagePos := win.Bounds().Center()
		if message.Opened() {
			message.Draw(win, mtk.Matrix().Moved(messagePos))
		}
		// Update.
		win.Update()
		message.Update(win)
	}
}
