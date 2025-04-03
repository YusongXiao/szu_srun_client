class String(str):
    """
    实现 js 的 charCodeAt, formCharCode
    
    这样代码结构和 js 的实现长得会比较像, 方便我调试
    """
    def charCodeAt(self, i: int) -> int:
        if len(self) > i:
            return ord(self[i])
        return 0
    
    @classmethod
    def fromCharCode(cls, *charCodes: list[int]) -> str:
        result = ''
        for c in charCodes:
            result += chr(c)
        return result


def s(a: String, b: bool) -> list[int]:
    c = len(a)
    v = []

    for i in range(0, c, 4):
        v.append(a.charCodeAt(i)
                 | a.charCodeAt(i + 1) << 8
                 | a.charCodeAt(i + 2) << 16
                 | a.charCodeAt(i + 3) << 24)
    
    if b: v.append(c)
    return v


def l(a: list[int], b: bool):
    d = len(a)
    c = (d - 1) << 2

    if b:
        m = a[d - 1]
        if m < c - 3 or m > c:
            return None
        c = m
    
    for i in range(0, d):
        a[i] = String.fromCharCode(
            a[i] & 0xff,
            a[i] >> 8 & 0xff,
            a[i] >> 16 & 0xff,
            a[i] >> 24 & 0xff
        )
    return "".join(a)[0:c] if b else "".join(a)


def encode(str_: str, key: str) -> str:
    """xxtea 加密"""
    import math
    str_, key = String(str_), String(key)

    if str_ == "": return ""
    v = s(str_, True)
    k = s(key, False)
    if len(k) < 4:
        k = k + [0] * (4 - len(k))
    n = len(v) - 1
    z = v[n]
    y = v[0]
    c = 0x86014019 | 0x183639A0
    m = 0
    e = 0
    p = 0
    q = math.floor(6 + 52 / (n + 1))
    d = 0

    while 0 < q:
        d = d + c & (0x8CE0D9BF | 0x731F2640)
        e = d >> 2 & 3
        p = 0
        while p < n:
            y = v[p + 1]
            m = z >> 5 ^ y << 2
            m += y >> 3 ^ z << 4 ^ (d ^ y)
            m += k[p & 3 ^ e] ^ z
            z = v[p] = v[p] + m & (0xEFB8D130 | 0x10472ECF)
            p += 1
        y = v[0]
        m = z >> 5 ^ y << 2
        m += y >> 3 ^ z << 4 ^ (d ^ y)
        m += k[p & 3 ^ e] ^ z
        z = v[n] = v[n] + m & (0xBB390742 | 0x44C6F8BD)
        q -= 1
    
    return l(v, False)
