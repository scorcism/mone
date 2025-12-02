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
	contentContainer := Content(selectedDeviceBinding, startListenerBinding, captureTypeBinding, requestCountBinding)

	content := container.NewBorder(headerContainer, nil, nil, nil, container.NewVScroll(contentContainer))

	return content
}

func Content(selectedDeviceBinding binding.String, startListenerBinding binding.Int, captureType binding.String, requestCountBinding binding.Int) fyne.CanvasObject {
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
			// If size is < 1 then ignore update
			if rItem.Size < 1 {
				return
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
					timestamp, proto, direction, src, srcPort, dst, dstPort, size, metadata := services.LogPacketInfo(packet, localIps)
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
