# AnonGT Porting Guide: Python to Go (Golang)

Este documento detalla la lógica interna de **AnonGT (Anonymous Ghost)** para facilitar su migración desde Python a Go. El objetivo es mantener la funcionalidad de anonimización total mientras se aprovecha la concurrencia y el rendimiento de Go.

---

## 1. Arquitectura del Sistema

El proyecto original está dividido en módulos de script (`core/scripts/`) y fuentes auxiliares (`core/sources/`). En Go, se recomienda una estructura de paquetes:

- `pkg/tor`: Gestión de servicios y túneles.
- `pkg/network`: Configuración de iptables, MAC y DNS.
- `pkg/crypto`: Cifrado AES y derivación de claves.
- `pkg/security`: Anti-MITM (Scapy replacement) y Keylogger.
- `cmd/anongt`: Punto de entrada (Main CLI).

---

## 2. Desglose de Funcionalidades y Equivalencias

### A. Modo Anónimo (Core Logic)
**Python:** Usa `subprocess.run` y `os.system` para manipular `iptables`, `sysctl` y servicios del sistema.
**Go (Port):** 
- Usar el paquete estándar `os/exec` para comandos externos.
- Para `iptables`, se recomienda la biblioteca `github.com/coreos/go-iptables`.
- **Lógica de Redirección:** Se debe configurar el puerto transparente de Tor (TransPort).
  - *Input:* Capturar todo el tráfico TCP/UDP y redirigir a `127.0.0.1:9040`.
  - *DNS:* Redirigir puerto 53 a `127.0.0.1:5353` (Tor DNS).

### B. Anti-MITM (ARP Poisoning Detection)
**Python:** Usa `scapy` para sniffear y enviar paquetes ARP.
**Go (Port):** 
- Biblioteca: `github.com/google/gopacket`.
- **Proceso:** 
  1. Escuchar paquetes ARP en modo promiscuo.
  2. Al recibir un `ARP Reply`, verificar si la IP del emisor coincide con la MAC real (haciendo un query previo).
  3. Si hay discrepancia, ejecutar comando de bloqueo (iptables).

### C. Cifrado AES-256
**Python:** Usa `cryptography` (PBKDF2HMAC + AES CBC).
**Go (Port):** 
- Biblioteca: Estándar `crypto/aes`, `crypto/cipher` y `golang.org/x/crypto/pbkdf2`.
- **Importante:** Go no añade padding PKCS7 automáticamente en el modo CBC; debe implementarse manualmente o usar el modo GCM (más moderno y seguro).

### D. Keylogger Cifrado
**Python:** Usa la librería `keyboard` y `fernet`.
**Go (Port):** 
- Biblioteca: `github.com/robotn/gohook` o `github.com/kindlyfire/go-keylogger`.
- **Estructura:** Usar una **Goroutine** para capturar eventos de teclado en segundo plano y un canal (`chan`) para enviar las teclas al módulo de cifrado.

### E. Dark Web & OnionShare
**Python:** Scraping con `BeautifulSoup` e interfaz con `onionshare-cli`.
**Go (Port):**
- **Scraping:** `github.com/PuerkitoBio/goquery`.
- **Tor Native:** Considerar usar `github.com/cretz/bine` para interactuar con Tor de forma nativa desde Go sin depender tanto de binarios externos.

---

## 3. Desafíos Técnicos en Go

1.  **Privilegios de Root:** Casi todas las funciones (iptables, gopacket, macchanger) requieren privilegios de superusuario. Go debe verificar el `os.Geteuid() == 0` al inicio.
2.  **Manejo de Señales:** Al cerrar el programa (Ctrl+C), es CRÍTICO revertir los cambios de red (limpiar iptables, restaurar DNS). Usar `os/signal` para capturar el cierre y ejecutar una función de "Cleanup".
3.  **Cross-Compilation:** Aunque Go es portable, este proyecto depende fuertemente de herramientas de Linux (`iptables`, `tor`, `macchanger`). El binario resultante solo será funcional en sistemas Linux (especialmente Debian/Kali).

---

## 4. Stack de Bibliotecas Recomendado (Go)

| Función | Biblioteca Go |
| :--- | :--- |
| CLI / UI | `github.com/spf13/cobra` o `github.com/pterm/pterm` |
| Iptables | `github.com/coreos/go-iptables/iptables` |
| Networking (Low Level) | `github.com/google/gopacket` |
| Crypto | `crypto/aes`, `crypto/cipher` |
| HTTP Client (Tor) | `net/http` con un Custom Transport SOCKS5 |
| Keylogging | `github.com/robotn/gohook` |

---

## 5. Roadmap Sugerido para el Port

1.  **Fase 1:** Implementar el paquete `network` (manejo de iptables y DNS).
2.  **Fase 2:** Implementar la lógica de `tor` (configuración del archivo `torrc` y control del servicio).
3.  **Fase 3:** Integrar `gopacket` para el monitor Anti-MITM.
4.  **Fase 4:** Desarrollar el sistema de cifrado y keylogger.
5.  **Fase 5:** Crear la interfaz de usuario con menús interactivos.
