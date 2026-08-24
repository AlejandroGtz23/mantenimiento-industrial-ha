# Despliegue demostrable en AlmaLinux 10

## 1. Crear la máquina virtual

1. Descarga la ISO mínima de AlmaLinux 10 y crea una VM de **2 vCPU, 4 GB RAM y 30 GB de disco**.
2. En VirtualBox entra en **Configuración > Red > Adaptador 1** y selecciona **Adaptador puente**; elige la tarjeta física que usa la red local. En Hyper-V crea/usa un *Virtual Switch* de tipo **External** y asígnalo a la VM.
3. Arranca desde la ISO, selecciona instalación mínima, zona horaria y contraseña de `root`. No instales entorno gráfico.
4. Después de iniciar, confirma que la VM obtuvo una IP de la red local: `ip -4 addr`. Desde otra PC se debe poder hacer `ping IP_DE_LA_VM`.

## 2. Preparar AlmaLinux y Docker

Ejecuta como `root` o con `sudo`:

```bash
sudo dnf update -y
sudo dnf install -y dnf-plugins-core firewalld git
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload

sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
```

Cierra e inicia sesión para aplicar el grupo `docker` y valida con `docker compose version`.
SELinux debe permanecer en modo *Enforcing*: las monturas locales de este proyecto ya usan el sufijo `:Z`, que etiqueta el contenido para los contenedores.

## 3. Levantar el proyecto

```bash
git clone URL_DEL_REPOSITORIO mantenimiento-ha
cd mantenimiento-ha
cp .env.docker.example .env.docker
# Edita .env.docker y cambia las claves; no lo subas al repositorio.
docker compose --env-file .env.docker up -d --build
docker compose ps
docker compose logs -f nginx backend-1 backend-2
```

> Si ya hay un archivo `.env` del desarrollo local, no lo reemplaces. El comando con `--env-file .env.docker` aísla las claves del despliegue.

Accede desde cualquier equipo de la LAN a `http://IP_DE_LA_VM/`. La instalación limpia permite entrar a la web con el usuario definido en `BOOTSTRAP_ADMIN_USER` y su contraseña correspondiente; después crea técnicos y máquinas desde el panel. El usuario inicial sólo permite arrancar una base Firebird vacía: si se migra la base existente, las cuentas de `PROFILE` conservan prioridad.

## 4. Demostrar el balanceo y Redis

Abre `http://IP_DE_LA_VM/balance-test/` y recarga varias veces: el nombre `demo-1`/`demo-2` y el ID de contenedor cambian porque Nginx alterna el upstream. También puedes mostrar `docker compose ps` y el bloque `upstream` de `deploy/nginx/nginx.conf`. Para demostrar el estado compartido, inicia sesión, revisa `docker compose exec redis redis-cli KEYS 'maintenance:session:*'`, detén una réplica con `docker compose stop backend-1` y navega/consulta la API con el mismo JWT: `backend-2` valida el token almacenado en Redis.

## Operación breve

```bash
docker compose --env-file .env.docker down       # detiene sin borrar datos
docker compose --env-file .env.docker up -d      # vuelve a iniciar
docker compose --env-file .env.docker logs -f    # diagnóstico
```

Los datos persistentes viven en `deploy/data/firebird`, `deploy/data/redis` y `deploy/data/uploads`. No uses `down -v` para esta entrega, ya que eliminaría los volúmenes nombrados; las carpetas enlazadas sobreviven de todas formas.
