package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"udesk/model"
	"udesk/resource/database"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(BasicInfo(r, "Success"))
	}))
}

type ViewFunc func(http.ResponseWriter, *http.Request)

func BasicAuth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := r.URL.Query().Get("timestamp")
		t := time.Now()
		tm, _ := time.Parse("20060102150405", timestamp)
		x := tm.Sub(t).Minutes()
		if x > 490 || x < 470 {
			fmt.Println("expired!")
			log.Panicln("++++++ 10 minutes expired! ++++++")
		}

		username, passwd := GetUsernameAndPasswd(r) //获得 http auth 参数
		if username == "" || passwd == "" {
			log.Panicln("username or passwd is nil")
		}
		if CheckSign(username, passwd, timestamp) {
			f(w, r)
			return
		}
		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		log.Println(BasicInfo(r, "Fail"))
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func CheckErr(err error, res http.ResponseWriter) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func BasicInfo(r *http.Request, info string) string {
	url := r.URL.String()
	info = "Method: " + r.Method + " Url: " + url + " ++++++Basic HTTP Auth " + info + "++++++"
	return info
}

func AuthToken(appid string) string {
	db, err := database.DB()
	if err != nil {
		log.Println("db establish fail")
		return ""
	}
	defer db.Close()
	account := model.Account{}
	db.Joins("JOIN apps ON apps.account_id = accounts.id").Where("apps.sid = ?", appid).Find(&account)
	if account.ID == 0 {
		log.Println("account not founded")
		return ""
	} else {
		return account.AuthToken
	}
}

func GetUsernameAndPasswd(r *http.Request) (username, passwd string) {
	basicAuthPrefix := "Basic "
	auth := r.Header.Get("Authorization") //获取 request header
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):]) // 解码认证信息
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			return string(pair[0]), string(pair[1])
		} else {
			return "", ""
		}
	}
	return "", ""
}

func CheckSign(username, passwd, timestamp string) bool {
	authToken := AuthToken(username)
	if authToken == "" {
		return false
	}
	sign_str := username + authToken + timestamp
	sha1Newer := sha1.New()
	sha1Newer.Write([]byte(sign_str))
	signResult := sha1Newer.Sum(nil)
	sign := fmt.Sprintf("%x", signResult)
	fmt.Println(sign)
	if sign == passwd {
		return true
	} else {
		return false
	}
}
