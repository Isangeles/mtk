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

// Example of creating and using list.
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
		Title:  "MTK list example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail to create mtk window: %v", err))
	}
  // Create list.
  listParams := mtk.Params{
    SizeRaw:     mtk.ConvVec(pixel.V(350, 500)),
    MainColor:   colornames.Grey,
    SecColor:    colornames.Red,
    AccentColor: colornames.Blue,
  }
  list := mtk.NewList(listParams)
  // Insert items to list.
  items := make(map[string]interface{})
  items["Item 1"] = "it1"
  items["Item 2"] = "it2"
  items["Item 3"] = "it3"
  list.InsertItems(items)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
    // Draw.
    listPos := win.Bounds().Center()
    list.Draw(win, mtk.Matrix().Moved(listPos))
		// Update.
		win.Update()
    list.Update(win)
	}
}
