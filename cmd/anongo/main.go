package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TeoDev1611/anongo/pkg/i18n"
	"github.com/TeoDev1611/anongo/pkg/logger"
	"github.com/TeoDev1611/anongo/pkg/network"
	"github.com/TeoDev1611/anongo/pkg/security"
	"github.com/TeoDev1611/anongo/pkg/tor"

	"github.com/pterm/pterm"
	"golang.org/x/net/proxy"
)

var (
	tunnelActive = false
	currentTorIP = "N/A"
	langFlag     = flag.String("lang", "es", "Language of the application (es/en)")
	logFlag      = flag.Bool("logs", false, "Enable saving session output to anongo_session.log")
	T            i18n.Translation
)

func main() {
	flag.Parse()

	// Cargar Idioma
	T = i18n.Get(*langFlag)

	// Configurar Logs solo si la flag -logs está activa (Desactivado por defecto)
	if *logFlag {
		f, err := os.OpenFile("anongo_session.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err == nil {
			multi := io.MultiWriter(os.Stdout, f)
			pterm.SetDefaultOutput(multi)
			fmt.Fprintf(f, "\n--- SESSION START: %s (%s) ---\n", time.Now().Format(time.RFC822), *langFlag)
		}
	}

	fmt.Print("\033[H\033[2J")

	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("ANO", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("NGO", pterm.NewStyle(pterm.FgMagenta))).
		Render()

	pterm.Println(pterm.FgGray.Sprint(T.AppDescription))
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	tm := tor.NewTorManager()
	nm, _ := network.NewNetworkManager(tm)
	af := security.NewAntiForensics()

	// Limpieza inicial
	nm.StopAnonymity()
	tm.Stop()

	go func() {
		<-sigs
		fmt.Println()
		pterm.Warning.Println(T.ExitSignalDetected + " " + T.RestoringSystem)
		nm.StopAnonymity()
		tm.Stop()
		cancel()
		pterm.Success.Println(T.SessionTerminated)
		os.Exit(0)
	}()

	for {
		renderStatus()

		selectedOption, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText(T.MenuTitle).
			WithOptions(T.MenuOptions).
			Show()

		switch selectedOption {
		case T.MenuOptions[0]: // Activar
			if tunnelActive {
				pterm.Warning.Println("Tunnel already active.")
				continue
			}

			if err := tm.Start(); err != nil {
				logger.LogError(err, "Tor Startup")
				pterm.Error.Printf("%s: %v\n", T.TorStartupFail, err)
				continue
			}

			if err := nm.StartAnonymity(ctx); err != nil {
				logger.LogError(err, "Network Activation")
				tm.Stop()
				pterm.Error.Printf("%s: %v\n", T.NetActivationFail, err)
				continue
			}

			tunnelActive = true
			spinner, _ := pterm.DefaultSpinner.Start(T.ValidatingIdentity)
			currentTorIP = getIPThroughTor()
			spinner.Success(fmt.Sprintf(T.SystemSecured, currentTorIP))

		case T.MenuOptions[1]: // Detener
			nm.StopAnonymity()
			tm.Stop()
			tunnelActive = false
			currentTorIP = "N/A"
			pterm.Success.Println(T.NetRestored)

		case T.MenuOptions[2]: // Cambiar IP
			if !tunnelActive {
				pterm.Error.Println(T.ProtectedText + "?")
				continue
			}
			nm.SetMaintenance(true)
			pterm.Info.Println(T.RestartingTor)
			tm.Stop()
			time.Sleep(1 * time.Second)
			if err := tm.Start(); err != nil {
				pterm.Error.Printf("%s: %v\n", T.TorStartupFail, err)
			} else {
				currentTorIP = getIPThroughTor()
				pterm.Success.Printf(T.NewIPAssigned+"\n", currentTorIP)
			}
			nm.SetMaintenance(false)

		case T.MenuOptions[3]: // Chequeo detallado
			displayIPTable()

		case T.MenuOptions[4]: // Anti-Forensics
			pterm.Info.Println(T.PurginTraces)
			af.CleanSystemTraces()
			af.WipeMemory()
			pterm.Success.Println(T.TracesPurged)

		case T.MenuOptions[5]: // Limpieza Failsafe
			pterm.Warning.Println(T.EmergencyCleanup + "...")
			nm.StopAnonymity()
			tm.Stop()
			tunnelActive = false
			currentTorIP = "N/A"
			pterm.Success.Println(T.CleanupSuccess)

		case T.MenuOptions[6]: // Salir
			nm.StopAnonymity()
			tm.Stop()
			pterm.FgMagenta.Println("\n[!] " + T.SessionTerminated)
			return
		}
	}
}

func renderStatus() {
	statusText := pterm.FgRed.Sprint(T.ExposedText)
	if tunnelActive {
		statusText = pterm.FgGreen.Sprint(T.ProtectedText)
	}

	panelContent := fmt.Sprintf(
		"%s: %s  |  %s: %s  |  %s: %s",
		T.StatusLabel,
		statusText,
		T.TorExitLabel,
		pterm.FgCyan.Sprint(currentTorIP),
		T.TimeLabel,
		time.Now().Format("15:04:05"),
	)

	pterm.DefaultBox.WithTitle("System Dashboard").Println(panelContent)
}

func getIPThroughTor() string {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		return "Error Proxy"
	}
	client := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}, Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "N/A"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func displayIPTable() {
	pterm.DefaultSection.Println(T.NetworkIntegrity)

	resp, err := http.Get("https://api.ipify.org")
	directIP := "Offline"
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		directIP = string(body)
		resp.Body.Close()
	}

	pterm.DefaultTable.WithData(pterm.TableData{
		{"Layer", "Identification", "Security"},
		{T.PublicInterface, directIP, pterm.FgRed.Sprint(T.Visible)},
		{T.TorGhostLayer, currentTorIP, pterm.FgGreen.Sprint(T.Encrypted)},
	}).Render()
}
