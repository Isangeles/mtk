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

// Example of creating and using slot list.
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/mtk"
)

var slots *mtk.SlotList

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK slot list example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to create MTK window: %v", err))
	}
	// Create slot list.
	slotsSize := mtk.ConvVec(pixel.V(250, 300))
	slots = mtk.NewSlotList(slotsSize, colornames.Grey, mtk.SizeMedium)
	// Create and add slots.
	for i := 0; i < 30; i++ {
		params := mtk.Params{
			Size:     mtk.SizeBig,
			FontSize: mtk.SizeMedium,
		}
		slot := mtk.NewSlot(params)
		slot.SetOnLeftClickFunc(dragSlot)
		slot.SetOnRightClickFunc(dropSlot)
		slots.Add(slot)
	}
	// Insert some content into fist two slots.
	slot := slots.EmptySlot()
	slot.AddValues("value 1")
	icon, err := loadPicture("icon1.png")
	if err != nil {
		panic(fmt.Errorf("Unable to load icon 1"))
	}
	slot.SetIcon(icon)
	slot.SetLabel("Slot 1")
	slot.SetInfo("Use left click to drag and\n right click to drop")
	slot = slots.EmptySlot()
	slot.AddValues("value 2")
	icon, err = loadPicture("icon2.png")
	if err != nil {
		panic(fmt.Errorf("Unable to load icon 2"))
	}
	slot.SetIcon(icon)
	slot.SetInfo("Use left click to drag and\n right click to drop")
	slot.SetLabel("Slot 2")
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		slotsPos := win.Bounds().Center()
		slots.Draw(win, mtk.Matrix().Moved(slotsPos))
		// Update.
		win.Update()
		slots.Update(win)
	}
}

// dragSlot sets drag mode for specified slot and disables
// drag mode on any other slot if enabled.
func dragSlot(slot *mtk.Slot) {
	if len(slot.Values()) < 1 {
		return
	}
	draggedSlot := draggedSlot()
	if draggedSlot != nil {
		draggedSlot.Drag(false)
	}
	slot.Drag(true)
}

// dropSlot drops currently dragged slot into specified slot
// by switching the contents between two slots.
func dropSlot(slot *mtk.Slot) {
	draggedSlot := draggedSlot()
	if draggedSlot != nil {
		mtk.SlotSwitch(slot, draggedSlot)
	}
}

// draggedSlot returns first currently dragged slot
// from the slot list.
func draggedSlot() *mtk.Slot {
	for _, slot := range slots.Slots() {
		if slot.Dragged() {
			return slot
		}
	}
	return nil
}

// loadPicture loads picture from file with specified path.
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open image file: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image: %v", err)
	}
	return pixel.PictureDataFromImage(img), nil
}
