
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

// Example for creating MTK switch.
package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

var (
	exitreq    bool        // menu exit request flag
	exitButton *mtk.Button // button to enable/disable
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
		Title:  "MTK switch example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create mtk window: %v", err))
	}
	// Create switch.
	switchParams := mtk.Params{
		Size:      mtk.SizeBig,
		MainColor: colornames.Grey,
		SecColor:  colornames.Red,
	}
	exitSwitch := mtk.NewSwitch(switchParams)
	// Add values.
	onValue := mtk.SwitchValue{"on", true}
	offValue := mtk.SwitchValue{"off", false}
	switchValues := []mtk.SwitchValue{onValue, offValue}
	exitSwitch.SetValues(switchValues...)
	exitSwitch.SetLabel("Toggle exit button")
	exitSwitch.SetOnChangeFunc(onExitSwitchChanged)
	// Create button to enable/disable.
	buttonParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: colornames.Red,
	}
	exitButton = mtk.NewButton(buttonParams)
	exitButton.SetLabel("Exit")
	exitButton.SetInfo("Exit example")
	exitButton.SetOnClickFunc(onExitButtonClicked)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw switch.
		switchPos := win.Bounds().Center()
		exitSwitch.Draw(win, mtk.Matrix().Moved(switchPos))
		// Draw button.
		buttonPos := mtk.BottomOf(exitSwitch.DrawArea(), exitButton.Size(), 10)
		exitButton.Draw(win, mtk.Matrix().Moved(buttonPos))
		// Update.
		win.Update()
		exitSwitch.Update(win)
		exitButton.Update(win)
		if exitreq {
			win.SetClosed(true)
		}
	}
}

// onExitSwitchChanged handles exit switch change event.
func onExitSwitchChanged(s *mtk.Switch, old, new *mtk.SwitchValue) {
	enabled, ok := new.Value.(bool)
	if !ok {
		fmt.Printf("Invalid switch value type")
		return
	}
	exitButton.Active(enabled)
}

// onExitButtonClicked handles exit button click event.
func onExitButtonClicked(b *mtk.Button) {
	exitreq = true
}
