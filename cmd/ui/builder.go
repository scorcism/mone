package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/scorcism/mone/cmd/services"
	"github.com/scorcism/mone/cmd/utils"
)

func BuildUI() {
	fmt.Printf("Builder is ONNNNNNNN")

	a := app.NewWithID(utils.AppID)
	w := a.NewWindow(utils.AppName)
	isDesktop := false

	ok := a.(desktop.App)
	if ok != nil {
		isDesktop = true
	}

	if isDesktop {
		fmt.Printf("Running on Desktop Environment\n")
		w.Resize(fyne.NewSize(800, 600))
	}

	w.SetContent(screen1())
	w.ShowAndRun()
}

func screen1() fyne.CanvasObject {
	title := widget.NewLabel("Welcome to Mone!")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true
	title.TextStyle.Monospace = true

	interfs := services.GetLocalInterfaces()
	interfsDesc := []string{}
	for _, interf := range interfs {
		f := fmt.Sprintf("Device: %s :: Description: %s", interf.Name, interf.Description)
		interfsDesc = append(interfsDesc, f)
	}
	// devices := []string{"Desktop", "Mobile", "Tablet"}
	radio := widget.NewRadioGroup(interfsDesc, func(value string) {
		fmt.Println("Radio set to", value)
		device := strings.Split(value, " :: ")[0][8:]
		fmt.Printf("Selected Device: %s\n", device)
	})

	bodyContent := container.NewVBox(
		widget.NewLabel("Select your device type:"),
		radio,
	)

	mainContent := container.NewVBox(
		title,
		bodyContent,
	)

	return mainContent
}

func screen2() fyne.CanvasObject {
	return widget.NewLabel("Screen 2")
}
