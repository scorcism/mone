package screens

import (
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/scorcism/mone/cmd/services"
	"github.com/scorcism/mone/cmd/utils"
)

func Screen1(win fyne.Window, selectedDeviceBinding binding.String) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Welcome to Mone!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})
	reloadInterfacesBinding := binding.NewBool()

	reloadBtn := widget.NewButtonWithIcon("Reload Interfaces", theme.ViewRefreshIcon(), func() {
		val, _ := reloadInterfacesBinding.Get()
		reloadInterfacesBinding.Set(!val)
	})
	moreInfoBtn := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		AppInfoWindow()
	})
	exitBtn := widget.NewButton("Exit", func() {
		win.Close()
	})

	leftContent := container.NewHBox(
		reloadBtn,
	)
	rightContent := container.NewHBox(
		exitBtn,
		moreInfoBtn,
	)

	btnContainer := container.NewHBox(leftContent, layout.NewSpacer(), rightContent)

	bodyContent := container.NewVBox(
		widget.NewLabel("Select your device type:"),
		handleInterfaceSelection(selectedDeviceBinding, reloadInterfacesBinding),
		btnContainer,
	)

	mainContent := container.NewVBox(
		title,
		bodyContent,
	)

	return container.New(layout.NewCenterLayout(), mainContent)
}

func handleInterfaceSelection(selectedDeviceBinding binding.String, reloadInterfacesBinding binding.Bool) *widget.RadioGroup {
	radio := widget.NewRadioGroup([]string{}, func(value string) {
		device := strings.Split(value, " :: ")[0][8:]
		selectedDeviceBinding.Set(device)
	})

	reloadInterfacesBinding.AddListener(binding.NewDataListener(func() {
		interfsDesc := []string{}
		interfs := services.GetLocalInterfaces()
		for _, interf := range interfs {
			f := fmt.Sprintf("Device: %s :: Description: %s", interf.Name, interf.Description)
			interfsDesc = append(interfsDesc, f)
		}
		radio.Options = interfsDesc
		radio.Refresh()
	}))

	return radio
}

func AppInfoWindow() {
	a := fyne.CurrentApp()
	win := a.NewWindow("About Mone")
	title := widget.NewLabelWithStyle("Mone - Network Packet Analyzer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: false})
	description := widget.NewLabel(utils.Description)
	description.Wrapping = fyne.TextWrapWord
	version := widget.NewLabel(fmt.Sprintf("Version: %s", utils.Version))
	developer := widget.NewLabel(fmt.Sprintf("Developer: %s", utils.Developer))
	u, _ := url.Parse(utils.GitHub)
	github := widget.NewHyperlink(fmt.Sprintf("GitHub: %s", utils.GitHub), u)
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		description,
		version,
		developer,
		github,
	)

	win.SetContent(content)
	win.Resize(fyne.NewSize(400, 300))
	win.Show()
}
