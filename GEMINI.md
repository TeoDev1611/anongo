# GEMINI.md - Contexto de Anongo (Golang Ghost)

## Estado del Proyecto
Este proyecto es un port funcional y optimizado de **AnonGT** (Python) a **Go**. Se ha completado la migración de las funcionalidades principales con un enfoque en **rendimiento paralelo** y **seguridad de red**.

## Ventajas respecto a AnonGT (Python)
1.  **Rendimiento Superior:** El uso de Goroutines permite que el Watchdog de red, el Monitor ARP y el Gestor de Tor se ejecuten en paralelo con un consumo de recursos mínimo, superando la ejecución secuencial o pesada de Python.
2.  **Seguridad Criptográfica Mejorada:** Migración de AES-CBC a **AES-256-GCM**, proporcionando cifrado autenticado y mayor resistencia a ataques de manipulación de datos.
3.  **Binario Único y Estático:** Elimina la necesidad de un intérprete de Python y múltiples dependencias de pip, facilitando la distribución y ejecución en sistemas Linux mínimos.
4.  **Monitoreo de Red de Bajo Nivel:** El uso de `gopacket` para el análisis ARP es significativamente más rápido y eficiente que el uso de Scapy en el port original.
5.  **Gestión de Ciclo de Vida Robusta:** Implementación estricta de `context` y manejo de señales de sistema para asegurar una limpieza total de las reglas de `iptables` y DNS al cerrar el programa.

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
