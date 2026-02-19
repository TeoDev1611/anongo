package tor

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/net/proxy"
)

type TorManager struct {
	cmd *exec.Cmd
}

func NewTorManager() *TorManager {
	return &TorManager{}
}

func (tm *TorManager) Start() error {
	pterm.Info.Printf("Iniciando Tor: %s\n", time.Now().Format(time.Kitchen))

	// 0. Identificar usuario Tor para evitar el bucle de tráfico
	torUser, err := user.Lookup("tor")
	if err != nil {
		torUser, _ = user.Lookup("debian-tor")
	}

	// 1. LIMPIEZA PREVENTIVA
	exec.Command("pkill", "-9", "tor").Run() 
	time.Sleep(1 * time.Second)

	dataDir := "/tmp/tor_anongo"
	os.RemoveAll(dataDir)
	os.MkdirAll(dataDir, 0700)

	// Dar permisos al usuario Tor si existe
	if torUser != nil {
		exec.Command("chown", "-R", torUser.Username+":"+torUser.Username, dataDir).Run()
	}

	emptyTorrc := dataDir + "/torrc"
	os.WriteFile(emptyTorrc, []byte(""), 0600)

	// Argumentos de optimización para velocidad y seguridad
	args := []string{
		"-f", emptyTorrc, 
		"--TransPort", "9040",
		"--DNSPort", "5353",
		"--SocksPort", "9050",
		"--DataDirectory", dataDir,
		"--Log", "notice stdout",
		"--AvoidDiskWrites", "1",           // Reduce I/O para mayor velocidad
		"--FastFirstHopPK", "1",            // Acelera la creación del primer salto
		"--MaxCircuitDirtiness", "30",      // Cambia de circuito cada 30s si es necesario
		"--CircuitBuildTimeout", "15",      // No esperar demasiado a nodos lentos
	}

	if torUser != nil {
		args = append(args, "--User", torUser.Username)
	}

	tm.cmd = exec.Command("tor", args...)

	stdout, err := tm.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creando pipe: %v", err)
	}

	if err := tm.cmd.Start(); err != nil {
		return fmt.Errorf("error al ejecutar binario tor: %v", err)
	}

	pterm.Info.Println("Estableciendo circuitos seguros (Ignorando torrc del sistema)...")

	// Escáner de progreso y captura de warnings
	scanner := bufio.NewScanner(stdout)
	bootstrapped := make(chan bool)
	errChan := make(chan string)
	var lastWarnings []string

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			
			// Guardar warnings para reportar si falla
			if strings.Contains(line, "[warn]") || strings.Contains(line, "[err]") {
				lastWarnings = append(lastWarnings, line)
			}

			// Mostrar progreso de forma limpia
			if strings.Contains(line, "Bootstrapped") {
				parts := strings.Split(line, "Bootstrapped")
				if len(parts) > 1 {
					progress := strings.TrimSpace(parts[1])
					pterm.Print(pterm.FgCyan.Sprintf("\r[TOR] Progreso: %s          ", progress))
				}
			}

			if strings.Contains(line, "Bootstrapped 100%") {
				fmt.Println() 
				bootstrapped <- true
				return
			}
		}
		// Si el scanner termina sin bootstrap, algo falló
		errChan <- strings.Join(lastWarnings, "\n")
	}()

	select {
	case <-bootstrapped:
		pterm.Success.Println("Tor conectado exitosamente.")
		return nil
	case warnings := <-errChan:
		tm.Stop()
		return fmt.Errorf("Tor falló al iniciar:\n%s", warnings)
	case <-time.After(120 * time.Second):
		tm.Stop()
		return fmt.Errorf("Tiempo de espera agotado (120s).")
	}
}

func (tm *TorManager) Reload() error {
	if tm.cmd != nil && tm.cmd.Process != nil {
		pterm.Info.Println("Solicitando nueva identidad (NewNYM)...")
		// Reiniciamos para forzar nueva IP (es lo más fiable sin ControlPort)
		tm.Stop()
		return tm.Start()
	}
	return fmt.Errorf("tor no está en ejecución")
}

func (tm *TorManager) CheckTorConnection() (bool, error) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Transport: &http.Transport{Dial: dialer.Dial},
		Timeout:   15 * time.Second,
	}

	resp, err := client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

func (tm *TorManager) Stop() {
	if tm.cmd != nil && tm.cmd.Process != nil {
		pterm.Warning.Println("Cerrando circuitos de Tor...")
		tm.cmd.Process.Kill()
		tm.cmd.Wait()
		os.RemoveAll("/tmp/tor_anongo")
	}
}
