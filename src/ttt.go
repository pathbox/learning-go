package main

import (
	"encoding/json"
	"fmt"

	"log"
)

func main() {
	fmt.Println("Hello, playground")

	replyByte := []byte("{\"items\":[{\"voice_type\":\"tts\", \"voice_content\":\"你想知道哪天的天气？\"},{\"voice_type\":\"tts\", \"voice_content\":\"我是第二个语音\"}], \"vars\":[{\"key\":\"weather_city\", \"type\":\"string\",\"value\":\"\"},{\"key\":\"weather_date\", \"type\":\"string\",\"value\":\"\"}]}")

	reply := &Reply{}
	err := json.Unmarshal(replyByte, reply)
	if err != nil {
		log.Println("========", err)
	}

	fmt.Println("reply: ", reply)
}

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
	Vars  []Var  `json:"vars"`
}

type Var struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Item struct {
	VoiceType    string `json:"voice_type"`
	VoiceContent string `json:"voice_content"`
}
