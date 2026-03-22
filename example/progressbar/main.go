/*
 * main.go
 *
 * Copyright 2026 Dariusz Sikora <ds@isangeles.dev>
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

// Example for MTK Progressbar.
package main

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

// Main function.
func main() {
	// Run Pixel graphic
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration
	cfg := pixelgl.WindowConfig{
		Title:  "MTK window example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create progressbar
	params := mtk.Params{Size: mtk.SizeMedium, MainColor: colornames.Red}
	bar := mtk.NewProgressBar(params)
	bar.SetMax(10)
	update := time.Now()
	var timer int64
	// Main loop
	for !win.Closed() {
		// Clear window
		win.Clear(colornames.Black)
		// Draw progress bar
		barPos := win.Bounds().Center()
		bar.Draw(win, mtk.Matrix().Moved(barPos))
		// Update progress
		if timer >= 10000 && bar.Value() < 10 {
			bar.SetValue(bar.Value() + 1)
			update = time.Now()
			timer = 0
		} else if timer >= 10000 {
			bar.SetValue(0)
			update = time.Now()
			timer = 0
		}
		timer += time.Since(update).Milliseconds()
		// Update
		win.Update()
		bar.Update(win)
	}
}
