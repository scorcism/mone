package screens

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

func (ri *RequestItem) ShowMetadataWindow() {
	app := fyne.CurrentApp()
	win := app.NewWindow("Request Metadata")

	sourceWhoIsBinding := binding.NewString()
	destinationWhoIsBinding := binding.NewString()

	sourceWhoIsBinding.Set("Fetching WHOIS information...")
	destinationWhoIsBinding.Set("Fetching WHOIS information...")

	sp, _ := strconv.ParseInt(ri.SrcPort, 10, 32)
	sourceAppInfo := utils.GetServiceByPort(uint32(sp))
	dp, _ := strconv.ParseInt(ri.DstPort, 10, 32)
	desAppInfo := utils.GetServiceByPort(uint32(dp))

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
			"Application Using This Port\n"+
				"Source: %v\n"+
				"Destination: %v\n",
			sourceAppInfo,
			desAppInfo,
		),
	)

	fullMeta := widget.NewAccordion(
		widget.NewAccordionItem("Full Metadata (Raw)", widget.NewLabel(fmt.Sprintf("%+v", meta))),
	)

	sourceWhoIsLbl := widget.NewLabelWithStyle(fmt.Sprintf("Source [%s] WHOIS Lookup", ri.Src), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})
	destinationWhoIsLbl := widget.NewLabelWithStyle(fmt.Sprintf("Destination [%s] WHOIS Lookup", ri.Dst), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	go func() {

		info, err := utils.PerformWhoisLookup(ri.Src)
		infoDst, errDst := utils.PerformWhoisLookup(ri.Dst)

		if err != nil {
			sourceWhoIsBinding.Set(fmt.Sprintf("Error fetching WHOIS: %v", err))
		} else {
			sourceWhoIsBinding.Set(fmt.Sprintf("Raw WHOIS Output:\n\n%s", info))
		}

		if errDst != nil {
			destinationWhoIsBinding.Set(fmt.Sprintf("Error fetching WHOIS: %v", errDst))
		} else {
			destinationWhoIsBinding.Set(fmt.Sprintf("Raw WHOIS Output:\n\n%s", infoDst))
		}
	}()

	whoIsContentSrc := container.NewVBox(
		sourceWhoIsLbl,
		widget.NewLabelWithData(sourceWhoIsBinding),
		destinationWhoIsLbl,
		widget.NewLabelWithData(destinationWhoIsBinding),
	)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Request Metadata", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		basicInfo,
		portAppInfo,
		captureInfo,
		fullMeta,
		whoIsContentSrc,
	)

	scroll := container.NewVScroll(content)

	win.SetContent(scroll)
	win.Resize(fyne.NewSize(400, 450))
	win.SetFixedSize(true)
	win.Show()
}
