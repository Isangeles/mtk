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

// Example of controlling audio player volume.
package main

import (
	"fmt"
	"os"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/beep"
	"github.com/faiface/beep/vorbis"

	"github.com/isangeles/mtk"
)

var audio *mtk.AudioPlayer

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK audio volume example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Volume buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: colornames.Red,
	}
	upButton := mtk.NewButton(buttonParams)
	upButton.SetLabel("+")
	upButton.SetOnClickFunc(onUpButtonClicked)
	downButton := mtk.NewButton(buttonParams)
	downButton.SetLabel("-")
	downButton.SetOnClickFunc(onDownButtonClicked)
	muteButton := mtk.NewButton(buttonParams)
	muteButton.SetLabel("Mute")
	muteButton.SetOnClickFunc(onMuteButtonClicked)
	// Volume info.
	textParams := mtk.Params{
		Size:     mtk.SizeMedium,
		FontSize: mtk.SizeBig,
	}
	volInfo := mtk.NewText(textParams)
	// Init MTK audio.
	audioFormat := beep.Format{44100, 2, 2}
	audio, err = mtk.NewAudioPlayer(audioFormat)
	if err != nil {
		panic(fmt.Sprintf("Unable to create audio player: %v", err))
	}
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
		// Draw.
		muteButtonPos := win.Bounds().Center()
		upButtonPos := mtk.RightOf(muteButton.DrawArea(), upButton.Size(), 10)
		downButtonPos := mtk.LeftOf(muteButton.DrawArea(), downButton.Size(), 10)
		volInfoPos := mtk.TopOf(muteButton.DrawArea(), volInfo.Size(), 10)
		muteButton.Draw(win, mtk.Matrix().Moved(muteButtonPos))
		upButton.Draw(win, mtk.Matrix().Moved(upButtonPos))
		downButton.Draw(win, mtk.Matrix().Moved(downButtonPos))
		volInfo.Draw(win, mtk.Matrix().Moved(volInfoPos))
		// Update.
		win.Update()
		muteButton.Update(win)
		upButton.Update(win)
		downButton.Update(win)
		if audio.Muted() {
			volInfo.SetText("Volume: muted")
		} else {
			volInfo.SetText(fmt.Sprintf("Volume: %f", audio.Volume()))
		}
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
	buffer := beep.NewBuffer(format)
	buffer.Append(stream)
	return buffer, nil
}

// Triggered on up button click event.
func onUpButtonClicked(b *mtk.Button) {
	audio.SetVolume(audio.Volume() + 1)
}

// Triggered on down button click event.
func onDownButtonClicked(b *mtk.Button) {
	audio.SetVolume(audio.Volume() - 1)
}

// Triggered on mute button click event.
func onMuteButtonClicked(b *mtk.Button) {
	audio.SetMute(!audio.Muted())
}
