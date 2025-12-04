package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"srunClient/encryptlib"
	"strings"
	"time"
)

/*
rad_user_info: 获取用户信息的地址
get_challenge: 获取加密 token 的地址
srun_portal: 身份认证地址
*/
var targets = map[string]string{
	"rad_user_info": "http://net.szu.edu.cn/cgi-bin/rad_user_info",
	"get_challenge": "http://net.szu.edu.cn/cgi-bin/get_challenge",
	"srun_portal":   "http://net.szu.edu.cn/cgi-bin/srun_portal",
}

const (
	// banner
	Banner = `
   ________  __  __  ____                ________          __ 
  / __/_  / / / / / / __/_____ _____    / ___/ (_)__ ___  / /_
 _\ \  / /_/ /_/ / _\ \/ __/ // / _ \  / /__/ / / -_) _ \/ __/
/___/ /___/\____/ /___/_/  \_,_/_//_/  \___/_/_/\__/_//_/\__/ 
`
	// jsonp 标志
	callback string = "jQueryCallback"
	// 模拟 UA 头
	userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	// 加密常量
	TYPE string = "1"
	N    string = "200"
	ENC  string = "srun_bx1"
	ACID string = "12"
	//
	httpErrorMessage = `[!] 访问 %s 的过程中出现网络问题, 请检查:
1. 网络环境是否正常
2. 配置的 URL 是否正确
`
)

type getIpResp struct {
	Ip string `json:"online_ip"`
}

type getChallengeResp struct {
	Challenge string `json:"challenge"`
}

type srunPortalRes struct {
	Res string `json:"res"`
}

type srunPortalErr struct {
	ErrMsg string `json:"err_msg"`
}

func checkErrMsg(body []byte) (string, error) {
	var respJsonErr srunPortalErr
	fmt.Println("[*] 可能 srun 服务发生了异常, 正在检查报错码...")
	err := json.Unmarshal(body, &respJsonErr)
	if err != nil {
		fmt.Println("[!] 获取报错码失败")
		return "", err
	}
	fmt.Printf("[*] srun 报错码: %s", respJsonErr.ErrMsg)
	return respJsonErr.ErrMsg, nil
}

func checkCallback(body []byte, callback string) {
	fmt.Println("[*] 正在检查是否存在 callback 信息...")
	if string(body[:len(callback)]) == callback {
		fmt.Println("[!] 正常获取了 callback 信息, 请人工检查返回内容是否符合 json 格式:")
		fmt.Print(string(body[:100]))
		if len(body) < 100 {
			fmt.Println("")
		} else {
			fmt.Println("...")
		}
	} else {
		fmt.Println("[!] 没有 callback 信息, 请检查配置的 URL 是否正确")
	}
}

func getIp(callback, path string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()
	query.Add("callback", callback)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf(httpErrorMessage, path)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respJson getIpResp
	err = json.Unmarshal(body[len(callback)+1:len(body)-1], &respJson)
	if err != nil {
		fmt.Println("[!] 成功获取了返回的数据, 但是不存在 online_ip 值")
		errMsg, err := checkErrMsg(body[len(callback)+1 : len(body)-1])
		if err != nil {
			checkCallback(body, callback)
			return "", err
		}
		return "", fmt.Errorf("srun err: %s", errMsg)
	}
	fmt.Printf("[*] 成功获取登录 ip: %s\n", respJson.Ip)
	return respJson.Ip, nil
}

func getChallenge(callback, username, ip, path string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()
	query.Add("callback", callback)
	query.Add("username", username)
	query.Add("ip", ip)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf(httpErrorMessage, path)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respJson getChallengeResp
	err = json.Unmarshal(body[len(callback)+1:len(body)-1], &respJson)
	if err != nil {
		fmt.Println("[!] 成功获取了返回的数据, 但是不存在 challenge 值")
		errMsg, err := checkErrMsg(body[len(callback)+1 : len(body)-1])
		if err != nil {
			checkCallback(body, callback)
			return "", err
		}
		return "", fmt.Errorf("srun err: %s", errMsg)
	}
	fmt.Printf("[*] 成功获取加密 token: %s\n", respJson.Challenge)
	return respJson.Challenge, nil
}

