package logger

import (
	"fmt"
	"os"
	"time"
)

const logFile = "anongo.log"

// LogError guarda un error con fecha y hora en el archivo de log
func LogError(err error, context string) {
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] ERROR [%s]: %v", timestamp, context, err)
	f.WriteString(logMsg)
}

// ReadLogs devuelve el contenido del archivo de logs
func ReadLogs() (string, error) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ClearLogs vacía el archivo de logs
func ClearLogs() error {
	return os.WriteFile(logFile, []byte(""), 0o644)
}
