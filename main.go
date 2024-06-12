package main

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
	"log"
	"os"
	"strings"
	"unicode"
)

type ProcArgs struct {
	labels    map[rune]string
	imgDir    string
	keyFilter []event.Filter
}

func processArgs() (*ProcArgs, error) {
	args := os.Args[1:]
	pathInfo, err := os.Stat(args[0])
	if err != nil || pathInfo.IsDir() != true {
		return nil, err
	}
	labelsMap := make(map[rune]string)
	filter := make([]event.Filter, 0, len(args[1:]))
	for _, val := range args[1:] {
		labelCharSplit := strings.Split(val, "=")
		if len(labelCharSplit) != 2 {
			return nil, errors.New(fmt.Sprintf("label and key should be separated by a single '=', %v", val))
		}
		// Need to add validation for label for directory
		// Check character hasn't been added before
		char := []rune(labelCharSplit[1])[0]
		if _, ok := labelsMap[char]; ok {
			return nil, errors.New(fmt.Sprintf("keyboard character has been assigned twice %v", char))
		}
		labelsMap[char] = labelCharSplit[0]
		filter = append(filter, key.Filter{Name: key.Name(unicode.ToUpper(char))})
	}
	filter = append(filter, key.Filter{Name: key.NameBack})
	return &ProcArgs{imgDir: args[0], labels: labelsMap, keyFilter: filter}, nil
}

func main() {
	procArgs, err := processArgs()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		window := new(app.Window)

		err := draw(window, procArgs)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func draw(window *app.Window, args *ProcArgs) error {
	//theme := material.NewTheme()
	var ops op.Ops
	imagesToView, _ := os.ReadDir(args.imgDir)
	currentImgIndex := 0
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			if imagesToView[currentImgIndex].IsDir() {
				currentImgIndex++
				continue
			}

			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			// register a global key listener for the escape key wrapping our entire UI.
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			event.Op(gtx.Ops, "tagVal")
			for {
				ev, ok := gtx.Event(args.keyFilter...)
				if !ok {
					return errors.New("KeyFilter error")
				}

				switch ev := ev.(type) {
				case key.Event:
					switch ev.Name {
					case key.NameBack:
						// TODO: Implement backspace
						fmt.Println("Pressed backspace")
					default:
						// Copy the file to the other locations
						fmt.Println("Pressed a labelling key")
						currentImgIndex++
					}
				}

				// Process the image
				imgFile, err := os.Open(args.imgDir + "/" + imagesToView[currentImgIndex].Name())
				if err != nil {
					return err
				}
				window.Option(app.Title(imgFile.Name()))

				img, _, err := image.Decode(imgFile)
				err = imgFile.Close()
				if err != nil {
					return err
				}
				imgOp := paint.NewImageOp(img)
				widget.Image{Src: imgOp, Fit: widget.ScaleDown}.Layout(gtx)

				area.Pop()
				e.Frame(gtx.Ops)
			}
		}
	}
}
