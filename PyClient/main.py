import time
import json
import requests
import sys
from typing import Union
from enum import Enum
from urllib.parse import urlencode
from encryptlib import hmd5
from encryptlib import sha1
from encryptlib import chkstr
from encryptlib import info_


version = sys.version_info
if version < (3, 0):
    print('The current version is not supported, you need to use python3')
    sys.exit()

# rad_user_info: 获取用户信息的地址
# get_challenge: 获取加密 token 的地址
# srun_portal: 身份认证地址
url = {
    'rad_user_info': 'http://net.szu.edu.cn/cgi-bin/rad_user_info',
    'get_challenge': 'http://net.szu.edu.cn/cgi-bin/get_challenge',
    'srun_portal'  : 'http://net.szu.edu.cn/cgi-bin/srun_portal'
}

# jsonp 的标志, 一般不用变动
callback = "jQueryCallback"
# 用于模拟浏览器访问的 UA 头, 一般不用变动
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

# 加密常量, 一般不用变动
TYPE = "1"
N = "200"
ENC = 'srun_bx1'
ACID = "12"

BANNER = """
   ________  __  __  ____                ________          __ 
  / __/_  / / / / / / __/_____ _____    / ___/ (_)__ ___  / /_
 _\ \  / /_/ /_/ / _\ \/ __/ // / _ \  / /__/ / / -_) _ \/ __/
/___/ /___/\____/ /___/_/  \_,_/_//_/  \___/_/_/\__/_//_/\__/ 
"""


class Mode(Enum):
    Login = '1'
    Logout = '2'


def get_ip(callback: str, path: str) -> Union[str, None]:
    params = urlencode({"callback": callback})
    try:
        print("[*] 正在尝试自动获取登录 ip...")
        resp = requests.get(
            url=f"{path}?{params}",
            headers={"User-Agent": UA}
        )
        info = resp.text[len(callback)+1:-1]
        ip = json.loads(info)['online_ip']
    except requests.RequestException:
        print(f"[!] 访问 {path} 的过程中出现网络问题, 请检查:")
        print("1. 网络环境是否正常")
        print("2. 配置的 URL 是否正确")
    except json.JSONDecodeError:
        print("[!] 成功获取返回值, 但是 json 解析失败")
        print("[*] 正在检查是否存在 callback 信息...")
        if resp.text[:len(callback)] == callback:
            print("[!] 正常获取了 callback 信息, 请人工检查返回内容是否符合 json 格式:")
            print(resp.text[:100] + ('' if len(resp.text) < 100 else '...'))
        else:
            print("[!] 没有 callback 信息, 请检查配置的 URL 是否正确")
    except KeyError:
        print("[!] 成功解析了返回的数据, 但是不存在 online_ip 值")
        print("[*] 可能 srun 服务发生了异常, 正在检查报错码...")
        try:
            error = json.loads(info)['error_msg']
            print(f"[!] srun 错误码为: {error}")
        except:
            print("[!] 获取报错码失败, 请人工核验 json 信息:")
            print(json.loads(info))
    # 其他类型的异常就直接抛出, 没法预知并处理所有异常
    except Exception as e:
        raise e
    else:
        print(f"[*] 成功获取登录 ip: {ip}")
        return ip


def get_challenge(
        callback: str,
        username: str,
        ip: str,
        path: str
    ) -> Union[str, None]:
    params = urlencode({
        "callback": callback,
        "username": username,
        "ip": ip
    })
    try:
        print("[*] 正在获取加密 token...")
        resp = requests.get(
            url=f"{path}?{params}",
            headers={"User-Agent": UA}
        )
        get_challenge = resp.text[len(callback)+1:-1]
        challenge = json.loads(get_challenge)['challenge']
    except requests.RequestException:
        print(f"[!] 访问 {path} 的过程中出现网络问题, 请检查:")
        print("1. 网络环境是否正常")
        print("2. 配置的 URL 是否正确")
    except json.JSONDecodeError:
        print("[!] 成功获取返回值, 但是 json 解析失败")
        print("[*] 正在检查是否存在 callback 信息...")
        if resp.text[:len(callback)] == callback:
            print("[!] 正常获取了 callback 信息, 请人工检查返回内容是否符合 json 格式:")
            print(resp.text[:100] + ('' if len(resp.text) < 100 else '...'))
        else:
            print("[!] 没有 callback 信息, 请检查配置的 URL 是否正确")
    except KeyError:
        print("[!] 成功解析了返回的数据, 但是不存在 challenge 值")
        print("[*] 可能 srun 服务发生了异常, 正在检查报错码...")
        try:
            error = json.loads(get_challenge)['error_msg']
            print(f"[!] srun 错误码为: {error}")
        except:
            print("[!] 获取报错码失败, 请人工核验 json 信息:")
            print(json.loads(get_challenge))
    # 其他类型的异常就直接抛出, 没法预知并处理所有异常
    except Exception as e:
        raise e
    else:
        print(f"[*] 成功获取加密 token: {challenge}")
        return challenge


