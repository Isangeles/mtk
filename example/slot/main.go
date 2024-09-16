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

// Example of creating and using slot.
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

var (
	// Colors.
	slotColor      = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	slotColorCheck = pixel.RGBA{255, 0, 0, 0.5}
)

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK slot example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create slot.
	slotParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		MainColor: slotColor,
	}
	slot := mtk.NewSlot(slotParams)
	slot.SetInfo("Click to check or uncheck")
	slot.AddValues(false)
	slot.SetOnLeftClickFunc(onSlotClicked)
	icon, err := loadPicture("icon.png")
	if err != nil {
		panic(fmt.Errorf("Unable to load icon picture: %v", err))
	}
	slot.SetIcon(icon)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		slotPos := win.Bounds().Center()
		slot.Draw(win, mtk.Matrix().Moved(slotPos))
		// Update.
		win.Update()
		slot.Update(win)
	}
}

// onSlotClicked handles slot click event.
func onSlotClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	// Get first slot value.
	value, ok := s.Values()[0].(bool)
	if !ok {
		return
	}
	s.SetValues(nil)
	// Check/uncheck slot.
	if !value {
		s.AddValues(true)
		s.SetColor(slotColorCheck)
		return
	}
	s.AddValues(false)
	s.SetColor(slotColor)
}

// loadPicture loads picture from file with specified path.
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open image file: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image: %v", err)
	}
	return pixel.PictureDataFromImage(img), nil
}
