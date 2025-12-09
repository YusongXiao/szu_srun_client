# SZU SRUN Client (Docker Version)

这是一个用于深圳大学校园网认证（SRUN）的轻量级 Docker 客户端。它经过精简，专注于提供稳定、自动化的登录服务，并包含防掉网机制。

## 注意
使用之前应该先用[Caterpie771881/szu_srun_client](https://github.com/Caterpie771881/szu_srun_client)登录一次，这样才能拉取docker镜像，如果只需要实现上网功能直接使用以上链接即可，本仓库旨在实现**防掉网**

## 功能特点

*   **自动登录**：容器启动时自动进行认证。
*   **防掉网**：内置守护进程，每 30 秒检测一次网络连接，断网自动重连。
*   **Docker 部署**：支持通过环境变量配置，易于部署和管理。
*   **轻量级**：基于 Alpine Linux 构建，镜像体积小。

## 鸣谢

本项目大部分源码来自 [Caterpie771881/szu_srun_client](https://github.com/Caterpie771881/szu_srun_client)，在此基础上进行了 Docker 化和功能精简。

## 用法

使用以下命令启动容器：

```bash
docker run --env username=你的学号 --env password=你的密码 --name srunclient --restart=always --network host -d ccr.ccs.tencentyun.com/songhappy/srunclient
```

### 参数说明

*   `--env username=...`: 设置你的学号。
*   `--env password=...`: 设置你的密码。
*   `--restart=always`: 设置容器自动重启，保证开机自启和崩溃重启。
*   `--network host`: 主机网络模式。
*   `-d`: 后台运行。

