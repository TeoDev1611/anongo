package network

import (
	"context"
	"os/user"
	"time"

	"github.com/TeoDev1611/anongo/pkg/tor"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pterm/pterm"
)

const (
	natChain    = "ANONGO_NAT"
	filterChain = "ANONGO_FILTER"
)

type NetworkManager struct {
	ipt            *iptables.IPTables
	ipt6           *iptables.IPTables
	tm             *tor.TorManager
	maintenance    bool
	watchdogActive bool
}

func NewNetworkManager(tm *tor.TorManager) (*NetworkManager, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, err
	}
	// Soporte para IPv6 (para bloquearlo)
	ipt6, _ := iptables.NewWithProtocol(iptables.ProtocolIPv6)

	return &NetworkManager{ipt: ipt, ipt6: ipt6, tm: tm, maintenance: false, watchdogActive: false}, nil
}

func (nm *NetworkManager) SetMaintenance(status bool) {
	nm.maintenance = status
}

func (nm *NetworkManager) StartAnonymity(ctx context.Context) error {
	pterm.Info.Println("Iniciando blindaje de red (Anti-Leak Mode)...")

	// 0. Crear cadenas personalizadas si no existen
	nm.setupChains()

	// 1. Bloquear IPv6 totalmente para evitar fugas
	if nm.ipt6 != nil {
		nm.ipt6.AppendUnique("filter", "OUTPUT", "-j", "DROP")
		nm.ipt6.AppendUnique("filter", "INPUT", "-j", "DROP")
	}

	// 2. Exclusión de Usuario Tor
	torUser, err := user.Lookup("tor")
	if err != nil {
		torUser, _ = user.Lookup("debian-tor")
	}
	if torUser != nil {
		nm.ipt.AppendUnique("nat", natChain, "-m", "owner", "--uid-owner", torUser.Uid, "-j", "RETURN")
	}

	// 3. Excluir Redes Locales (Para no perder acceso al Router/LAN)
	localNets := []string{"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	for _, net := range localNets {
		nm.ipt.AppendUnique("nat", natChain, "-d", net, "-j", "RETURN")
	}

	// 4. Redirigir DNS (UDP 53)
	nm.ipt.AppendUnique("nat", natChain, "-p", "udp", "--dport", "53", "-j", "REDIRECT", "--to-ports", "5353")

	// 5. Bloquear el resto de UDP
	nm.ipt.AppendUnique("filter", filterChain, "-p", "udp", "!", "--dport", "53", "-j", "DROP")

	// 6. Redirigir TODO el tráfico TCP a Tor
	nm.ipt.AppendUnique("nat", natChain, "-p", "tcp", "-j", "REDIRECT", "--to-ports", "9040")

	// 7. Insertar saltos quirúrgicos
	nm.ipt.Insert("nat", "OUTPUT", 1, "-j", natChain)
	nm.ipt.Insert("filter", "OUTPUT", 1, "-j", filterChain)

	pterm.Success.Println("Iptables aplicadas quirúrgicamente.")
	
	// Solo lanzar un vigilante si no hay uno corriendo
	if !nm.watchdogActive {
		nm.watchdogActive = true
		go nm.watchdog(ctx)
	}

	return nil
}

func (nm *NetworkManager) setupChains() {
	_ = nm.ipt.NewChain("nat", natChain)
	_ = nm.ipt.NewChain("filter", filterChain)
}

func (nm *NetworkManager) StopAnonymity() {
	nm.maintenance = true // Bloquear cualquier acción del vigilante de inmediato
	pterm.Warning.Println("Desactivando blindaje y restaurando red...")

	// 1. Eliminar saltos
	_ = nm.ipt.Delete("nat", "OUTPUT", "-j", natChain)
	_ = nm.ipt.Delete("filter", "OUTPUT", "-j", filterChain)

	// 2. Limpiar y eliminar cadenas
	_ = nm.ipt.ClearChain("nat", natChain)
	_ = nm.ipt.DeleteChain("nat", natChain)
	_ = nm.ipt.ClearChain("filter", filterChain)
	_ = nm.ipt.DeleteChain("filter", filterChain)

	if nm.ipt6 != nil {
		_ = nm.ipt6.Delete("filter", "OUTPUT", "-j", "DROP")
		_ = nm.ipt6.Delete("filter", "INPUT", "-j", "DROP")
	}

	pterm.Success.Println("Red restaurada al estado original.")
}

func (nm *NetworkManager) watchdog(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
		nm.watchdogActive = false
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Si estamos en mantenimiento o cerrando, ignorar
			if nm.maintenance {
				continue
			}

			// Verificar reglas básicas
			exists, _ := nm.ipt.Exists("nat", "OUTPUT", "-j", natChain)
			if !exists && !nm.maintenance {
				pterm.Error.Println("\n[!] ALERTA: Reglas eliminadas. Restaurando...")
				nm.StartAnonymity(ctx)
				continue
			}

			// Verificar conexión
			ok, _ := nm.tm.CheckTorConnection()
			if !ok && !nm.maintenance {
				pterm.Error.Println("\n[!] FUGA DETECTADA. Re-aplicando blindaje...")
				nm.StartAnonymity(ctx)
			}
		}
	}
}
