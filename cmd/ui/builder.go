package ui

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/scorcism/mone/cmd/services"
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
				showInfoWindow()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		ok.SetSystemTrayMenu(m)
		w.Resize(fyne.NewSize(800, 600))
	}

	selectedDeviceBinding.AddListener(binding.NewDataListener(func() {
		device, _ := selectedDeviceBinding.Get()
		if device == "" {
			w.SetContent(screen1(w, selectedDeviceBinding))
		} else {
			w.SetContent(screen2(selectedDeviceBinding))
		}
	}))

	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.ShowAndRun()
}

func screen1(win fyne.Window, selectedDeviceBinding binding.String) fyne.CanvasObject {
	title := widget.NewLabel("Welcome to Mone!")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true
	title.TextStyle.Monospace = true
	device := ""

	interfs := services.GetLocalInterfaces()
	interfsDesc := []string{}
	for _, interf := range interfs {
		f := fmt.Sprintf("Device: %s :: Description: %s", interf.Name, interf.Description)
		interfsDesc = append(interfsDesc, f)
	}
	radio := widget.NewRadioGroup(interfsDesc, func(value string) {
		device = strings.Split(value, " :: ")[0][8:]
	})

	confirmBtn := widget.NewButton("Confirm", func() {
		if device == "" {
			dialog.ShowError(fmt.Errorf("no device selected"), win)
			return
		}
		selectedDeviceBinding.Set(device)
	})

	// More info btn
	moreInfoBtn := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		showInfoWindow()
	})

	exitBtn := widget.NewButton("Exit", func() {
		win.Close()
	})

	confirmICont := container.NewHBox(
		confirmBtn,
		exitBtn,
		moreInfoBtn,
	)

	bodyContent := container.NewVBox(
		widget.NewLabel("Select your device type:"),
		radio,
		confirmICont,
	)

	mainContent := container.NewVBox(
		title,
		bodyContent,
	)

	return container.New(layout.NewCenterLayout(), mainContent)
}

func screen2(selectedDeviceBinding binding.String) fyne.CanvasObject {

	startListenerBinding := binding.NewInt()
	startListenerBinding.Set(0)
	requestCountBinding := binding.NewInt()
	requestCountBinding.Set(0)
	captureTypeBinding := binding.NewString()

	headerContainer := screen2Header(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)
	contentContainer := screen2Content(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)

	content := container.NewBorder(headerContainer, nil, nil, nil, container.NewVScroll(contentContainer))

	return content
}

func screen2Content(selectedDeviceBinding binding.String, startListenerBinding binding.Int, captureType binding.String, requestCountBinding binding.Int) fyne.CanvasObject {
	device, _ := selectedDeviceBinding.Get()

	requestCaptureRulesBinding := binding.NewString()
	requestCaptureRulesBinding.Set("")

	hostIps := services.GetLocalIps()
	snapshotlen := int32(65535)
	promiscuous := true
	timeout := pcap.BlockForever
	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
	}

	localIps := services.GetLocalIps()
	data := binding.NewUntypedList()

	dataList := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return container.New(nil)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			raw, _ := item.(binding.Untyped).Get()

			cnt := raw.(*fyne.Container)
			row := obj.(*fyne.Container)
			row.Objects = nil
			row.Add(cnt)
		},
	)

	startListenerBinding.AddListener(binding.NewDataListener(func() {
		mode, _ := startListenerBinding.Get()
		switch mode {
		case 0:
			// Nothing is selected
		case 1:
			go func() {
				packets := gopacket.NewPacketSource(handle, handle.LinkType())
				for packet := range packets.Packets() {
					timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata := services.LogPacketInfo(packet, localIps)
					rItem := NewRequestItem(timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata)
					data.Append(rItem.Container)
					time.Sleep(1 * time.Millisecond)
					currentRequestsCount, _ := requestCountBinding.Get()
					requestCountBinding.Set(currentRequestsCount + 1)
				}
			}()
		default:
			// Stop Listening
			fmt.Printf("Stopped Listening on device: %s\n", device)
			handle.Close()
		}
	}))

	captureType.AddListener(binding.NewDataListener(func() {
		value, _ := captureType.Get()
		filters := ""

		switch value {
		case "INCOMING":
			for i, ip := range hostIps {
				if i > 0 {
					filters += " or "
				}
				filters += "dst host " + ip.String()
			}
			handle.SetBPFFilter(filters)
		case "OUTGOING":
			for i, ip := range hostIps {
				if i > 0 {
					filters += " or "
				}
				filters += "src host " + ip.String()
			}
			handle.SetBPFFilter(filters)
		case "BOTH":
			filters = ""
			handle.SetBPFFilter(filters)
		default:
			filters = ""
			handle.SetBPFFilter(filters)
		}
	}))

	return dataList
}

func screen2Header(selectedDeviceBinding binding.String,
	startListenerBinding binding.Int,
	captureTypeBinding binding.String,
	requestCountBinding binding.Int) fyne.CanvasObject {

	// Incoming Button
	ib := widget.NewButton("Incoming", func() {
		captureTypeBinding.Set("INCOMING")
	})
	ib.Importance = widget.LowImportance

	// Outgoing Button
	ob := widget.NewButton("Outgoing", func() {
		captureTypeBinding.Set("OUTGOING")
	})
	ob.Importance = widget.LowImportance

	// Capture All Button
	cb := widget.NewButton("Capture All", func() {
		captureTypeBinding.Set("BOTH")
	})
	cb.Importance = widget.LowImportance

	c := widget.NewLabel("")

	// Start Btn
	sb := widget.NewButton("Start", func() {
		startListenerBinding.Set(1)
	})

	// Stop
	stb := widget.NewButton("Stop", func() {
		startListenerBinding.Set(2)
	})

	// Back
	bb := widget.NewButton("Back", func() {
		selectedDeviceBinding.Set("")
	})

	captureTypeBinding.AddListener(binding.NewDataListener(func() {
		captureTypeValue, _ := captureTypeBinding.Get()
		if captureTypeValue == "" {
			sb.Disable()
			stb.Disable()
			return
		} else {
			sb.Enable()
		}
		switch captureTypeValue {
		case "INCOMING":
			ib.Importance = widget.HighImportance
			ob.Importance = widget.LowImportance
			cb.Importance = widget.LowImportance
		case "OUTGOING":
			ob.Importance = widget.HighImportance
			ib.Importance = widget.LowImportance
			cb.Importance = widget.LowImportance
		case "BOTH":
			cb.Importance = widget.HighImportance
			ib.Importance = widget.LowImportance
			ob.Importance = widget.LowImportance
		default:

		}
	}))

	startListenerBinding.AddListener(binding.NewDataListener(func() {
		mode, _ := startListenerBinding.Get()
		switch mode {
		case 0:
			sb.Disable()
			stb.Disable()
		case 1:
			sb.Disable()
			stb.Enable()
		default:
			sb.Enable()
			stb.Disable()
		}
	}))

	requestCountBinding.AddListener(binding.NewDataListener(func() {
		count, _ := requestCountBinding.Get()
		c.SetText(fmt.Sprintf("Requests: %d", count))
	}))

	lbg := container.NewHBox(
		ib,
		ob,
		cb,
	)

	rbg := container.NewHBox(
		c,
		sb,
		stb,
		bb,
	)

	// Header
	header := container.NewBorder(nil, nil, nil, nil, container.NewHBox(lbg, layout.NewSpacer(), rbg))
	return header
}

func showInfoWindow() {
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
