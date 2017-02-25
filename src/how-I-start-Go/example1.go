package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := weatherQuery(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data) // 将data解析为json后赋值给w  response 返回为json结构
	})

	log.Fatal(http.ListenAndServe(":9090", nil))

}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println(string([]byte("Say Hello!")))
	w.Write([]byte("Hello!"))
}

func weatherQuery(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil { // 对resp.Body进行json解析，并且返回结果json给d d就是weatherData struct的变量
		return weatherData{}, err
	}

	return d, nil
}

type weatherData struct { // 嵌套json
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

// example:
// curl http://localhost:8080/weather/tokyo
{"name":"Tokyo","main":{"temp":295.9}}
