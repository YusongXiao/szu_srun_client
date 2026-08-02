# SZU Srun Client

一个用于在命令行环境下登陆深大校园网的客户端, 适用于 srun 认证系统 (drcom 认证系统请移步 [login_szu_network](https://github.com/Caterpie771881/login_szu_network))

~~需要 python3 环境, 由于认证逻辑比较复杂, 所以没用 shell 脚本来写~~

最新客户端使用 golang 编写, 无需依赖 python 环境即可实现跨平台

若所在平台支持 python3, 也可以使用 `PyClient` 文件夹下的脚本

# 使用方式

## golang 版客户端

### 快速启动

原理：SRUN 未登录时不会封禁 HTTPS 的 IPv6 站点，所以可以直接走 IPv6 下载脚本和二进制完成一键安装

> [SH](https://github.com/YusongXiao) 已经帮你们构建好了一键脚本（顺便做了下开机自启）

```bash
curl -fsSL https://oss.songhappy.cn/archive/srun.sh | sudo SRUN_USERNAME=你的学号 SRUN_PASSWORD=你的密码 bash
```

### 手动下载
在 [Releases](https://github.com/Caterpie771881/szu_srun_client/releases) 中下载对应平台的二进制文件并运行即可

PS: 运行前请保证二进制文件有可执行权限

## python 版客户端

保证所在机器有 python3 环境, 将本项目中的 [PyClient](./PyClient/) 文件夹复制到机器上

然后执行以下命令

```
cd PyClient & chmod +x main.py & python3 main.py
```
