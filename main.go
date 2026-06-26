package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Selection struct {
	goos   string
	goarch string
	cgosup string
}

func main() {
	// CLI

	qs := len(os.Args) > 1 && os.Args[1] == "qs"
	if qs {
		cgoenabled := ""
		if len(os.Args) > 2 {
			cgoenabled = os.Args[2]
		}
		name := ""
		if len(os.Args) > 3 {
			name = os.Args[3]
		}
		Clibuild(name, cgoenabled)
		return
	}

	cc := len(os.Args) > 1 && os.Args[1] == "cc"
	if cc {
		c := exec.Command("go", "clean", "-cache")
		err := c.Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// TUI

	tview.Styles.BorderColor = tcell.ColorRebeccaPurple
	tview.Styles.PrimaryTextColor = tcell.ColorDeepSkyBlue
	tview.Styles.SecondaryTextColor = tcell.ColorDeepSkyBlue
	tview.Styles.TitleColor = tcell.ColorRebeccaPurple
	tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(20, 20, 20)

	app := tview.NewApplication()
	pages := tview.NewPages()

	// Show targets

	l, err := List()
	if err != nil {
		log.Fatal(err)
	}

	var tableData [][]string

	tableData = append(tableData, []string{"GOOS", "GOARCH", "CgoSupport", "FirstClass"})

	for x := range len(l) {
		tableData = append(tableData, []string{l[x].Os, l[x].Arch, fmt.Sprint(l[x].Cgo), fmt.Sprint(l[x].Fclass)})
	}

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	table.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRebeccaPurple))

	for row, cols := range tableData {
		for col, cell := range cols {
			c := tview.NewTableCell(cell).SetExpansion(1)
			if row == 0 {
				c.SetTextColor(tcell.ColorCornflowerBlue).SetSelectable(false)
			}
			table.SetCell(row, col, c)
		}
	}

	var s Selection

	table.SetFixed(1, 0)
	table.SetSelectedFunc(func(row, col int) {
		goos := table.GetCell(row, 0).Text
		goarch := table.GetCell(row, 1).Text
		cgosup := table.GetCell(row, 2).Text

		s.goos = goos
		s.goarch = goarch
		s.cgosup = cgosup

		pages.SwitchToPage("form")
	})

	table.SetBorder(true).SetTitle(" Targets ")

	// Show options

	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Custom Args", "", 20, nil, nil).
		AddCheckbox("Cgo", false, nil).
		AddCheckbox("Strip Symbols", false, nil). // -ldflags="-s -w"
		AddCheckbox("Trimpath", false, nil).      // -trimpath
		AddCheckbox("Verbose", false, nil).       // -v
		AddCheckbox("Race Detector", false, nil). // -race
		AddButton("Compile", func() {
			custargs := form.GetFormItemByLabel("Custom Args").(*tview.InputField).GetText()
			cgocheck := form.GetFormItemByLabel("Cgo").(*tview.Checkbox).IsChecked()
			stripsymbols := form.GetFormItemByLabel("Strip Symbols").(*tview.Checkbox).IsChecked()
			trimpath := form.GetFormItemByLabel("Trimpath").(*tview.Checkbox).IsChecked()
			verbose := form.GetFormItemByLabel("Verbose").(*tview.Checkbox).IsChecked()
			race := form.GetFormItemByLabel("Race Detector").(*tview.Checkbox).IsChecked()

			args := []string{} // Args to use

			if custargs != "" {
				args = append(args, custargs)
			}
			if stripsymbols {
				args = append(args, `-ldflags="-s -w"`)
			}
			if trimpath {
				args = append(args, "-trimpath")
			}
			if verbose {
				args = append(args, "-v")
			}
			if race {
				args = append(args, "-race")
			}

			passargs := strings.Join(args, " ")

			cgo := "0"

			if s.cgosup == "true" && cgocheck {
				cgo = "1"
			}

			msg, err := Compile(s.goos, s.goarch, cgo, passargs)
			if err != nil {
				app.Stop()
				fmt.Println(msg)
				log.Fatal(err)
			}
			app.Stop()
			fmt.Println(msg)
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("table")
		})

	form.SetBorder(true).SetTitle(" Options ")

	// Add pages

	pages.AddPage("table", table, true, true)
	pages.AddPage("form", form, true, false)

	e := app.SetRoot(pages, true).EnableMouse(true).Run()
	if e != nil {
		app.Stop()
		panic(e)
	}
}
