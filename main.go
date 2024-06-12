package main

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/draw"
	"log"
	"os"
)

func main() {
	go func() {
		window := new(app.Window)
		imageDir := ""
		files, err := os.ReadDir(imageDir)
		if err != nil {
			os.Exit(1)
		}
		for _, imagePath := range files {
			if imagePath.IsDir() {
				continue
			}
			err := run(window, imageDir+"/"+imagePath.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, imagePath string) error {
	theme := material.NewTheme()
	var ops op.Ops
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	err = imgFile.Close()
	if err != nil {
		return err
	}
	imgRGBA, ok := img.(*image.RGBA)
	if !ok {
		b := img.Bounds()
		imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)
	}
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			// Define a large label with an appropriate text
			title := material.H1(theme, "Label to place the image")
			widget.Image{}
			title.Layout(gtx)
		}
	}
}
