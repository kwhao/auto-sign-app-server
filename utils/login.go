package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

// BaseString 加密用的字符串
const BaseString = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"

// LoginUrl 登录地址
const LoginUrl = "https://nbu.campusphere.net/portal/login"

// @title 指定长度随机字符串
func randStr(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = BaseString[rand.Intn(len(BaseString))]
	}
	return string(b)
}

// pkcs5Padding 用于AES padding 填充
func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// aesCbcEncrypt AES-128-CBC-PKCS5Padding 加密主体
func aesCbcEncrypt(plainText []byte, key []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = pkcs5Padding(plainText, 16)
	//指定初始向量vi
	//指定分组模式，返回一个BlockMode接口对象
	iv := []byte(randStr(16))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

// headerHelp 用于构造请求头
func headerHelp(request *http.Request) {
	pHeaders := map[string]string{
		"Host":                      "uis.nbu.edu.cn",
		"Connection":                "close",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"Origin":                    "https://uis.nbu.edu.cn",
		"Content-Type":              "application/x-www-form-urlencoded",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Apple Mac OS X 13_1_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.88 Safari/605.1.15",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"X-Forwarded-For":           "10.22.99.98,10.22.97.98,114.114.114.114",
		"X-Real-IP":                 "10.22.99.98,10.22.97.98,114.114.114.114",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-User":            "?1",
		"Sec-Fetch-Dest":            "document",
		"Referer":                   "https://uis.nbu.edu.cn/authserver/login?service=http%3A%2F%2Fehall.nbu.edu.cn%2Flogin%3Fservice%3Dhttps%3A%2F%2Fehall.nbu.edu.cn%2Fnew%2Findex.html",
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
	}
	for hKey, hValue := range pHeaders {
		request.Header.Set(hKey, hValue)
	}
}

// logRedir 用于处理重定向
func logRedir(loginReq *http.Request, logRes *http.Response, method string) (*http.Request, *http.Response) {
	finalUrl := logRes.Header.Get("Location")
	//fmt.Println(finalUrl)
	loFinalReq, _ := http.NewRequest(method, finalUrl, nil)
	headerHelp(loFinalReq)
	for _, cookie := range loginReq.Cookies() {
		loFinalReq.AddCookie(cookie)
	}
	for _, cookie := range logRes.Cookies() {
		loFinalReq.AddCookie(cookie)
	}
	client := &http.Transport{}
	FinalRes, _ := client.RoundTrip(loFinalReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(FinalRes.Body)

	return loFinalReq, FinalRes
}

// ValidateAccount 判断账号登录是否需要验证码
func ValidateAccount(username string) bool {
	validUrl := "https://uis.nbu.edu.cn/authserver/needCaptcha.html?username=" + username + "&pwdEncrypt2=pwdEncryptSalt"
	getReq, _ := http.NewRequest("GET", validUrl, nil)
	getReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57")
	client0 := &http.Client{}
	res, err := client0.Do(getReq)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)
	result, _ := ioutil.ReadAll(res.Body)
	if find := strings.Contains(string(result), "false"); find {
		return true
	}
	return false
}

// Login 登录获取 ehall 首页 cookie
func Login(username string, password string, loginUrl string) ([]*http.Cookie, error) {
	getReq, err := http.NewRequest("GET", loginUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("err %s", "加载登录信息发生错误")
		//panic("Something GET error")
	}
	getReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57")
	client0 := &http.Client{}
	res, err := client0.Do(getReq)
	if err != nil {
		return nil, fmt.Errorf("err %s", "加载登录信息发生错误")
		//panic("Something GET error")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	if res.StatusCode != 200 {
		panic("res status error")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		//panic(err)
		return nil, fmt.Errorf("err %s", "解析登录页面发生错误")
	}
	cryKey, _ := doc.Find("input#pwdDefaultEncryptSalt").Attr("value")
	lt, _ := doc.Find("input[name=lt]").Attr("value")
	dllt, _ := doc.Find("input[name=dllt]").Attr("value")
	execution, _ := doc.Find("input[name=execution]").Attr("value")
	rmShown, _ := doc.Find("input[name=rmShown]").Attr("value")
	_eventId, _ := doc.Find("input[name=_eventId]").Attr("value")

	rndPwd := randStr(64) + password
	enPwd := base64.StdEncoding.EncodeToString(aesCbcEncrypt([]byte(rndPwd), []byte(cryKey)))
	enPwd = url.QueryEscape(enPwd)

	strTmp := "captchaResponse=&" + "username=" + username + "&password=" + enPwd + "&lt=" + lt + "&dllt=" + dllt + "&execution=" + execution + "&rmShown=" + rmShown + "&_eventId=" + _eventId
	bodyP := strings.NewReader(strTmp)
	req, _ := http.NewRequest("POST", loginUrl, bodyP)
	for _, cookie := range res.Cookies() {
		req.AddCookie(cookie)
	}
	headerHelp(req)
	//第一次发送带用户名和密码的payload进行验证
	client := &http.Transport{}
	response, err := client.RoundTrip(req)
	if err != nil {
		//panic(err)
		return nil, fmt.Errorf("err %s", "提交数据可能有误")
	}
	cont, _ := goquery.NewDocumentFromReader(response.Body)
	//用户名密码验证成功后的响应状态为302
	if response.StatusCode != 302 {
		msg := cont.Find("span#msg").Text()
		fmt.Println(msg)
		return nil, fmt.Errorf("err %s", "提交数据可能有误")
		//panic("提交数据可能有误")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)

	//获取验证成功后的重定向地址，拼接不同阶段的cookie进行重定向
	//可能发生多次302/301重定向，且需要多次重定向拼接
	locationUrl := response.Header.Get("Location")
	checkReq, _ := http.NewRequest("GET", locationUrl, nil)
	headerHelp(checkReq)
	//拼接上次请求的payload中的cookie
	for _, cookie := range req.Cookies() {
		checkReq.AddCookie(cookie)
	}
	//拼接上次请求后的response中的set-cookie
	for _, cookie := range response.Cookies() {
		checkReq.AddCookie(cookie)
	}
	client1 := &http.Transport{}
	checkRes, err := client1.RoundTrip(checkReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(checkRes.Body)

	fReq := checkReq
	fRes := checkRes
	for {
		redirReq, redirRes := logRedir(fReq, fRes, "GET")
		if redirRes.StatusCode != 302 {
			if redirRes.StatusCode != 301 {
				fReq = redirReq
				break
			}
		}
		fReq = redirReq
		fRes = redirRes
	}
	return fReq.Cookies(), nil
}
