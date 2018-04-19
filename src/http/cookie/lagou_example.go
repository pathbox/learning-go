package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

var cookies_lagou []*http.Cookie

const (
	login_url_lagou string = "https://passport.lagou.com/login/login.html"

	post_login_info_url_lagou string = "https://passport.lagou.com/login/login.json"

	username_lagou string = " "
	password_lagou string = " "
)

func getToken(contents io.Reader) (string, string) {

	data, _ := ioutil.ReadAll(contents)
	regCode := regexp.MustCompile(`X_Anti_Forge_Code\s+\=(.+?);`)
	if regCode == nil {
		log.Fatal("解析Code出错...")
	}

	//提取关键信息
	code := regCode.FindAllStringSubmatch(string(data), -1)[0][1]

	regToken := regexp.MustCompile(`X_Anti_Forge_Token\s+\=(.+?);`)
	if regToken == nil {
		fmt.Println("MustCompile err")
	}

	//提取关键信息
	token := regToken.FindAllStringSubmatch(string(data), -1)[0][1]

	return token, code
}

func login_lagou() {
	//获取登陆界面的cookie
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	req, _ := http.NewRequest("GET", login_url_lagou, nil)
	res, _ := client.Do(req)
	for k, v := range res.Cookies() {
		fmt.Printf("%v=%v\n", k, v)
	}
	token, code := getToken(res.Body)
	//post数据
	postValues := url.Values{}
	postValues.Add("isValidate", "true")
	postValues.Add("username", username_lagou)
	postValues.Add("password", password_lagou)
	postValues.Add("request_form_verifyCode", "")
	postValues.Add("submit", "")
	body := ioutil.NopCloser(strings.NewReader(postValues.Encode())) //把form数据编下码
	requ, _ := http.NewRequest("POST", post_login_info_url_lagou, body)

	requ.Header.Set("X-Requested-With", "XMLHttpRequest")
	requ.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	requ.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	// requ.Header.Set("Host", "passport.lagou.com")
	// requ.Header.Set("Origin", "https://passport.lagou.com")
	requ.Header.Add("X-Anit-Forge-Token", token)
	requ.Header.Add("X-Anit-Forge-Code", code)
	// requ.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	// requ.Header.Set("Connection", "keep-alive")
	// requ.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	// requ.Header.Set("Accept-Encoding", "gzip, deflate, br")
	requ.Header.Set("Referer", "https://passport.lagou.com/login/login.html")
	// for _, v := range res.Cookies() {
	// 	requ.AddCookie(v)
	// }

	res, _ = client.Do(requ)
	//cookies_lagou = res.Cookies()
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(string(data))
}

func main() {
	login_lagou()
}
