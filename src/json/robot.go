package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type Result struct {
	Reply       string `json:"_reply"`
	SessionID   string `json:"_session_id"`
	FaqQuestion string `json:"_faq_question"`
	Intent      string `json:"_intent"`
	Action      string `json:"_action"`
	FaqAnswer   string `json:"_faq_answer"`
}

type Reply struct {
	Items []Item `json:"items"`
}

type Item struct {
	VoiceType    string `json:"voice_type"`
	VoiceContent string `json:"voice_content"`
}

func main() {

	pa := make(map[string]interface{})
	data := make(map[string]string)
	pa["timestamp"] = "1528252981"
	data["query_text"] = "天气"
	data["udesk_call_id"] = ""
	pa["data"] = data

	url := "url"

	res, err := PostJson(url, pa)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	rs, err := ioutil.ReadAll(res.Body)

	m := make(map[string]interface{})

	e := json.Unmarshal(rs, &m)

	if e != nil {
		fmt.Println(e)
	}

	fmt.Println(m["data"])

	mData := m["data"]

	mb, e := json.Marshal(mData)
	if e != nil {
		fmt.Println(e)
	}

	result := &Result{}

	err = json.Unmarshal(mb, result)
	if err != nil {
		fmt.Println("===", err)
	}

	log.Println(result)

	rp := result.Reply

	fmt.Println("rprprp:", rp)

	rpb := []byte(rp) // Nice Way 将 json字符串强转为 byte json string => byte，之后就能进行json.Unmarshal  而不是通过json.Marshal 转为byte

	reply := &Reply{}

	err = json.Unmarshal(rpb, &reply)

	if err != nil {
		fmt.Println("++++", err)
	}

	for _, item := range reply.Items {
		fmt.Printf("voice_type:%s,voice_content:%s\n", item.VoiceType, item.VoiceContent)
	}
}

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	return client
}

func PostJson(url string, paramMap map[string]interface{}) (*http.Response, error) {
	reqBody, _ := json.Marshal(paramMap)
	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Println("PostJson", err)
		return nil, err
	}

	client := NewHTTPClient()
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Println("PostJson", err)
	}

	return resp, err
}

/*
json 字符串：{\"items\":[{\"voice_type\":\"tts\", \"voice_content\":\"哪天的天气？\"},{\"voice_type\":\"tts\", \"voice_content\":\"have a nice day\"}]}


map[_faq_answer: _reply:{"items":[{"voice_type":"tts", "voice_content":"哪天的天气？"}]} _session_id: _faq_question: _intent:sys_weather _action:]
2018/06/06 16:02:39 &{{"items":[{"voice_type":"tts", "voice_content":"哪天的天气？"}]}   sys_weather  }

json 字符串 => rprprp: {"items":[{"voice_type":"tts", "voice_content":"哪天的天气？"},{"voice_type":"tts", "voice_content":"have a nice day"}]}

voice_type:tts,voice_content:哪天的天气？
voice_type:tts,voice_content:have a nice day
*/
