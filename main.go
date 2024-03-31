package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("creating new screen: %v", err)
	}

	err = screen.Init()
	if err != nil {
		log.Fatalf("initializing screen: %v", err)
	}
	defer screen.Fini()

	screen.DisableMouse()
	screen.SetStyle(tcell.StyleDefault)

	enterCommandMode(screen)

	var (
		helpText = []string{
			"type  /<type>[.<fieldName>]<Enter> to see an explaination of a Kubernetes resource",
			"type  u                            to scroll down one half of the screen size",
			"type  u                            to scroll up one half of the screen size",
			"type  q<Enter>                     to exit",

			"type  h                            to see this message again",
		}
		displayText = helpText
		offset      int
		_, sheight  = screen.Size()
		height      = sheight - 1 // displayable height
		resource    string
	)

	for {
		screen.Clear()
		enterCommandMode(screen)
		display(screen, displayText[offset:])

		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			// TODO: set screen size
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Rune() {
			case 'q':
				return
			case 'h':
				displayText = helpText
				offset = 0
			case '/':
				enterLastLineMode(screen, resource)

				input := receiveInput(screen, resource)
				enterCommandMode(screen)
				if len(input) == 0 {
					break
				}

				expl, err := getExplanation(input)
				if err != nil {
					displayText = []string{err.Error()}
				} else {
					displayText = expl
					resource = input
				}

				offset = 0

			case 'd':
				no, ok := scrollDownText(displayText, offset, height)

				if !ok {
					_ = screen.Beep()
				} else {
					offset = no
				}
			case 'u':
				no, ok := scrollUpText(displayText, offset, height)

				if !ok {
					_ = screen.Beep()
				} else {
					offset = no
				}
			default:
				screen.Beep()
			}
		}
	}
}

func scrollDownText(lines []string, offset, height int) (int, bool) {
	if len(lines[offset:]) <= height {
		return 0, false
	}

	// check if you are about to display less than the whole screen
	mid := int(math.Round(float64(height) / float64(2)))
	offset += mid
	if len(lines[offset:]) <= height {
		offset = len(lines) - height
	}

	return offset, true
}

func scrollUpText(lines []string, offset, height int) (int, bool) {
	if offset == 0 {
		return 0, false
	}

	mid := int(math.Round(float64(height) / float64(2)))
	offset -= mid
	if offset < 0 {
		offset = 0
	}

	return offset, true
}

func display(screen tcell.Screen, lines []string) {
	width, height := screen.Size()
	limit := len(lines)
	if limit >= height-1 {
		limit = height - 1
	}

	for i := 0; i < limit; i++ {
		drawText(screen, 0, i, width, height, tcell.StyleDefault, lines[i])
	}
}

func logS(screen tcell.Screen, msg string) {
	width, height := screen.Size()
	b := strings.Builder{}
	for i := 0; i < width; i++ {
		b.WriteString(" ")
	}

	drawText(screen, 0, 0, width, height, tcell.StyleDefault, b.String())
	drawText(screen, 0, 0, width, height, tcell.StyleDefault, msg)
}

func enterCommandMode(screen tcell.Screen) {
	screen.Clear()
	width, height := screen.Size()
	drawText(screen, 0, height-1, width, height, tcell.StyleDefault, string(":"))
	screen.ShowCursor(1, height-1)
}

func enterLastLineMode(screen tcell.Screen, resource string) {
	width, height := screen.Size()

	res := fmt.Sprintf("/%s", resource)
	drawText(screen, 0, height-1, width, height, tcell.StyleDefault, res)
	screen.ShowCursor(len(res), height-1)
}

func receiveInput(screen tcell.Screen, resource string) string {
	input := []rune(resource)

	width, height := screen.Size()
	col := len(input) + 1
	for {
		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			width, height = screen.Size()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return string(input)
			case tcell.KeyBackspace2:
				if len(input) == 0 {
					return ""
				}

				input = input[:len(input)-1]
				col--
				drawText(screen, col, height-1, width, height, tcell.StyleDefault, string(" "))
				screen.ShowCursor(col, height-1)
			}

			r := ev.Rune()

			if unicode.IsLetter(r) || r == '.' {
				drawText(screen, col, height-1, width, height, tcell.StyleDefault, string(r))
				col++
				input = append(input, r)
				screen.ShowCursor(col, height-1)
			}
		}
	}
}

func drawText(screen tcell.Screen, x, y, width, height int, style tcell.Style, text string) {
	row := y
	col := x

	for _, r := range []rune(text) {
		screen.SetContent(col, row, r, nil, style)

		col++
		if col >= width {
			col = x
			row++
		}

		if row > height {
			break
		}
	}
}

func getExplanation(resource string) ([]string, error) {
	cmd := exec.Command("kubectl", "explain", resource)

	out, err := cmd.Output()
	if err != nil {
		ee := err.(*exec.ExitError)
		return nil, errors.New(string(ee.Stderr))
	}

	return splitByLines(out)
}

func splitByLines(buf []byte) ([]string, error) {
	outlines := []string{}
	scan := bufio.NewScanner(bytes.NewReader(buf))
	scan.Split(bufio.ScanLines)
	for scan.Scan() {
		outlines = append(outlines, scan.Text())
	}

	if scan.Err() != nil {
		return nil, fmt.Errorf("scan: %w", scan.Err())
	}

	return outlines, nil
}
