package services

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/scorcism/mone/cmd/types"
)

func Run() {
	fmt.Printf("Service is running\n")
	interfs := getLocalInterfaces()
	for _, interf := range interfs {
		fmt.Printf("Interface: %s :: %s\n", interf.Name, interf.Description)
	}

	device := "\\Device\\NPF_{4C95BE1E-B86A-4CB9-AB63-095864C9E90B}"
	snapshotlen := int32(65535)
	promiscuous := true
	timeout := pcap.BlockForever
	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
		return
	}
	defer handle.Close()

	localIps := getLocalIps()
	fmt.Printf("LocalIps: %v", localIps)

	packets := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packets.Packets() {
		logPacketInfo(packet, localIps)
		break
	}
}

func logPacketInfo(packet gopacket.Packet, localIps []net.IP) {
	netLayer := packet.NetworkLayer()
	transLayer := packet.TransportLayer()

	if netLayer == nil || transLayer == nil {
		fmt.Println("No network or transport layer found in packet")
		return
	}
	src := netLayer.NetworkFlow().Src().String()
	dst := netLayer.NetworkFlow().Dst().String()
	fmt.Printf("Network Layer: %s -> %s\n", src, dst)
	// srcIP, dstIP := netLayer.NetworkFlow().Endpoints()
	srcPort, dstPort := transLayer.TransportFlow().Endpoints()

	direction := getDirection(src, dst, localIps)
	proto := transLayer.LayerType().String()
	size := len(packet.Data())

	timestamp := packet.Metadata().Timestamp.Format(time.RFC3339)

	fmt.Printf("Packet Info: [%s] %s | %s | %s:%s -> %s:%s | Size: %d bytes\n", timestamp, proto, direction, src, srcPort.String(), dst, dstPort.String(), size)
}

func getLocalIps() []net.IP {
	localIPs := []net.IP{}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIPs = append(localIPs, ipNet.IP)
			}
		}
	}
	return localIPs
}

func getDirection(srcIP, dstIP string, localIps []net.IP) string {
	fmt.Printf("Checking direction for SrcIP: %s, DstIP: %s\n", srcIP, dstIP)
	if isLocalIP(srcIP, localIps) && !isLocalIP(dstIP, localIps) {
		return "OUTGOING"
	}
	if !isLocalIP(srcIP, localIps) && isLocalIP(dstIP, localIps) {
		return "INCOMING"
	}
	return "UNKNOWN"
}

func isLocalIP(ipStr string, localIps []net.IP) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, localIP := range localIps {
		if ip.Equal(localIP) {
			return true
		}
	}
	return false
}

func getLocalInterfaces() []types.LocalInterface {
	interfs := []types.LocalInterface{}
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Printf("Error finding devices: %v\n", err)
		return nil
	}

	for _, device := range devices {
		// fmt.Printf("%s :: %s\n", device.Name, device.Description)
		interfs = append(interfs, types.LocalInterface{
			Name:        device.Name,
			Description: device.Description,
		})
	}
	return interfs
}
