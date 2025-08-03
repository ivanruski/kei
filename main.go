package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Main text viewer.
	view := tview.NewTextView().SetDynamicColors(true)
	view.SetBorder(true).
		SetTitle(" kubectl explain ")

	// Input field for `:`
	explainInput := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	explainInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			path := strings.TrimSpace(explainInput.GetText())
			view.SetText(runExplain(path))
			app.SetFocus(view)

		case tcell.KeyEscape:
			app.SetFocus(view)
		}
	})

	// Input field for `/`
	searchExplainOutputInput := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	searchExplainOutputInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			searchTerm := searchExplainOutputInput.GetText()
			if searchTerm == "" {
				app.SetFocus(view)
				return
			}

			lines := strings.Split(view.GetText(true), "\n")
			var highlighted strings.Builder
			for _, line := range lines {
				for {
					idx := strings.Index(line, searchTerm)
					if idx == -1 {
						highlighted.WriteString(line)
						break
					}

					highlighted.WriteString(line[:idx])
					highlighted.WriteString("[black:white]")
					highlighted.WriteString(searchTerm)
					highlighted.WriteString("[white:black]")

					line = line[idx+len(searchTerm):]
				}

				highlighted.WriteString("\n")
			}

			view.SetText(highlighted.String())
			app.SetFocus(view)

		case tcell.KeyEscape:
			app.SetFocus(view)
		}
	})

	inputPages := tview.NewPages().
		AddPage("explain", explainInput, true, false).
		AddPage("search", searchExplainOutputInput, true, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 0, 1, true).       // height is flexible
		AddItem(inputPages, 1, 0, false) // height is 1 line, a the bottom

	view.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case ':': // Focus input field to type <type>.<fieldName>[.<fieldName>]
				inputPages.SwitchToPage("explain")
				explainInput.SetLabel(":")
				app.SetFocus(explainInput)
				return nil
			case '/': // Focus input field to type string to search in explain output buffer
				inputPages.SwitchToPage("search")
				searchExplainOutputInput.SetLabel("/")
				app.SetFocus(searchExplainOutputInput)
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
