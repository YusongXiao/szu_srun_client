package encryptlib

func h(s string, i int) int64 {
	if i < len(s) {
		return int64(s[i])
	}
	return 0
}

func s(a string, b bool) []int64 {
	v := make([]int64, 0)
	for i := 0; i < len(a); i += 4 {
		n := h(a, i) | (h(a, i+1) << 8) | (h(a, i+2) << 16) | (h(a, i+3) << 24)
		v = append(v, n)
	}
	if b {
		v = append(v, int64(len(a)))
	}
	return v
}

func XxteaEncrypt(msg, key string) []byte {
	var m, e, p, d int64
	if msg == "" {
		return []byte{}
	}
	v := s(msg, true)
	k := s(key, false)
	if len(k) < 4 {
		k = append(k, make([]int64, 4-len(k))...)
	}
	n := len(v) - 1
	z := v[n]
	var c int64 = 0x9e3779b9
	q := 6 + 52/(n+1)
	for 0 < q {
		d = d + c&0xffffffff
		e = d >> 2 & 3
		p = 0
		for p < int64(n) {
			y := v[p+1]
			m = (z>>5 ^ y<<2) + (y>>3 ^ z<<4 ^ (d ^ y)) + (k[p&3^e] ^ z)
			v[p] = (v[p] + m) & 0xffffffff
			z = v[p]
			p += 1
		}
		y := v[0]
		m = (z>>5 ^ y<<2) + (y>>3 ^ z<<4 ^ (d ^ y)) + (k[p&3^e] ^ z)
		v[n] = (v[n] + m) & 0xffffffff
		z = v[n]
		q -= 1
	}
	result := make([]byte, 0)
	for _, i := range v {
		result = append(result, byte(i&0xff), byte(i>>8&0xff), byte(i>>16&0xff), byte(i>>24&0xff))
	}
	return result
}
