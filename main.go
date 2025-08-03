package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func runExplain(path string) string {
	cmd := exec.Command("kubectl", "explain", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return out.String()
		}

		return err.Error()
	}

	return out.String()
}

func main() {
	app := tview.NewApplication()

	// Main text viewer.
	view := tview.NewTextView()
	view.SetBorder(true).
		SetTitle(" kubectl explain ")

	// Input field for `:`
	explainInput := tview.NewInputField()

	explainInput.
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				path := strings.TrimSpace(explainInput.GetText())
				view.SetText(runExplain(path))
				app.SetFocus(view)

			case tcell.KeyEscape:
				app.SetFocus(view)
			}
		})

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 0, 1, true).         // height is flexible
		AddItem(explainInput, 1, 0, false) // height is 1 line, a the bottom

	view.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case ':': // Show input
				explainInput.SetLabel(":")
				app.SetFocus(explainInput)
				return nil
			case 'q': // Quit
				app.Stop()
				return nil
			}
		}
		return ev
	})

	// Show default help or root doc
	view.SetText(runExplain("--help"))

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
