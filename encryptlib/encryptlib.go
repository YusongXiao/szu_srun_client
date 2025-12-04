package encryptlib

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Acid     string `json:"acid"`
	EncVer   string `json:"enc_ver"`
}

func Hmd5(msg, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha1(msg string) string {
	h := sha1.New()
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Chkstr(token, username, hmd5, ac_id, ip, n, type_, info string) string {
	result := token + username
	result += token + hmd5
	result += token + ac_id
	result += token + ip
	result += token + n
	result += token + type_
	result += token + info
	return result
}

func TransB64encode(s []byte, alpha *string) string {
	var encoder string
	if alpha == nil {
		encoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	} else {
		encoder = *alpha
	}
	b := base64.NewEncoding(encoder)
	return b.EncodeToString(s)
}

func GetInfo(info Info, token string) string {
	alpha := "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
	data, err := json.Marshal(info)
	if err != nil {
		return ""
	}
	result := TransB64encode(XxteaEncrypt(string(data), token), &alpha)
	return "{SRBX1}" + result
}
