# SZU SRUN Client

这是一个用于深圳大学校园网认证（SRUN）的轻量级客户端，目前保留三种方式：

*   **一键安装版（推荐）**：面向 Linux amd64，原生二进制 + 安装脚本，不依赖 Docker。
*   **从源码安装版**：适合 IPv6 不可达、需要自行编译运行的用户。
*   **Docker 版（Legacy）**：保留原有容器运行方式，适合已经具备 Docker 环境的用户。

## 鸣谢

本项目大部分源码来自 [Caterpie771881/szu_srun_client](https://github.com/Caterpie771881/szu_srun_client)，在此基础上进行了精简和服务化封装。

## 一键安装版（推荐，Linux amd64）

SRUN 未登录时不会封禁 HTTPS 的 IPv6 站点，所以才可以直接走 IPv6 下载脚本完成一键安装。

```bash
curl -fsSL https://oss.songhappy.cn/archive/srun.sh | sudo env SRUN_USERNAME=你的学号 SRUN_PASSWORD=你的密码 bash
```

或：

```bash
curl -fsSL https://oss.songhappy.cn/archive/srun.sh | sudo bash
```

## 运行与停止

```bash
sudo systemctl status srun-client
sudo systemctl stop srun-client
sudo systemctl start srun-client
sudo systemctl restart srun-client
sudo journalctl -u srun-client -f
```

机器重启后会自动启动。

重新执行时还是同一条命令：

```bash
curl -fsSL https://oss.songhappy.cn/archive/srun.sh | sudo env SRUN_USERNAME=你的学号 SRUN_PASSWORD=你的密码 bash
```

## 卸载

```bash
sudo systemctl disable --now srun-client
sudo rm -f /etc/systemd/system/srun-client.service /etc/default/srun-client /opt/srun-client/srunClient
sudo rmdir /opt/srun-client 2>/dev/null || true
sudo systemctl daemon-reload
```

## 从源码安装版

IPv6 不可达时，直接使用 Python3 版本：

[Caterpie771881/szu_srun_client](https://github.com/Caterpie771881/szu_srun_client)

## Docker 版（Legacy）

适合已经有 Docker 环境的用户：

```bash
docker run --env username=你的学号 --env password=你的密码 --name srunclient --restart=always --network host -d ccr.ccs.tencentyun.com/songhappy/srunclient
```

### 参数说明

*   `--env username=...`: 设置你的学号。
*   `--env password=...`: 设置你的密码。
*   `--restart=always`: 设置容器自动重启，保证开机自启和崩溃重启。
*   `--network host`: 主机网络模式。
*   `-d`: 后台运行。
