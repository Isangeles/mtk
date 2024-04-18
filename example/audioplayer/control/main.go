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

// Example of controlling audio player.
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
		Title:  "MTK audio control example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Control buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: colornames.Red,
	}
	playButton := mtk.NewButton(buttonParams)
	playButton.SetLabel("Play")
	playButton.SetOnClickFunc(onPlayButtonClicked)
	stopButton := mtk.NewButton(buttonParams)
	stopButton.SetLabel("Stop")
	stopButton.SetOnClickFunc(onStopButtonClicked)
	nextButton := mtk.NewButton(buttonParams)
	nextButton.SetLabel("Next")
	nextButton.SetOnClickFunc(onNextButtonClicked)
	prevButton := mtk.NewButton(buttonParams)
	prevButton.SetLabel("Previous")
	prevButton.SetOnClickFunc(onPrevButtonClicked)
	upButton := mtk.NewButton(buttonParams)
	upButton.SetLabel("+")
	upButton.SetOnClickFunc(onUpButtonClicked)
	downButton := mtk.NewButton(buttonParams)
	downButton.SetLabel("-")
	downButton.SetOnClickFunc(onDownButtonClicked)
	muteButton := mtk.NewButton(buttonParams)
	muteButton.SetLabel("Mute")
	muteButton.SetOnClickFunc(onMuteButtonClicked)
	// Track, status, volume info.
	textParams := mtk.Params{
		Size:     mtk.SizeMedium,
		FontSize: mtk.SizeBig,
	}
	trackInfo := mtk.NewText(textParams)
	statusInfo := mtk.NewText(textParams)
	volInfo := mtk.NewText(textParams)
	// Init MTK audio.
	audioFormat := beep.Format{44100, 2, 2}
	err = mtk.InitAudio(audioFormat)
	if err != nil {
		panic(fmt.Sprintf("Unable to init MTK audio: %v", err))
	}
	// Create audio player and load example music.
	audio = mtk.NewAudioPlayer()
	music, err := audioBuffer("../res/music.ogg")
	if err != nil {
		panic(fmt.Sprintf("Unable to load example music: %v", err))
	}
	audio.AddAudio(music)
	music, err = audioBuffer("../res/music2.ogg")
	if err != nil {
		panic(fmt.Sprintf("Unable to load exmaple music: %v", err))
	}
	audio.AddAudio(music)
	// Play music.
	audio.ResumePlaylist()
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		playButtonPos := win.Bounds().Center()
		stopButtonPos := mtk.LeftOf(playButton.DrawArea(), stopButton.Size(), 10)
		nextButtonPos := mtk.RightOf(playButton.DrawArea(), nextButton.Size(), 10)
		prevButtonPos := mtk.LeftOf(stopButton.DrawArea(), prevButton.Size(), 10)
		muteButtonPos := mtk.BottomOf(playButton.DrawArea(), muteButton.Size(), 10)
		upButtonPos := mtk.RightOf(muteButton.DrawArea(), upButton.Size(), 10)
		downButtonPos := mtk.LeftOf(muteButton.DrawArea(), downButton.Size(), 10)
		volInfoPos := mtk.TopOf(playButton.DrawArea(), volInfo.Size(), 10)
		statusInfoPos := mtk.TopOf(volInfo.DrawArea(), statusInfo.Size(), 10)
		trackInfoPos := mtk.TopOf(statusInfo.DrawArea(), trackInfo.Size(), 10)
		playButton.Draw(win, mtk.Matrix().Moved(playButtonPos))
		stopButton.Draw(win, mtk.Matrix().Moved(stopButtonPos))
		nextButton.Draw(win, mtk.Matrix().Moved(nextButtonPos))
		prevButton.Draw(win, mtk.Matrix().Moved(prevButtonPos))
		muteButton.Draw(win, mtk.Matrix().Moved(muteButtonPos))
		upButton.Draw(win, mtk.Matrix().Moved(upButtonPos))
		downButton.Draw(win, mtk.Matrix().Moved(downButtonPos))
		volInfo.Draw(win, mtk.Matrix().Moved(volInfoPos))
		statusInfo.Draw(win, mtk.Matrix().Moved(statusInfoPos))
		trackInfo.Draw(win, mtk.Matrix().Moved(trackInfoPos))
		// Update.
		win.Update()
		playButton.Update(win)
		stopButton.Update(win)
		nextButton.Update(win)
		prevButton.Update(win)
		muteButton.Update(win)
		upButton.Update(win)
		downButton.Update(win)
		statusInfo.SetText(fmt.Sprintf("Playing: %v", audio.Playing()))
		trackInfo.SetText(fmt.Sprintf("Track: %d", audio.PlayIndex()))
		volText := fmt.Sprintf("%f", audio.Volume())
		if audio.Muted() {
			volText = "muted"
		}
		volInfo.SetText(fmt.Sprintf("Volume: %s", volText))
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

// Triggered on play button click event.
func onPlayButtonClicked(b *mtk.Button) {
	audio.ResumePlaylist()
}

// Triggered on stop button click event.
func onStopButtonClicked(b *mtk.Button) {
	audio.Stop()
}

// Triggered on next button click event.
func onNextButtonClicked(b *mtk.Button) {
	audio.Stop()
	audio.SetPlayIndex(audio.PlayIndex()+1)
}

// Triggered on previous button click event.
func onPrevButtonClicked(b *mtk.Button) {
	audio.Stop()
	audio.SetPlayIndex(audio.PlayIndex()-1)
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
