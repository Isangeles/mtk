/*
 * main.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	
	"github.com/isangeles/mtk"
)

var (
	// Exit request flag.
	exitreq bool
	// Colors.
	slotColor      = colornames.Grey
	slotColorCheck = colornames.Red
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
		Title:  "MTK slot example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create button for exit.
	buttonParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: colornames.Red,
	}
	exitButton := mtk.NewButton(buttonParams)
	exitButton.SetLabel("Exit")
	exitButton.SetInfo("Exit menu")
	// Set function for exit button click event.
	exitButton.SetOnClickFunc(onExitButtonClicked)
	// Create slot.
	slotParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		MainColor: slotColor,
	}
	slot := mtk.NewSlot(slotParams)
	slot.AddValues(false)
	slot.SetOnLeftClickFunc(onSlotClicked)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		slotPos := win.Bounds().Center()
		slot.Draw(win, mtk.Matrix().Moved(slotPos))
		exitButtonPos := mtk.BottomOf(slot.DrawArea(), exitButton.Size(), 10)
		exitButton.Draw(win, mtk.Matrix().Moved(exitButtonPos))
		// Update.
		win.Update()
		slot.Update(win)
		exitButton.Update(win)
		// On exit request.
		if exitreq {
			win.SetClosed(true)
		}
	}
}

// onExitButtonClicked handles exit button click
// event.
func onExitButtonClicked(b *mtk.Button) {
	exitreq = true
}

// onSlotClicked handles slot click event.
func onSlotClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	// Get first slot value.
	v, ok := s.Values()[0].(bool)
	if !ok {
		return
	}
	// Check/uncheck slot.
	if v {
		s.Clear()
		s.AddValues(false)
		s.SetColor(slotColor)
	} else {
		s.Clear()
		s.AddValues(true)
		s.SetColor(slotColorCheck)
	}
}
