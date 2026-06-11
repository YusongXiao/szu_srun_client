#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="srun-client"
INSTALL_DIR="/opt/srun-client"
BIN_PATH="${INSTALL_DIR}/srunClient"
ENV_FILE="/etc/default/srun-client"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
SCRIPT_URL="https://oss.songhappy.cn/archive/srun.sh"
DOWNLOAD_URL="${SRUN_DOWNLOAD_URL:-https://oss.songhappy.cn/archive/srunClient_linux_x64}"
TTY_PATH="/dev/tty"

log() {
	printf '[srun-install] %s\n' "$*"
}

fail() {
	printf '[srun-install] %s\n' "$*" >&2
	exit 1
}

has_cmd() {
	command -v "$1" >/dev/null 2>&1
}

require_root() {
	if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
		fail "请使用 root 运行，例如: curl -fsSL ${SCRIPT_URL} | sudo bash"
	fi
}

check_platform() {
	if [[ "$(uname -s)" != "Linux" ]]; then
		fail "当前只支持 Linux"
	fi

	case "$(uname -m)" in
		x86_64|amd64) ;;
		*) fail "当前只支持 Linux amd64" ;;
	esac

	if ! has_cmd systemctl; then
		fail "未检测到 systemd，当前最小化方案无法自动托管服务"
	fi

	if ! has_cmd curl && ! has_cmd wget; then
		fail "需要 curl 或 wget 来下载二进制文件"
	fi
}

download_file() {
	local url="$1"
	local dest="$2"

	if has_cmd curl; then
		curl -fsSL "$url" -o "$dest"
		return
	fi

	wget -qO "$dest" "$url"
}

prompt_nonempty() {
	local prompt_text="$1"
	local secret_mode="$2"
	local value=""

	if [[ ! -r "$TTY_PATH" ]]; then
		fail "当前没有可用终端，无法交互输入账号密码。请连接终端后重试，或通过 SRUN_USERNAME / SRUN_PASSWORD 传入。"
	fi

	while [[ -z "$value" ]]; do
		printf '%s' "$prompt_text" > "$TTY_PATH"
		if [[ "$secret_mode" == "yes" ]]; then
			IFS= read -r -s value < "$TTY_PATH"
			printf '\n' > "$TTY_PATH"
		else
			IFS= read -r value < "$TTY_PATH"
		fi
	done

	printf '%s' "$value"
}

resolve_credentials() {
	if [[ -f "$ENV_FILE" ]]; then
		# shellcheck disable=SC1090
		. "$ENV_FILE"
	fi

	USERNAME="${SRUN_USERNAME:-${username:-}}"
	PASSWORD="${SRUN_PASSWORD:-${password:-}}"

	if [[ -z "$USERNAME" ]]; then
		USERNAME="$(prompt_nonempty '请输入校园网账号: ' no)"
	fi

	if [[ -z "$PASSWORD" ]]; then
		PASSWORD="$(prompt_nonempty '请输入校园网密码: ' yes)"
	fi
}

escape_single_quotes() {
	printf '%s' "$1" | sed "s/'/'\\''/g"
}

write_env_file() {
	local escaped_username
	local escaped_password

	escaped_username="$(escape_single_quotes "$USERNAME")"
	escaped_password="$(escape_single_quotes "$PASSWORD")"

	install -d -m 0755 "$(dirname "$ENV_FILE")"
	umask 077
	cat > "$ENV_FILE" <<EOF
username='${escaped_username}'
password='${escaped_password}'
EOF
	chmod 600 "$ENV_FILE"
}

install_binary() {
	local tmp_bin
	tmp_bin="$(mktemp)"
	trap 'rm -f "$tmp_bin"' EXIT

	log "下载 Linux amd64 客户端..."
	download_file "$DOWNLOAD_URL" "$tmp_bin"

	install -d -m 0755 "$INSTALL_DIR"
	install -m 0755 "$tmp_bin" "$BIN_PATH"

	rm -f "$tmp_bin"
	trap - EXIT
}

write_service_file() {
	cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=SZU SRUN Client
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
EnvironmentFile=$ENV_FILE
ExecStart=$BIN_PATH
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
}

enable_service() {
	log "注册并启动 systemd 服务..."
	systemctl daemon-reload
	systemctl enable "$SERVICE_NAME" >/dev/null
	systemctl restart "$SERVICE_NAME"
}

show_summary() {
	cat <<EOF

安装完成。

首次和后续都可以使用同一条命令：
curl -fsSL $SCRIPT_URL | sudo bash

常用命令：
	启动服务: sudo systemctl start $SERVICE_NAME
	停止服务: sudo systemctl stop $SERVICE_NAME
	重启服务: sudo systemctl restart $SERVICE_NAME
	查看状态: sudo systemctl status $SERVICE_NAME
	查看日志: sudo journalctl -u $SERVICE_NAME -f

服务已设置为开机自启。
如需更新账号密码，重新执行脚本并传入 SRUN_USERNAME / SRUN_PASSWORD 即可。
EOF
}

main() {
	require_root
	check_platform
	resolve_credentials
	install_binary
	write_env_file
	write_service_file
	enable_service
	show_summary
}

main "$@"