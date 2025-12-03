package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/scorcism/mone/cmd/services"
)

func Screen2(selectedDeviceBinding binding.String) fyne.CanvasObject {

	startListenerBinding := binding.NewInt()
	startListenerBinding.Set(0)
	requestCountBinding := binding.NewInt()
	requestCountBinding.Set(0)
	captureTypeBinding := binding.NewString()

	headerContainer := Header(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)
	contentContainer := Body(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)

	content := container.NewBorder(headerContainer, nil, nil, nil, container.NewVScroll(contentContainer))

	return content
}

func Body(selectedDeviceBinding binding.String, startListenerBinding binding.Int, captureType binding.String, requestCountBinding binding.Int) fyne.CanvasObject {
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
			btn := widget.NewButtonWithIcon("", theme.InfoIcon(), nil)
			label := widget.NewLabel("")
			return container.NewHBox(btn, label)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			raw, _ := item.(binding.Untyped).Get()
			rItem := raw.(*RequestItem)

			row := obj.(*fyne.Container)
			btn := row.Objects[0].(*widget.Button)
			label := row.Objects[1].(*widget.Label)

			btn.OnTapped = func() {
				rItem.ShowMetadataWindow()
			}
			label.SetText(fmt.Sprintf("[%s] [%s] %s %s:%s -> %s:%s Size: %d bytes",
				rItem.Timestamp, rItem.Direction, rItem.Proto, rItem.Src, rItem.SrcPort,
				rItem.Dst, rItem.DstPort, rItem.Size))
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
					timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata := services.PacketInfo(packet, localIps)
					rItem := NewRequestItem(timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata)
					data.Append(rItem)
					time.Sleep(1 * time.Millisecond)
					currentRequestsCount, _ := requestCountBinding.Get()
					requestCountBinding.Set(currentRequestsCount + 1)
				}
			}()
		default:
			// Stop Listening
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

func Header(selectedDeviceBinding binding.String,
	startListenerBinding binding.Int,
	captureTypeBinding binding.String,
	requestCountBinding binding.Int) fyne.CanvasObject {

	c := widget.NewLabel("")

	// Incoming Button
	incomingBtn := widget.NewButton("Incoming", func() {
		captureTypeBinding.Set("INCOMING")
	})
	// Outgoing Button
	outgoingBtn := widget.NewButton("Outgoing", func() {
		captureTypeBinding.Set("OUTGOING")
	})
	// Capture All Button
	captureAllBtn := widget.NewButton("Capture All", func() {
		captureTypeBinding.Set("BOTH")
	})

	incomingBtn.Importance = widget.LowImportance
	outgoingBtn.Importance = widget.LowImportance
	captureAllBtn.Importance = widget.LowImportance

	// Start Btn
	startBtn := widget.NewButton("Start", func() {
		startListenerBinding.Set(1)
	})
	// Stop
	stopBtn := widget.NewButton("Stop", func() {
		startListenerBinding.Set(2)
	})
	// Back
	backBtn := widget.NewButton("Back", func() {
		selectedDeviceBinding.Set("")
	})

	captureTypeBinding.AddListener(binding.NewDataListener(func() {
		captureTypeValue, _ := captureTypeBinding.Get()
		if captureTypeValue == "" {
			startBtn.Disable()
			stopBtn.Disable()
			return
		} else {
			startBtn.Enable()
		}
		switch captureTypeValue {
		case "INCOMING":
			incomingBtn.Importance = widget.HighImportance
			outgoingBtn.Importance = widget.LowImportance
			captureAllBtn.Importance = widget.LowImportance
		case "OUTGOING":
			outgoingBtn.Importance = widget.HighImportance
			incomingBtn.Importance = widget.LowImportance
			captureAllBtn.Importance = widget.LowImportance
		case "BOTH":
			captureAllBtn.Importance = widget.HighImportance
			incomingBtn.Importance = widget.LowImportance
			outgoingBtn.Importance = widget.LowImportance
		default:

		}
	}))

	startListenerBinding.AddListener(binding.NewDataListener(func() {
		mode, _ := startListenerBinding.Get()
		switch mode {
		case 0:
			startBtn.Disable()
			stopBtn.Disable()
		case 1:
			startBtn.Disable()
			stopBtn.Enable()
		default:
			startBtn.Enable()
			stopBtn.Disable()
		}
	}))

	requestCountBinding.AddListener(binding.NewDataListener(func() {
		count, _ := requestCountBinding.Get()
		c.SetText(fmt.Sprintf("Requests: %d", count))
	}))

	lbg := container.NewHBox(
		incomingBtn,
		outgoingBtn,
		captureAllBtn,
	)

	rbg := container.NewHBox(
		c,
		startBtn,
		stopBtn,
		backBtn,
	)

	// Header
	header := container.NewBorder(nil, nil, nil, nil, container.NewHBox(lbg, layout.NewSpacer(), rbg))
	return header
}
