package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	json "github.com/json-iterator/go"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var MyDB *gorm.DB

type RespResult struct {
	ErrCode string      `json:"ErrCode"`
	ErrDesc string      `json:"ErrDesc"`
	Result  interface{} `json:"Result"`
}

type Base struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
	// UpdatedAt time.Time `gorm:"type:timestamp;not null;default:'1971-01-01 01:01:01'"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type IvrVoice struct {
	Base
	AppID     string `gorm:"size:191;unique_index;not null"` // app sid
	DefVoices string `gorm:"type:text"`
}

func init() {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"",
		"127.0.0.1",
		3306,
		"udesk_cti_cloud",
	)

	MyDB, err = gorm.Open("mysql", dsn)
	// 连接池
	if err == nil {
		MyDB.DB().SetMaxIdleConns(10)
		MyDB.DB().SetMaxOpenConns(100)
		MyDB.DB().Ping()
		MyDB.LogMode(false)
	} else {
		log.Panic("Gorm Open Error", err)
	}
}

func main() {
	addr := ":9090"
	router := mux.NewRouter()
	router.HandleFunc("/files", filesHandler).Methods("GET")
	router.HandleFunc("/files/{id:[0-9]+}", filesGetHandler).Methods("GET")
	router.HandleFunc("/files", fileCreate).Methods("POST")

	log.Fatal(http.ListenAndServe(addr, router))
}

func filesHandler(w http.ResponseWriter, r *http.Request) {

}

type File struct {
	AppID     string      `json:"AppID"`
	DefVoices interface{} `json:"DefVoices"`
}

func fileCreate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	rs := RespResult{"OK", "", nil}
	defer r.Body.Close()

	var p File

	log.Println(string(body))

	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Println("Unmarshal", err)
		rs = RespResult{"err", err.Error(), nil}
		ResponseJson(w, rs)
		return
	}

	voice, _ := json.Marshal(p.DefVoices)
	log.Println("Voice", string(voice))

	f := &IvrVoice{AppID: p.AppID, DefVoices: string(voice)}
	err = MyDB.Create(f).Error
	if err != nil {
		log.Println("Create", err)
		rs = RespResult{"err", err.Error(), nil}
		ResponseJson(w, rs)
		return
	}
	log.Println("+++")
	rs.Result = f
	ResponseJson(w, rs)

}

type VoiceInfo struct {
	AppID     string                       `json:"AppId"`
	DefVoices map[string]map[string]string `json:"DefVoices"`
}

func filesGetHandler(w http.ResponseWriter, r *http.Request) {
	appID := GetURLAppID(r)

	v := &IvrVoice{}

	err := MyDB.Where("app_id = ?", appID).First(v).Error

	if err != nil {
		log.Println(err)
	}
	m := make(map[string]map[string]string)
	ms := "{\"Welcome\":{\"Type\":\"wav\", \"Content\":\"/content/hello.wav\"}}"
	e := json.Unmarshal([]byte(ms), &m) // don't forget & 用取地址符
	log.Println(e)

	vf := VoiceInfo{AppID: v.AppID, DefVoices: m}

	rs := RespResult{"OK", "OK", vf}
	ResponseJson(w, rs)
}

func GetURLAppID(r *http.Request) string {
	return r.URL.Query().Get("AppId")
}

func ResponseJson(w http.ResponseWriter, rs interface{}) {
	rsJson, _ := json.Marshal(rs)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	log.Println("Response: ", string(rsJson))
	w.Write(rsJson)
}
