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
	r := widget.NewLabel(fmt.Sprintf("[%s] %s %s:%s -> %s:%s Size: %d bytes", ri.Timestamp, ri.Proto, ri.Src, ri.SrcPort, ri.Dst, ri.DstPort, ri.Size))
	fmt.Printf("Request: %v\n", r.Text)
	content.Add(r)
	return content
}

func (ri *RequestItem) showMetadataWindow() {
	a := fyne.CurrentApp()
	win := a.NewWindow("Packet Metadata")
	win.SetContent(widget.NewLabel(fmt.Sprintf("Full Metadata: %+v", ri.Metadata)))
	win.Resize(fyne.NewSize(100, 300))
	win.Show()
}
