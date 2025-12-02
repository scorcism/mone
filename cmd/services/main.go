package services

import (
	"net"
	"slices"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/scorcism/mone/cmd/types"
)

func LogPacketInfo(packet gopacket.Packet, localIps []net.IP) (string, string, string, string, string, string, string, int, gopacket.PacketMetadata) {
	netLayer := packet.NetworkLayer()
	transLayer := packet.TransportLayer()

	if netLayer == nil || transLayer == nil {
		return "", "", "", "", "", "", "", 0, gopacket.PacketMetadata{}
	}
	src := netLayer.NetworkFlow().Src().String()
	dst := netLayer.NetworkFlow().Dst().String()
	srcPort, dstPort := transLayer.TransportFlow().Endpoints()

	direction := getDirection(src, dst, localIps)
	proto := transLayer.LayerType().String()
	size := len(packet.Data())

	timestamp := packet.Metadata().Timestamp.Local().Format(time.RFC3339)
	metadata := packet.Metadata()

	return timestamp, proto, direction, src, srcPort.String(), dst, dstPort.String(), size, *metadata
}

func GetLocalIps() []net.IP {
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
	return slices.ContainsFunc(localIps, ip.Equal)
}

func GetLocalInterfaces() []types.LocalInterface {
	interfs := []types.LocalInterface{}
	devices, err := pcap.FindAllDevs()
	if err != nil {
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
