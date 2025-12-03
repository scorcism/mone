package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/scorcism/mone/cmd/ui/screens"
	"github.com/scorcism/mone/cmd/utils"
)

func BuildUI() {
	a := app.NewWithID(utils.AppID)
	w := a.NewWindow(utils.AppName)
	isDesktop := false
	selectedDeviceBinding := binding.NewString()

	ok := a.(desktop.App)
	if ok != nil {
		isDesktop = true
	}

	if isDesktop {
		m := fyne.NewMenu("Mone",
			fyne.NewMenuItem("Open Mone", func() {
				w.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About Mone", func() {
				screens.AppInfoWindow()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		ok.SetSystemTrayMenu(m)
		w.Resize(fyne.NewSize(840, 600))
	}

	selectedDeviceBinding.AddListener(binding.NewDataListener(func() {
		device, _ := selectedDeviceBinding.Get()
		if device == "" {
			w.SetContent(screens.Screen1(w, selectedDeviceBinding))
		} else {
			w.SetContent(screens.Screen2(selectedDeviceBinding))
		}
	}))

	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.ShowAndRun()
}
