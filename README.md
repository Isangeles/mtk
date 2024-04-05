## Introduction
MTK([Mural](https://github.com/Isangeles/mural) Toolkit) is a simple widget toolkit for creating graphical user interfaces written in Go with [Pixel](https://github.com/gopxl/pixel) library.

The toolkit provides basic UI elements(buttons, text boxes, switches, lists, animations, etc.), simple audio handling and automatic scaling of UI elements.

Originally created as a part of [Mural](https://github.com/Isangeles/mural) GUI.

MTK highly relies on [Pixel](https://github.com/gopxl/pixel) library, make sure to check [Pixel wiki](https://github.com/gopxl/pixel/wiki) first.

## Dependencies
Basic dependencies are OpenGL development libraries, and some audio development libraries.

### Linux
On Fedora-like distribution install: `go` `libX11-devel` `libXcursor-devel` `libXrandr-devel` `libXinerama-devel` `mesa-libGL-devel` `libXi-devel` `libXxf86vm-devel` `alsa-lib-devel`.

On other distributions, you need to install the equivalence of these packages.
### macOS
Install [Go](https://go.dev/) and Xcode or Command Line Tools for Xcode.
### Windows
Install [Go](https://go.dev/)

## Examples
All UI elements are automatically scaled to size specified in Pixel WindowConfig.

Create window:
```
// Create Pixel window configuration.
cfg := pixelgl.WindowConfig{
       Title:  "MTK window",
       Bounds: pixel.R(0, 0, 1600, 900),
}
// Create MTK warpper for Pixel window.
win, err := mtk.NewWindow(cfg)
if err != nil {
       panic(fmt.Errorf("Unable to create MTK window: %v", err))
}
```
Create button:
```
// Specify button parameters.
params := mtk.Params{
	Size:      mtk.SizeBig,
	FontSize:  mtk.SizeMedium,
	Shape:     mtk.ShapeRectangle,
	MainColor: colornames.Red,
}
// Create button.
button := mtk.NewButton(params)
// Set label and on-hover info.
button.SetLabel("Button")
button.SetInfo("Button info")
// Set some function on button click event.
button.SetOnClickFunc(onButtonClickedFunc)
```
Draw button in window:
```
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
```
Check [example](https://github.com/Isangeles/mtk/tree/master/example) package for more detailed examples.

## Documentation
Source code documentation can be easily browsed with `go doc` command.

## Contributing
You are welcome to contribute to project development.

If you looking for things to do, then check [TODO file](https://github.com/Isangeles/mtk/blob/master/TODO) or contact maintainer(ds@isangeles.dev).

When you find something to do, create new branch for your feature.
After you finish, open pull request to merge your changes with master branch.

## Contact
* Isangeles <<ds@isangeles.dev>>

## License
Copyright 2018-2024 Dariusz Sikora <<ds@isangeles.dev>>

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
MA 02110-1301, USA.
