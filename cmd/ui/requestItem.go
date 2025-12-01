package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	"github.com/scorcism/mone/cmd/utils"
)

type RequestItem struct {
	Timestamp string
	Proto     string
	Direction string
	Src       string
	SrcPort   string
	Dst       string
	DstPort   string
	Size      int
	Metadata  gopacket.PacketMetadata

	Container fyne.CanvasObject
}

func NewRequestItem(timestamp, proto, direction, src, srcPort, dst, dstPort string, size int, metadata gopacket.PacketMetadata) *RequestItem {
	item := &RequestItem{
		Timestamp: timestamp,
		Proto:     proto,
		Direction: direction,
		Src:       src,
		SrcPort:   srcPort,
		Dst:       dst,
		DstPort:   dstPort,
		Size:      size,
		Metadata:  metadata,
	}
	return item
}
func (ri *RequestItem) showMetadataWindow() {
	app := fyne.CurrentApp()
	win := app.NewWindow("Request Metadata")

	sp, _ := strconv.ParseInt(ri.SrcPort, 10, 32)
	sourceAppInfo := utils.GetServiceByPort(uint32(i))
	dp, _ := strconv.ParseInt(ri.DstPort, 10, 32)
	desAppInfo := utils.GetServiceByPort(uint32(j))

	meta := ri.Metadata
	capture := meta.CaptureInfo

	basicInfo := widget.NewRichTextWithText(
		fmt.Sprintf(
			"Timestamp: %s\n"+
				"Protocol: %s\n"+
				"Direction: %s\n"+
				"Src: %s:%s\n"+
				"Dst: %s:%s\n"+
				"Packet Size: %d bytes\n",
			ri.Timestamp,
			ri.Proto,
			ri.Direction,
			ri.Src, ri.SrcPort,
			ri.Dst, ri.DstPort,
			ri.Size,
		),
	)

	captureInfo := widget.NewRichTextWithText(
		fmt.Sprintf(
			"Capture Info\n"+
				"Length: %d\n"+
				"Capture Length: %d\n",
			capture.Length,
			capture.CaptureLength,
		),
	)

	portAppInfo := widget.NewRichTextWithText(
		fmt.Sprintf(
			"Port app Info\n"+
				"Source: %v\n"+
				"Destination: %v\n",
			sourceAppInfo,
			desAppInfo,
		),
	)

	fullMeta := widget.NewAccordion(
		widget.NewAccordionItem("Full Metadata (Raw)", widget.NewLabel(fmt.Sprintf("%+v", meta))),
	)

	srcBtn := widget.NewButtonWithIcon("Source WHOIS Lookup", theme.InfoIcon(), func() {
		ri.ShowWhoIsWindow(ri.Src)
	})

	dstBtn := widget.NewButtonWithIcon("Destination WHOIS Lookup", theme.InfoIcon(), func() {
		ri.ShowWhoIsWindow(ri.Dst)
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle("Request Metadata", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		basicInfo,
		portAppInfo,
		captureInfo,
		fullMeta,
		srcBtn,
		dstBtn,
	)

	scroll := container.NewVScroll(content)

	win.SetContent(scroll)
	win.Resize(fyne.NewSize(400, 450))
	win.Show()
}

func (ri *RequestItem) ShowWhoIsWindow(ip string) {
	app := fyne.CurrentApp()
	win := app.NewWindow(fmt.Sprintf("WHOIS: %s", ip))
	infoBinding := binding.NewString()

	title := widget.NewLabelWithStyle(
		fmt.Sprintf("WHOIS Information for: %s", ip),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	info, err := utils.PerformWhoisLookup(ip)

	if err != nil {
		infoBinding.Set(fmt.Sprintf("Error fetching WHOIS: %v", err))
	} else {
		infoBinding.Set(fmt.Sprintf("Raw WHOIS Output:\n\n%s", info))
	}

	l := widget.NewLabelWithData(infoBinding)
	scroll := container.NewScroll(container.NewVBox(title, l))

	win.SetContent(scroll)
	win.Resize(fyne.NewSize(600, 400))
	win.Show()
}
