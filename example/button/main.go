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

// Example for creating simple MTK button with draw background,
// custom on-click function, and click sound effect.
package main

import (
	"fmt"
	"os"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"

	"github.com/isangeles/mtk"
)

var (
	exitreq bool // menu exit request flag
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
		Title:  "MTK button example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Initialize MTK audio player(for button click sound).
	audioFormat := beep.Format{44100, 2, 2}
	mtk.Audio, err = mtk.NewAudioPlayer(audioFormat)
	if err != nil {
		panic(fmt.Errorf("Unable to create audio player: %v", err))
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
	// Set button click sound effect.
	soundEffect, err := audioBuffer("res/click.flac")
	if err != nil {
		panic(fmt.Errorf("Unable to load click sound effect: %v", err))
	}
	exitButton.SetClickSound(soundEffect)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw exit button.
		exitButtonPos := win.Bounds().Center()
		exitButton.Draw(win, mtk.Matrix().Moved(exitButtonPos))
		// Update.
		win.Update()
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

// audioBuffer loads flac audio file with specified path.
func audioBuffer(path string) (*beep.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()
	stream, format, err := flac.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("unable to decode flac data: %v", err)
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(stream)
	return buffer, nil
}
