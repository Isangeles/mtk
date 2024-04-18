/*
 * mtk.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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

// Toolkit package for Mural GUI.
package mtk

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/golang/freetype/truetype"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/text"
)

const (
	// Sizes.
	SizeMini Size = iota
	SizeSmall
	SizeMedium
	SizeBig
	SizeHuge
	// Shapes.
	ShapeRectangle Shape = iota
	ShapeSquare
	// Aligns
	AlignCenter Align = iota
	AlignRight
	AlignLeft
)

var (
	// Toolkit audio player used to play various sound effects,
	// like button click sound for example.
	audio            *AudioPlayer = NewAudioPlayer()
	buttonClickSound *beep.Buffer
	// Font.
	fallbackFont font.Face = basicfont.Face7x13
	mainFontBase *truetype.Font
	// Time.
	secTimer = time.Tick(time.Second)
)

// Type for shapes of UI elements.
// Shapes: rectangle(0), square(1).
type Shape int

// Type for sizes of UI elements, like buttons, switches, etc.
// Sizes: mini(0), small(1), medium(2), big(3), huge(4).
type Size int

// Type for aligns.
// Directions: center(0), right(1), left(2)
type Align int

// Interface for all 'focusable' UI elements, like buttons,
// switches, etc.
type Focuser interface {
	Focus(focus bool)
	Focused() bool
}

// Focus represents user focus on UI element.
type Focus struct {
	element Focuser
}

// Focus sets focus on specified focusable element.
// Previously focused element(if exists) is unfocused before
// focusing specified one.
func (f *Focus) Focus(e Focuser) {
	if f.element != nil {
		f.element.Focus(false)
	}
	f.element = e
	if e == nil {
		return
	}
	f.element.Focus(true)
}

// ButtonSize returns szie parameters for button with
// this size and with specifed shape.
func (s Size) ButtonSize(sh Shape) pixel.Vec {
	switch {
	case s <= SizeMini && sh == ShapeSquare:
		return ConvVec(pixel.V(30, 30))
	case s == SizeSmall && sh == ShapeSquare:
		return ConvVec(pixel.V(45, 45))
	case s == SizeMedium && sh == ShapeSquare:
		return ConvVec(pixel.V(60, 60))
	case s >= SizeBig && sh == ShapeSquare:
		return ConvVec(pixel.V(70, 70))
	case s <= SizeMini && sh == ShapeRectangle:
		return ConvVec(pixel.V(30, 15))
	case s == SizeSmall && sh == ShapeRectangle:
		return ConvVec(pixel.V(70, 35))
	case s == SizeMedium && sh == ShapeRectangle:
		return ConvVec(pixel.V(100, 50))
	case s >= SizeBig && sh == ShapeRectangle:
		return ConvVec(pixel.V(120, 70))
	default:
		return ConvVec(pixel.V(70, 35))
	}
}

// SwitchSize return rectangele parameters for switch
// with this size.
func (s Size) SwitchSize() pixel.Vec {
	switch {
	case s <= SizeSmall:
		return ConvVec(pixel.V(170, 50))
	case s == SizeMedium:
		return ConvVec(pixel.V(210, 80))
	case s == SizeBig:
		return ConvVec(pixel.V(260, 140))
	case s >= SizeHuge:
		return ConvVec(pixel.V(270, 130))
	default:
		return ConvVec(pixel.V(70, 35))
	}
}

// MessageWindowSize returns size parameters for message window.
func (s Size) MessageWindowSize() pixel.Vec {
	switch {
	case s <= SizeSmall:
		return ConvVec(pixel.V(400, 300))
	default:
		return ConvVec(pixel.V(400, 300))
	}
}

// ListSize returns size parameters for list.
func (s Size) ListSize() pixel.Vec {
	switch {
	case s <= SizeMedium:
		return ConvVec(pixel.V(600, 500))
	default:
		return ConvVec(pixel.V(600, 500))
	}
}

// BarSize returns size parameters for progress bar.
func (s Size) BarSize() pixel.Vec {
	switch {
	case s <= SizeMini:
		return ConvVec(pixel.V(100, 10))
	default:
		return ConvVec(pixel.V(100, 10))
	}
}

// SlotSize returns size parameters for slot.
func (s Size) SlotSize() pixel.Vec {
	switch {
	case s <= SizeSmall:
		return ConvVec(pixel.V(25, 25))
	case s == SizeMedium:
		return ConvVec(pixel.V(30, 30))
	case s >= SizeBig:
		return ConvVec(pixel.V(40, 40))
	default:
		return ConvVec(pixel.V(25, 25))
	}
}

// InitAudio initializes the system audio with specified format.
// It need to be called before using the audio player struct.
func InitAudio(format beep.Format) error {
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("Unable to init the system speaker: %v", err)
	}
	return nil
}

// Audio returns toolkit audio player.
// Returned audio player will not be operational if the audio
// was not initialized before with InitAudio function.
func Audio() *AudioPlayer {
	return audio
}

// Sets specified truetype font as current main font of
// the interface.
func SetMainFont(font *truetype.Font) {
	mainFontBase = font
}

// SetButtonClickSound sets specified audio buffer
// as on-click audio effect for all buttons.
func SetButtonClickSound(s *beep.Buffer) {
	buttonClickSound = s
}

// MainFont returns main font in specified size from
// data package.
func MainFont(s Size) font.Face {
	switch {
	case s <= SizeMini:
		return createMainFont(10)
	case s == SizeSmall:
		return createMainFont(15)
	case s == SizeMedium:
		return createMainFont(20)
	case s >= SizeBig:
		return createMainFont(30)
	default:
		return createMainFont(10)
	}
}

// Atlas returns atlas for UI text with specified
// font.
func Atlas(f *font.Face) *text.Atlas {
	return text.NewAtlas(*f, text.ASCII)
}

// Matrix return scaled identity matrix.
func Matrix() pixel.Matrix {
	return pixel.IM.Scaled(pixel.V(0, 0), Scale())
}

// DrawRectangle draw rectangle on specified target with
// specified draw area(position and size) and color.
func DrawRectangle(t pixel.Target, drawArea pixel.Rect, color color.Color) {
	draw := imdraw.New(nil)
	draw.Clear()
	draw.Color = color
	draw.Push(drawArea.Min)
	draw.Push(drawArea.Max)
	draw.Rectangle(0)
	draw.Draw(t)
}

// createMainFont creates new main font face with
// specified size.
func createMainFont(size float64) font.Face {
	if mainFontBase == nil {
		return fallbackFont
	}
	return truetype.NewFace(mainFontBase, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}
