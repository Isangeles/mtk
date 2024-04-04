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

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"

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
	// Create the button.
	buttonParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: colornames.Red,
	}
	button := mtk.NewButton(buttonParams)
	button.SetLabel("Click")
	button.SetInfo("Click me!")
	// Set function for button click event.
	button.SetOnClickFunc(onButtonClicked)
	// Set button click sound effect.
	soundEffect, err := audioBuffer("res/click.flac")
	if err != nil {
		panic(fmt.Errorf("Unable to load click sound effect: %v", err))
	}
	button.SetClickSound(soundEffect)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw button.
		buttonPos := win.Bounds().Center()
		button.Draw(win, mtk.Matrix().Moved(buttonPos))
		// Update.
		win.Update()
		button.Update(win)
	}
}

// onButtonClicked handles button click
// event.
func onButtonClicked(b *mtk.Button) {
	fmt.Println("Click!")
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
