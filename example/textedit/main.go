/*
 * main.go
 *
 * Copyright 2022 Dariusz Sikora <dev@isangeles.pl>
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

// Example of creating and using MTK textedit.
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
		Title:  "MTK textedit example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Load & set main UI font.
	font, err := loadFont("SIMSUN.ttf")
	if err != nil {
		// MTK has fallback font, so we don't need to panic.
		fmt.Printf("Unable to load main font: %v\n", err)
	}
	mtk.SetMainFont(font)
	// Create textedit.
	textedit := mtk.NewTextedit(mtk.SizeMedium, colornames.Grey)
	textedit.SetSize(pixel.V(200, 30))
	textedit.SetOnInputFunc(onTexteditInput)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw textedit.
		texteditPos := win.Bounds().Center()
		textedit.Draw(win, mtk.Matrix().Moved(texteditPos))
		// Update.
		win.Update()
		textedit.Update(win)
	}
}

// loadFont reads font file from specified path
// and returns font face or error if file was
// not found.
func loadFont(path string) (*truetype.Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}
	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse font: %v", err)
	}
	return font, nil
}

// Triggered after textedit input.
func onTexteditInput(textedit *mtk.Textedit) {
	fmt.Printf("Input: %s\n", textedit.Text())
}
