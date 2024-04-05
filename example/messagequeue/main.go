/*
 * main.go
 *
 * Copyright 2024 Dariusz Sikora <ds@isangeles.dev>
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

// Example for the MTK Message Queue.
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
		Title:  "MTK message queue example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(winConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create message queue and add some messages.
	messages := mtk.NewMessageQueue(new(mtk.Focus))
	messages.Append(createMessage("Third message"))
	messages.Append(createMessage("Second message"))
	messages.Append(createMessage("First message"))
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		msgPos := win.Bounds().Center()
		messages.Draw(win, mtk.Matrix().Moved(msgPos))
		// Update.
		win.Update()
		messages.Update(win)
	}
}

// createMessage creates new message window with specified
// text as the message content.
func createMessage(text string) *mtk.MessageWindow {
	params := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		MainColor: colornames.Grey,
		SecColor:  colornames.Red,
		Info:      text,
	}
	message := mtk.NewMessageWindow(params)
	message.SetAcceptLabel("OK")
	message.Show(true)
	return message
}
