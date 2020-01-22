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

// Example of creating and using MTK textbox.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/image/colornames"

	"github.com/golang/freetype/truetype"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

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
		Title:  "MTK textbox example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail_to_create_mtk_window:%v", err))
	}
	// Load & set main UI font.
	font, err := loadFont("SIMSUN.ttf")
	if err != nil {
		// MTK has fallback font, so we don't need to panic.
		fmt.Printf("fail_to_load_main_font:%v\n", err)
	}
	mtk.SetMainFont(font)
	// Create textbox.
	textboxParams := mtk.Params{
		SizeRaw:     pixel.V(300, 300),
		FontSize:    mtk.SizeMedium,
		MainColor:   colornames.Grey,
		AccentColor: colornames.Red,
	}
	textbox := mtk.NewTextbox(textboxParams)
	textbox.Focus(true)
	// Insert text to textbox.
	text := `This is multi-line text,
you can scroll it with
UP
and
DOWN
keys!
`
	textbox.SetText(text)
	textbox.AddText("veeeeeeeeeeeeeeery veeeeeeeeery loooooooooooooooong linnnnnnnnnnnnnnnnnne? No problem!\n")
	textbox.AddText("line 1\n")
	textbox.AddText("line 2\n")
	textbox.AddText("line 3\n")
	textbox.AddText("line 4\n")
	textbox.AddText("line 5\n")
	textbox.AddText("line 6\n")
	textbox.AddText("line 7\n")
	textbox.AddText("line 8\n")
	textbox.AddText("line 9\n")
	textbox.AddText("line 10\n")
	textbox.AddText("line 11\n")
	textbox.AddText("line 12\n")
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw textbox.
		textboxPos := win.Bounds().Center()
		textbox.Draw(win, mtk.Matrix().Moved(textboxPos))
		// Update.
		win.Update()
		textbox.Update(win) // update made scrolling possible
	}
}

// loadFont reads font file from specified path
// and returns font face or error if file was
// not found.
func loadFont(path string) (*truetype.Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_file:%v", err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_file:%v", err)
	}
	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_font:%v", err)
	}
	return font, nil
}
