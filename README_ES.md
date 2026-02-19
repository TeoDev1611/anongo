# Anongo 👻

[![Buy Me A Coffee](https://img.shields.io/badge/Invítame%20a%20un%20café-Donar-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

**Anongo** (Go Ghost) es un port altamente optimizado y centrado en la seguridad del original [**AnonGT**](https://github.com/gt0day/AnonGT) (desarrollado por **gt0day** en Python) rescrito en **Go**. Crea una capa de anonimato transparente al forzar todo el tráfico de red del sistema a través de la red Tor, implementando protecciones avanzadas contra fugas y limpieza forense.

## 🌟 ¿Por qué el Port en Go? (Ventajas sobre Python)

- **Ejecución Paralela:** Utiliza Goroutines de Go para ejecutar el watchdog de red, el monitor ARP y el gestor de Tor de forma concurrente con un consumo mínimo de recursos.
- **Criptografía Mejorada:** Migración de AES-CBC a **AES-256-GCM**, proporcionando cifrado autenticado y mayor resistencia a la manipulación de datos.
- **Binario Estático:** Sin necesidad de un intérprete de Python ni dependencias complejas de `pip`; un solo binario rápido y ligero.
- **Menor Consumo de Recursos:** Uso significativamente más eficiente de memoria y CPU en comparación con la implementación original en Python.
- **Gestión de Red Robusta:** Utiliza `gopacket` para el monitoreo de red de bajo nivel, siendo más rápido y confiable que Scapy.


## 🚀 Cómo Funciona

Anongo crea un "Túnel Fantasma" combinando el **Proxy Transparente de Tor** y la **Redirección Quirúrgica de Iptables**:

1. **Blindaje de Red:** Crea cadenas personalizadas de `iptables` (`ANONGO_NAT` y `ANONGO_FILTER`) para redirigir todo el tráfico TCP al `TransPort` (9040) de Tor y el tráfico DNS al `DNSPort` (5353).
2. **Motor Anti-Fugas:** 
   - **Killswitch IPv6:** Deshabilita automáticamente todo el tráfico IPv6 para evitar fugas comunes.
   - **Filtrado UDP:** Bloquea el tráfico UDP no-Tor (ej. QUIC, STUN) mientras permite la redirección DNS.
   - **Exención de LAN:** Detecta y excluye automáticamente las redes locales (127.0.0.1, 192.168.x.x, etc.) para no perder acceso a tu router o dispositivos locales.
3. **Baja de Privilegios:** Tor se ejecuta bajo un usuario de sistema específico (`tor` o `debian-tor`). Esto permite que Anongo le diga a `iptables`: "Redirige todo EXCEPTO el tráfico que proviene del usuario Tor", evitando bucles de tráfico infinitos.
4. **Monitorización Watchdog:** Una goroutine en segundo plano verifica la integridad de la conexión cada 15 segundos. Si se detectan fugas o eliminación de reglas, re-aplica el blindaje instantáneamente.

## 🏗️ Arquitectura Técnica

El proyecto está organizado en paquetes modulares para asegurar el mantenimiento y el alto rendimiento:

- **`pkg/network`**: Gestiona las reglas de `iptables`, la desactivación de IPv6 y el Watchdog de red.
- **`pkg/tor`**: Controla el ciclo de vida del proceso Tor, el cambio de identidad y la salud de los circuitos.
- **`pkg/security`**: Implementa medidas anti-forenses (limpieza de RAM, truncado de historial, limpieza de logs).
- **`pkg/crypto`**: Proporciona cifrado de alto grado para datos locales (AES-256-GCM).
- **`pkg/i18n`**: Soporte multi-idioma (Inglés/Español).

## 🛠️ Requisitos y Dependencias

### Binarios de Sistema

Debes tener lo siguiente instalado en tu sistema Linux:

- **tor**: El servicio principal de anonimato.
- **iptables / ip6tables**: Para la redirección de red.
- **procps (pkill)**: Para la limpieza de procesos.
- **kmod**: Para asegurar que los módulos de iptables estén cargados.

### Dependencias Go

El proyecto utiliza las siguientes librerías:

- `github.com/pterm/pterm`: Para la interfaz CLI interactiva y el dashboard.
- `github.com/coreos/go-iptables`: Para la gestión quirúrgica del firewall.
- `golang.org/x/net/proxy`: Para verificación segura de SOCKS5.

## 📥 Instalación

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/TeoDev1611/anongo.git
   cd anongo
   ```

2. **Compilar el binario:**
   ```bash
   go build -o anongo ./cmd/anongo/main.go
   ```

## 🎮 Uso

Anongo requiere privilegios de **root** para gestionar las interfaces de red y el estado de los procesos.

### Banderas (Flags) de Línea de Comandos

| Bandera | Descripción | Por Defecto |
| :--- | :--- | :--- |
| `-lang` | Establece el idioma de la interfaz (`en` o `es`). | `es` |
| `-logs` | Activa el guardado del historial de sesión en `anongo_session.log`. | `false` |

### Ejecución Básica

```bash
# Ejecutar en inglés con logs activados
sudo ./anongo -lang en -logs
```

### Opciones del Menú

- **Activar Anonymity:** Asegura la conexión e inicia el túnel.
- **Detener Túnel:** Restaura de forma segura la configuración de red original sin borrar tus reglas personales.
- **Cambiar Identidad:** Reinicia los circuitos de Tor para obtener una nueva IP pública.
- **Chequeo Detallado:** Muestra una tabla comparando tu interfaz pública vs. cifrada.
- **Anti-Forensics:** Borra rastros del sistema, limpia cachés de RAM (`drop_caches`) y trunca el historial de bash/zsh.
- **Limpieza de Emergencia:** Opción failsafe para forzar la restauración de todos los ajustes.

## 🔒 Características de Seguridad

- **Iptables Quirúrgicas:** Utiliza cadenas aisladas. No interferirá con tus reglas de firewall personalizadas.
- **Seguridad de Memoria:** Escrito en Go, eliminando riesgos de desbordamiento de búfer presentes en alternativas basadas en C.
- **Limpieza Forense:** Limpia `/proc/sys/vm/drop_caches`, logs del sistema (`auth.log`, `syslog`) e historiales de sesión para minimizar la huella forense.

## ⚠️ Descargo de Seguridad

**Anongo es una herramienta para la investigación de seguridad y privacidad.** Aunque proporciona un fuerte anonimato, ninguna herramienta es 100% infalible.
- El uso de esta herramienta no te hace inmune al fingerprinting o fugas a nivel de aplicación (ej. plugins de navegador).
- Utiliza siempre un navegador endurecido para la privacidad (como Tor Browser) incluso cuando el túnel esté activo.
- **Los desarrolladores no son responsables de cualquier mal uso o daño causado por esta herramienta.**

## ☕ Soporte

Si encuentras útil **Anongo** y quieres apoyar su desarrollo, ¡puedes invitarme a un café!

[![Buy Me A Coffee](https://img.shields.io/badge/Invítame%20a%20un%20café-Donar-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

## 📄 Licencia
Este proyecto está bajo la licencia **GPL v3 License**.
