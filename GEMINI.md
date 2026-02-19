# GEMINI.md - Contexto de Anongo (Golang Ghost)

## Estado del Proyecto
Este proyecto es un port funcional y optimizado de **AnonGT** (Python) a **Go**. Se ha completado la migración de las funcionalidades principales con un enfoque en **rendimiento paralelo** y **seguridad de red**.

## Arquitectura de Goroutines
Anongo utiliza un sistema multihilo coordinado mediante `context` para gestionar el ciclo de vida de los servicios:
1.  **Watchdog (network.go):** Monitorea cada 7 segundos que las reglas de `iptables` sigan activas y que la IP pública sea una IP de Tor.
2.  **Monitor ARP (security.go):** Escucha paquetes ARP en modo promiscuo (`gopacket`) para detectar ataques MITM sin bloquear la interfaz de usuario.
3.  **Gestor de Tor (tor.go):** Controla el proceso binario de Tor (`--TransPort 9040`).

## Protocolos de Seguridad Implementados
- **AES-256-GCM:** Cifrado autenticado para proteger datos, superior al CBC original.
- **PBKDF2 SHA-256:** Derivación de claves con 100,000 iteraciones para resistencia contra fuerza bruta.
- **Anti-Forensics:** Limpieza de `drop_caches` y truncado de historiales para evitar rastreo post-sesión.

## Notas para Desarrollo Futuro
- **Manejo de Root:** El programa DEBE ejecutarse como root para las llamadas a `pcap` e `iptables`.
- **Compatibilidad:** Exclusivamente Linux (Debian/Arch). No portar a Windows sin rediseñar el motor de red.
- **Dependencies:** `pterm` (UI), `gopacket` (Red), `go-iptables` (Firewall).

## Cómo Contribuir
Sigue el patrón de paquetes en `pkg/` y utiliza siempre `context.Context` para asegurar que las nuevas funcionalidades se detengan correctamente al salir del programa.
