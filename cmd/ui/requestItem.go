package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
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
	item.Container = item.Build()
	return item
}

func (ri *RequestItem) Build() fyne.CanvasObject {
	content := container.NewHBox()
	btn := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ri.showMetadataWindow()
	})
	content.Add(btn)
	content.Add(widget.NewLabel(fmt.Sprintf("[%s]", ri.Timestamp)))
	content.Add(widget.NewLabel(ri.Proto))
	content.Add(widget.NewLabel(ri.Direction))
	content.Add(widget.NewLabel(ri.Src + ":" + ri.SrcPort))
	content.Add(widget.NewLabel("->"))
	content.Add(widget.NewLabel(ri.Dst + ":" + ri.DstPort))
	content.Add(widget.NewLabel(fmt.Sprintf("Size: %d bytes", ri.Size)))
	return content
}

func (ri *RequestItem) showMetadataWindow() {
	a := fyne.CurrentApp()
	win := a.NewWindow("Packet Metadata")
	win.SetContent(widget.NewLabel(fmt.Sprintf("Full Metadata: %+v", ri.Metadata)))
	win.Resize(fyne.NewSize(100, 300))
	win.Show()
}