func srunPortalLogin(callback, username, password, path, token, ip, os string) bool {
	fmt.Println("[*] 正在加密用户信息...")
	hmd5_password := encryptlib.Hmd5(password, token)
	info := encryptlib.GetInfo(encryptlib.Info{
		Username: username,
		Password: password,
		Ip:       ip,
		Acid:     ACID,
		EncVer:   ENC,
	}, token)
	chksum := encryptlib.Sha1(
		encryptlib.Chkstr(token, username, hmd5_password, ACID, ip, N, TYPE, info))
	fmt.Println("[*] 已完成用户信息加密, 准备进入身份认证")
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return false
	}
	current_time := fmt.Sprintf("%d000", time.Now().Unix())
	var loginQueryParams = map[string]string{
		"action":       "login",
		"callback":     callback,
		"username":     username,
		"password":     "{MD5}" + hmd5_password,
		"os":           os,
		"name":         os,
		"nas_ip":       "",
		"double_stack": "0",
		"chksum":       chksum,
		"info":         info,
		"ac_id":        ACID,
		"ip":           ip,
		"n":            N,
		"type":         TYPE,
		"captchaVal":   "",
		"_":            current_time,
	}
	query := req.URL.Query()
	for k, v := range loginQueryParams {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf(httpErrorMessage, path)
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respJsonRes srunPortalRes
	err = json.Unmarshal(body[len(callback)+1:len(body)-1], &respJsonRes)
	if err != nil || respJsonRes.Res != "ok" {
		fmt.Println("[!] 登陆失败")
		_, err = checkErrMsg(body[len(callback)+1 : len(body)-1])
		if err != nil {
			checkCallback(body, callback)
			return false
		}
		return false
	}
	fmt.Println("[*] 登录成功")
	return true
}

func srunPortalLogout(callback, username, ip, path string) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return
	}
	var logoutQueryParams = map[string]string{
		"action":   "logout",
		"callback": callback,
		"username": username,
		"ip":       ip,
	}
	query := req.URL.Query()
	for k, v := range logoutQueryParams {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("[*] 登出失败")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respJson srunPortalRes
	err = json.Unmarshal(body[len(callback)+1:len(body)-1], &respJson)
	if err != nil || respJson.Res != "ok" {
		fmt.Println("[*] 登出失败")
		return
	}
	fmt.Println("[*] 登出成功")
}

func login() {
	// 从环境变量读取学号与密码
	username := os.Getenv("username")
	if username == "" {
		// 兼容用户可能的输入拼写
		username = os.Getenv("userneme")
	}
	if username == "" {
		username = os.Getenv("SRUN_USERNAME")
	}
	passwordStr := os.Getenv("password")
	if passwordStr == "" {
		passwordStr = os.Getenv("SRUN_PASSWORD")
	}

	if username == "" || passwordStr == "" {
		fmt.Println("[!] 环境变量 username/SRUN_USERNAME 或 password/SRUN_PASSWORD 未设置")
		return
	}

	// 自动获取 ip
	ip, err := getIp(callback, targets["rad_user_info"])
	if err != nil {
		fmt.Println("获取 IP 失败, 请检查配置的 URL")
		return
	}

	// 设备类型默认为 Windows (避免与包名冲突, 使用 devOs 变量名)
	devOs := "Windows"

	token, err := getChallenge(callback, username, ip, targets["get_challenge"])
	if err != nil {
		return
	}
	ok := srunPortalLogin(callback, username, passwordStr, targets["srun_portal"], token, ip, devOs)
	if !ok {
		return
	}

	// start keepalive goroutine: 每隔30秒访问指定 URL, 若发生证书错误则重连
	go func() {
		for {
			time.Sleep(30 * time.Second)
			client := http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get("https://test.ustc.edu.cn/")
			if err != nil {
				// 仅在证书相关错误时触发重连
				if strings.Contains(err.Error(), "certificate") || strings.Contains(err.Error(), "x509") {
					fmt.Println("[*] 检测到证书错误, 尝试重新登录...")
					token2, err := getChallenge(callback, username, ip, targets["get_challenge"])
					if err != nil {
						fmt.Println("[!] 重新获取 token 失败:", err)
						continue
					}
					srunPortalLogin(callback, username, passwordStr, targets["srun_portal"], token2, ip, devOs)
				}
			} else {
				resp.Body.Close()
			}
		}
	}()

	fmt.Println("[*] 已启动防掉网守护进程 (每 30 秒检测)。按 Ctrl+C 退出。")
	select {}
}

func main() {
	fmt.Println(Banner)
	login()
}
