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

// Example of using audio player.
package main

import (
	"fmt"
	"os"
	
	"golang.org/x/image/colornames"
	
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/vorbis"
	
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
		Title:  "MTK audio example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create audio player.
	audioFormat := beep.Format{44100, 2, 2}
	err = mtk.InitAudio(audioFormat)
	if err != nil {
		panic(fmt.Sprintf("Unable to init MTK audio: %v", err))
	}
	audio := mtk.NewAudioPlayer()
	// Load example music.
	music, err := audioBuffer("../res/music.ogg")
	if err != nil {
		panic(fmt.Sprintf("Unable to load example music: %v", err))
	}
	// Play music.
	audio.AddAudio(music)
	audio.ResumePlaylist()
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Update.
		win.Update()
	}
}

// audioBuffer loads vorbis(.ogg) file with specified path.
func audioBuffer(path string) (*beep.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()
	stream, format, err := vorbis.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode vorbis data: %v", err)
	}
	defer stream.Close()
	buffer := beep.NewBuffer(format)
	buffer.Append(stream)
	return buffer, nil
}
