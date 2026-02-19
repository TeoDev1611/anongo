package security

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pterm/pterm"
)

type AntiForensics struct{}

func NewAntiForensics() *AntiForensics {
	return &AntiForensics{}
}

// CleanSystemTraces elimina rastros de actividad en el sistema
func (af *AntiForensics) CleanSystemTraces() {
	pterm.Info.Println("Iniciando limpieza de rastro del sistema...")

	// 1. Limpiar historical de Bash y Zsh
	home, _ := os.UserHomeDir()
	histFiles := []string{
		filepath.Join(home, ".bash_history"),
		filepath.Join(home, ".zsh_history"),
	}

	for _, file := range histFiles {
		if err := os.Truncate(file, 0); err == nil {
			pterm.Success.Printf("Historical limpiado: %s", file)
		}
	}

	// 2. Limpiar logs del sistema (require root)
	systemLogs := []string{
		"/var/log/auth.log",
		"/var/log/syslog",
		"/var/log/messages",
	}

	for _, log := range systemLogs {
		if _, err := os.Stat(log); err == nil {
			exec.Command("truncate", "-s", "0", log).Run()
			pterm.Success.Printf("Log vaciado: %s", log)
		}
	}

	// 3. Limpiar archivos temporales de Anongo
	os.RemoveAll("/tmp/tor_anongo")
	pterm.Success.Println("Archivos temporales eliminados.")

	// 4. Limpiar caché de DNS
	exec.Command("systemd-resolve", "--flush-caches").Run()
	pterm.Success.Println("Caché DNS purgada.")
}

// WipeMemory (Placeholder para lógica de sobreescritura de RAM si fuera necesario)
func (af *AntiForensics) WipeMemory() {
	pterm.Info.Println("Sincronizando y liberando buffers de memoria (Drop Caches)...")
	exec.Command("sync").Run()
	// Liberar pagecache, dentries e inodes (require root)
	exec.Command("sh", "-c", "echo 3 > /proc/sys/vm/drop_caches").Run()
	pterm.Success.Println("Memoria RAM sincronizada y caches liberadas.")
}
