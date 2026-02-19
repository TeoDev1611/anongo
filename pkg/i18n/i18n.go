package i18n

import "fmt"

type Translation struct {
	AppDescription  string
	StatusLabel     string
	TorExitLabel    string
	TimeLabel       string
	ExposedText     string
	ProtectedText   string
	MenuTitle       string
	MenuOptions     []string
	SecuringConn    string
	TorStartupFail  string
	NetActivationFail string
	ValidatingIdentity string
	SystemSecured   string
	NetRestored     string
	RestoringNet    string
	RestartingTor   string
	NewIPAssigned   string
	NetworkIntegrity string
	PublicInterface string
	TorGhostLayer   string
	Visible         string
	Encrypted       string
	PurginTraces    string
	WipingMemory    string
	TracesPurged    string
	EmergencyCleanup string
	CleanupSuccess   string
	SessionTerminated string
	ExitSignalDetected string
	RestoringSystem string
	CleaningTraces  string
	LanguageChanged string
}

var es = Translation{
	AppDescription:  "— EL PUERTO DEL FANTASMA ANÓNIMO —",
	StatusLabel:     "ESTADO",
	TorExitLabel:    "SALIDA TOR",
	TimeLabel:       "HORA",
	ExposedText:     "EXPUESTO",
	ProtectedText:   "PROTEGIDO",
	MenuTitle:       "Control de Terminal",
	MenuOptions: []string{
		"Activar Anonymity (Tor + Iptables)",
		"Detener Túnel (Limpiar Reglas)",
		"Cambiar Identidad (Nueva IP)",
		"Chequeo de Red Detallado",
		"Anti-Forensics (Borrar Rastros)",
		"Limpieza de Emergencia (Failsafe)",
		"Salir",
	},
	SecuringConn:    "Asegurando Conexión...",
	TorStartupFail:  "Fallo al iniciar Tor",
	NetActivationFail: "Error en reglas de red",
	ValidatingIdentity: "Validando Identidad Anónima...",
	SystemSecured:   "Sistema Asegurado. Tu IP es: %s",
	NetRestored:     "Red restaurada con éxito.",
	RestoringNet:    "Restaurando red...",
	RestartingTor:   "Reiniciando circuitos de Tor...",
	NewIPAssigned:   "Nueva IP asignada: %s",
	NetworkIntegrity: "Chequeo de Integridad de Red",
	PublicInterface: "Interfaz Pública",
	TorGhostLayer:   "Capa Fantasma Tor",
	Visible:         "VISIBLE",
	Encrypted:       "CIFRADO",
	PurginTraces:    "Purgando rastros y logs del sistema...",
	WipingMemory:    "Limpiando memoria RAM...",
	TracesPurged:    "Huellas borradas.",
	EmergencyCleanup: "Limpieza de Emergencia (Failsafe)",
	CleanupSuccess:   "Limpieza completada.",
	SessionTerminated: "Sesión terminada. Mantente seguro.",
	ExitSignalDetected: "Señal de salida detectada.",
	RestoringSystem: "Restaurando sistema...",
	CleaningTraces:  "Limpiando huellas...",
	LanguageChanged: "Idioma cambiado a Español",
}

var en = Translation{
	AppDescription:  "— THE ANONYMOUS GHOST PORT —",
	StatusLabel:     "STATUS",
	TorExitLabel:    "TOR EXIT",
	TimeLabel:       "TIME",
	ExposedText:     "EXPOSED",
	ProtectedText:   "PROTECTED",
	MenuTitle:       "Terminal Control",
	MenuOptions: []string{
		"Activate Anonymity (Tor + Iptables)",
		"Stop Tunnel (Clear Rules)",
		"Change Identity (New IP)",
		"Detailed Network Check",
		"Anti-Forensics (Wipe Traces)",
		"Emergency Cleanup (Failsafe)",
		"Exit",
	},
	SecuringConn:    "Securing Connection...",
	TorStartupFail:  "Tor startup failed",
	NetActivationFail: "Network rules error",
	ValidatingIdentity: "Validating Anonymous Identity...",
	SystemSecured:   "System Secured. Your IP is: %s",
	NetRestored:     "Network restored successfully.",
	RestoringNet:    "Restoring network...",
	RestartingTor:   "Restarting Tor circuits...",
	NewIPAssigned:   "New IP assigned: %s",
	NetworkIntegrity: "Network Integrity Check",
	PublicInterface: "Public Interface",
	TorGhostLayer:   "Tor Ghost Layer",
	Visible:         "VISIBLE",
	Encrypted:       "ENCRYPTED",
	PurginTraces:    "Purging system traces and logs...",
	WipingMemory:    "Wiping RAM memory...",
	TracesPurged:    "Traces wiped.",
	EmergencyCleanup: "Emergency Cleanup (Failsafe)",
	CleanupSuccess:   "Cleanup completed.",
	SessionTerminated: "Session terminated. Stay safe.",
	ExitSignalDetected: "Exit signal detected.",
	RestoringSystem: "Restoring system...",
	CleaningTraces:  "Cleaning traces...",
	LanguageChanged: "Language changed to English",
}

func Get(lang string) Translation {
	if lang == "en" {
		return en
	}
	return es
}
