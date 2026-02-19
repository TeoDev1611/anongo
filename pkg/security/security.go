package security

import (
	"context"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pterm/pterm"
)

type SecurityMonitor struct {
	device string
	handle *pcap.Handle
}

func NewSecurityMonitor(device string) (*SecurityMonitor, error) {
	// Abrir dispositivo de red en modo promiscuo
	handle, err := pcap.OpenLive(device, 1024, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("error al abrir interfaz %s: %v", device, err)
	}
	return &SecurityMonitor{device: device, handle: handle}, nil
}

// StartARPMonitor vigila ataques de ARP Poisoning en una Goroutine
func (sm *SecurityMonitor) StartARPMonitor(ctx context.Context) {
	pterm.Info.Printf("Iniciando monitor Anti-MITM en la interfaz: %s\n", sm.device)

	packetSource := gopacket.NewPacketSource(sm.handle, sm.handle.LinkType())

	// Mapa para recordar IPs y sus MACs "legítimas" (aprendidas dinámicamente)
	arpCache := make(map[string]string)

	go func() {
		defer sm.handle.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case packet := <-packetSource.Packets():
				arpLayer := packet.Layer(layers.LayerTypeARP)
				if arpLayer == nil {
					continue
				}

				arp := arpLayer.(*layers.ARP)
				if arp.Operation != layers.ARPReply {
					continue
				}

				ipSource := net.IP(arp.SourceProtAddress).String()
				macSource := net.HardwareAddr(arp.SourceHwAddress).String()

				// Lógica de detección: Si la MAC cambia para la misma IP en poco tiempo, hay sospecha
				if oldMac, exists := arpCache[ipSource]; exists && oldMac != macSource {
					pterm.Fatal.Printf("\n[!] ALERTA DE SEGURIDAD: Possible ataque ARP Poisoning detectado para IP %s.\nOld MAC: %s | New MAC: %s\n", ipSource, oldMac, macSource)
					// TODO: Implementar bloqueo automático del atacante vía iptables
				} else {
					arpCache[ipSource] = macSource
				}
			}
		}
	}()
}

// GetDefaultInterface intenta adivinar la interfaz activa (wifi o ethernet)
func GetDefaultInterface() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// Ignorar loopback y dispositivos inactivos
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			return iface.Name
		}
	}
	return "wlan0" // Valor por defecto
}
