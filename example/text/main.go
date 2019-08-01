/*
 * main.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

// Example of creating text to draw.
package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

// Main function.
func main () {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK text example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail to create mtk window: %v", err))
	}
	// Create text.
	text := mtk.NewText(mtk.SizeHuge, 100)
	//fmt.Fprintf(text, )
	text.SetText("Hello MTK!\nnew line\ntooo loooooooooong linnnnnne")
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw text.
		textPos := win.Bounds().Center()
		text.Draw(win, mtk.Matrix().Moved(textPos))
		// Update.
		win.Update()
	}
}
