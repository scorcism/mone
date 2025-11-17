package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/scorcism/mone/cmd/services"
	"github.com/scorcism/mone/cmd/utils"
)

func BuildUI() {
	fmt.Printf("Builder is ONNNNNNNN")

	a := app.NewWithID(utils.AppID)
	w := a.NewWindow(utils.AppName)
	isDesktop := false
	selectedDeviceBinding := binding.NewString()

	ok := a.(desktop.App)
	if ok != nil {
		isDesktop = true
	}

	if isDesktop {
		fmt.Printf("Running on Desktop Environment\n")
		w.Resize(fyne.NewSize(800, 600))
	}

	selectedDeviceBinding.AddListener(binding.NewDataListener(func() {
		device, _ := selectedDeviceBinding.Get()
		fmt.Printf("Device changed to: %v\n", device)
		if device == "" {
			w.SetContent(screen1(w, selectedDeviceBinding))
		} else {
			w.SetContent(screen2(a, w, selectedDeviceBinding))
		}
	}))

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

	confirmICont := container.NewHBox(
		confirmBtn,
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

	return mainContent
}

func screen2(a fyne.App, win fyne.Window, selectedDeviceBinding binding.String) fyne.CanvasObject {

	startListenerBinding := binding.NewInt()
	startListenerBinding.Set(0)
	requestCountBinding := binding.NewInt()
	requestCountBinding.Set(0)
	captureTypeBinding := binding.NewString()

	headerContainer := screen2Header(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)
	contentContainer := screen2Content(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)

	content := container.NewBorder(headerContainer, nil, nil, nil, container.NewVScroll(container.NewStack(contentContainer)))

	return content
}

func screen2Content(selectedDeviceBinding binding.String, startListenerBinding binding.Int, captureType binding.String, requestCountBinding binding.Int) fyne.CanvasObject {
	requests := binding.NewStringList()
	device, _ := selectedDeviceBinding.Get()

	snapshotlen := int32(65535)
	promiscuous := true
	timeout := pcap.BlockForever
	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
	}

	localIps := services.GetLocalIps()
	fmt.Printf("LocalIps: %v", localIps)
	content := container.NewVBox()

	startListenerBinding.AddListener(binding.NewDataListener(func() {
		mode, _ := startListenerBinding.Get()
		switch mode {
		case 0:
			// Do nothing
		case 1:
			fmt.Printf("Started Listening on device: %s\n", device)
			go func() {
				packets := gopacket.NewPacketSource(handle, handle.LinkType())
				for packet := range packets.Packets() {
					timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata := services.LogPacketInfo(packet, localIps)
					rItem := NewRequestItem(timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata)
					content.Add(rItem.Container)
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

	updateContent := func() {
		requestsList, _ := requests.Get()
		content.Objects = nil
		for _, r := range requestsList {
			content.Add(widget.NewLabel(r))
		}
		content.Refresh()
	}

	requests.AddListener(binding.NewDataListener(func() {
		updateContent()
	}))

	updateContent()

	return content
}

func screen2Header(selectedDeviceBinding binding.String, startListenerBinding binding.Int, captureTypeBinding binding.String, requestCountBinding binding.Int) fyne.CanvasObject {
	// Incoming Button
	ib := widget.NewButton("Incoming", func() {
		captureTypeBinding.Set("INCOMING")
	})

	// Outgoing Button
	ob := widget.NewButton("Outgoing", func() {
		captureTypeBinding.Set("OUTGOING")
	})
	// Capture All Button
	cb := widget.NewButton("Capture All", func() {
		captureTypeBinding.Set("BOTH")
	})

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
			ib.Importance = widget.LowImportance
			ob.Importance = widget.HighImportance
			cb.Importance = widget.LowImportance
		case "BOTH":
			ib.Importance = widget.LowImportance
			ob.Importance = widget.LowImportance
			cb.Importance = widget.HighImportance
		default:
			ib.Importance = widget.LowImportance
			ob.Importance = widget.LowImportance
			cb.Importance = widget.LowImportance
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
