# Anongo 👻

**Anongo** (Go Ghost) es un port altamente optimizado y centrado en la seguridad del original **AnonGT** (Python) rescrito en **Go**. Crea una capa de anonimato transparente al forzar todo el tráfico de red del sistema a través de la red Tor, implementando protecciones avanzadas contra fugas y limpieza forense.

## 🚀 Cómo Funciona

Anongo crea un "Túnel Fantasma" combinando el **Proxy Transparente de Tor** y la **Redirección Quirúrgica de Iptables**:
1. **Blindaje de Red:** Crea cadenas personalizadas de `iptables` para redirigir todo el tráfico TCP al `TransPort` (9040) de Tor y el tráfico DNS al `DNSPort` (5353).
2. **Motor Anti-Fugas:** Deshabilita automáticamente IPv6 (fuente común de fugas) y bloquea el tráfico UDP no-Tor (evitando que protocolos como QUIC omitan el proxy).
3. **Baja de Privilegios:** Tor se ejecuta bajo un usuario de sistema específico (`tor` o `debian-tor`). Esto permite que Anongo le diga a `iptables`: "Redirige todo EXCEPTO el tráfico que proviene del usuario Tor", evitando bucles de tráfico infinitos.
4. **Monitorización Watchdog:** Una goroutine en segundo plano verifica la integridad de la conexión cada 15 segundos. Si se detectan fugas o eliminación de reglas, re-aplica el blindaje instantáneamente.

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

Anongo requiere privilegios de **root** para gestionar las interfaces de red.

```bash
# Ejecución básica (Español por defecto)
sudo ./anongo

# Ejecución en Inglés
sudo ./anongo -lang en

# Ejecución con logs activados
sudo ./anongo -logs
```

### Opciones del Menú
- **Activar Anonymity:** Asegura la conexión e inicia el túnel.
- **Detener Túnel:** Restaura de forma segura la configuración de red original sin borrar tus reglas personales.
- **Cambiar Identidad:** Reinicia los circuitos de Tor para obtener una nueva IP pública.
- **Chequeo Detallado:** Muestra una tabla comparando tu interfaz pública vs. cifrada.
- **Anti-Forensics:** Borra rastros del sistema y limpia cachés de RAM.
- **Limpieza de Emergencia:** Opción failsafe para forzar la restauración de todos los ajustes.

## 🔒 Características de Seguridad
- **Iptables Quirúrgicas:** Utiliza cadenas `ANONGO_NAT` y `ANONGO_FILTER`. No romperá tu firewall existente.
- **Seguridad de Memoria:** Escrito en Go, proporcionando mejor gestión de memoria que la versión original en Python.
- **Limpieza Forense:** Limpia `drop_caches` y logs de sesión para minimizar la huella dejada en la máquina.

## 📄 Licencia
Este proyecto está bajo la licencia **GPL v3 License**.