def srun_portal_login(
        callback: str,
        username: str,
        password: str,
        path: str,
        token: str,
        ip: str,
        os: str,
    ) -> None:
    try:
        print("[*] 正在加密用户信息...")
        hmd5_password = hmd5(password, token)
        info = info_({
            "username": username,
            "password": password,
            "ip": ip,
            "acid": ACID,
            "enc_ver": ENC
        }, token)
        chksum = sha1(chkstr(
            token, username, hmd5_password, ACID, ip, N, TYPE, info
        ))
        print("[*] 已完成用户信息加密, 准备进入身份认证")
    except Exception as e:
        print("[!] 加密失败, 即将抛出异常:")
        raise e
    params = urlencode({
        "action": 'login',
        "callback": callback,
        "username": username,
        "password": '{MD5}' + hmd5_password,
        "os": os,
        "name": os,
        "nas_ip": '',
        "double_stack": 0,
        "chksum": chksum,
        "info": info,
        "ac_id": ACID,
        "ip": ip,
        "n": N,
        "type": TYPE,
        "captchaVal": '',
        '_': int(time.time() * 1000)
    })
    try:
        resp = requests.get(
            url=f"{path}?{params}",
            headers={"User-Agent": UA}
        )
        result = json.loads(resp.text[len(callback)+1:-1])
        assert result.get('res') == 'ok'
        print("[*] 登录成功")
    except requests.RequestException:
        print(f"[!] 访问 {path} 的过程中出现网络问题, 请检查:")
        print("1. 网络环境是否正常")
        print("2. 配置的 URL 是否正确")
    except Exception:
        print("[!] 登录失败")
        print("[*] 正在检查 srun 报错码...")
        try:
            error = result['error_msg']
            print(f"[*] srun 报错码: {error}")
        except:
            print("[!] 获取报错码失败, 请人工核验:")
            print(resp.text[:100] + ('' if len(resp.text) < 100 else '...'))


def srun_portal_logout(
        callback: str,
        username: str,
        ip: str,
        path: str
    ) -> None:
    params = urlencode({
        "action": "logout",
        "callback": callback,
        "username": username,
        "ip": ip
    })
    resp = requests.get(
        f"{path}?{params}"
    )
    result = json.loads(resp.text[len(callback)+1:-1])
    if result.get('res') == 'ok':
        print("[*] 登出成功")
    else:
        print("[*] 登出失败")


def logout():
    username = input("[+] 请输入您的学号: ")
    ip = get_ip(callback, url['rad_user_info'])
    assert ip != None, "获取 IP 失败, 请检查配置的 URL"
    srun_portal_logout(callback, username, ip, url['srun_portal'])


def login():
    import getpass
    username = input("[+] 请输入您的学号: ")
    password = getpass.getpass("[+] 请输入您的密码: ")
    auto_ip = input("[?] 是否需要自动获取登录 ip [Y/n]: ") or 'Y'
    if auto_ip.lower() == 'y':
        ip = get_ip(callback, url['rad_user_info'])
        assert ip != None, "获取 IP 失败, 请检查配置的 URL 或手动指定 ip"
    else:
        ip = input("[+] 请输入您的登录ip: ")
    auto_os = input("[?] 是否需要指定设备os(默认为 Windows) [y/N]: ") or 'N'
    if auto_os.lower() == 'y':
        os = input("[+] 请输入设备os: ")
    else:
        os = "Windows"

    token = get_challenge(
        callback, username,
        ip, url["get_challenge"]
    )
    if token:
        srun_portal_login(
            callback, username,
            password, url["srun_portal"],
            token, ip, os
        )


if __name__ == "__main__":
    print(BANNER)
    print("[1]登录 [2]登出 [other]退出")
    mode = input("[+] 请选择工作模式: ")
    if mode == Mode.Login.value:
        login()
    elif mode == Mode.Logout.value:
        logout()
    else:
        print("[*] BYE")
