import hashlib
from typing import Union


def hmd5(msg: str, key: str) -> str:
    """hmacmd5 哈希"""
    import hmac
    msg = msg.encode()
    key = key.encode()
    return hmac.new(key, msg, hashlib.md5).hexdigest()


def sha1(msg: str) -> str:
    """sha1 哈希"""
    return hashlib.sha1(msg.encode()).hexdigest()


def chkstr(
        token: str,
        username: str,
        hmd5: str,
        ac_id: str,
        ip: str,
        n: str,
        type_: str,
        info: str
    ) -> str:
    result = token + username
    result += token + hmd5
    result += token + ac_id
    result += token + ip
    result += token + n
    result += token + type_
    result += token + info
    return result


def trans_b64encode(s: str, alpha: Union[str, None] = None) -> str:
    """
    换表 base64

    这里用原生 base64 + 字符替换实现换表

    这种实现的代码量很少, 但是性能较低
    """
    import base64
    result = base64.b64encode(s.encode(encoding='latin-1')).decode()
    if not alpha:
        return result
    assert len(alpha) == 64, "base64字母表的长度必须为64"
    table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
    trans_table = str.maketrans(table, alpha)
    return result.translate(trans_table)


def info_(info: dict, token: str) -> str:
    import json, xxtea
    alpha = 'LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA'
    json_data = json.dumps(info).replace(' ', '')
    result = trans_b64encode(xxtea.encode(json_data, token), alpha)
    return f"{{SRBX1}}{result}"

