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
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type ProcArgs struct {
	labels     map[rune]string
	imgDir     string
	keyFilter  []event.Filter
	labelIndex map[rune]int
}

func processArgs() (*ProcArgs, error) {
	args := os.Args[1:]
	pathInfo, err := os.Stat(args[0])
	if err != nil || !pathInfo.IsDir() {
		return nil, err
	}
	labelsMap := make(map[rune]string)
	labelIndex := make(map[rune]int)
	filter := make([]event.Filter, 0, len(args[1:]))
	for _, val := range args[1:] {
		labelCharSplit := strings.Split(val, "=")
		if len(labelCharSplit) != 2 {
			return nil, fmt.Errorf("label and key should be separated by a single '=', %v", val)
		}
		// Need to add validation for label for directory
		// Check character hasn't been added before
		char := unicode.ToUpper([]rune(labelCharSplit[1])[0])
		if _, ok := labelsMap[char]; ok {
			return nil, fmt.Errorf("keyboard character has been assigned twice %v", char)
		}
		labelsMap[char] = labelCharSplit[0]
		labelIndex[char] = 0
		filter = append(filter, key.Filter{Name: key.Name(char)})
	}
	filter = append(filter, key.Filter{Name: key.NameBack})
	return &ProcArgs{imgDir: args[0], labels: labelsMap, keyFilter: filter, labelIndex: labelIndex}, nil
}

func copyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sourceFile)

	// Create the destination file for writing
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(destFile)

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions from source to destination
	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	err = os.Chmod(dst, info.Mode())
	if err != nil {
		return err
	}

	return nil
}

func main() {
	procArgs, err := processArgs()
	if err != nil {
		log.Fatal(err)
	}
	createDirectories(procArgs)
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

func createDirectories(args *ProcArgs) {
	for _, val := range args.labels {
		err := os.Mkdir(args.imgDir+"/"+val, 0750)
		if err != nil {
			log.Fatal(err)
		}
	}
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
						k := []rune(ev.Name)[0]
						label := args.labels[k]
						err := copyFile(args.imgDir+"/"+imagesToView[currentImgIndex].Name(), args.imgDir+"/"+label+"/"+strconv.Itoa(args.labelIndex[k])+".jpg")
						if err != nil {
							return err
						}
						args.labelIndex[k]++
						for {
							currentImgIndex++
							if imagesToView[currentImgIndex].IsDir() {
								continue
							}
						}
					}
				}

				// Process the image
				imgFile, err := os.Open(args.imgDir + "/" + imagesToView[currentImgIndex].Name())
				if err != nil {
					return err
				}
				window.Option(app.Title(imgFile.Name()))

				img, _, err := image.Decode(imgFile)
				if err != nil {
					return err
				}

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
