package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"
  "time"
)

func main() {
  mw := multiWeatherProvider{
    openWeatherMap{},
    weatherUnderground{apiKey: "your-key-here"}
  }

  http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
    begin := time.Now()
    city := strings.SplitN(r.URL.Path, "/", 3)[2]

    temp, err := mw.temprature(city)  //  为什么mw可以调用temprature方法, 因为mw中嵌套了 openWeatherMap{}, 可以简单的理解为继承或mix-in
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(map[string]interface{}{
      "city": city,
      "temp": temp,
      "took": time.Since(begin).String(),
      })
    })

  log.Fatal(http.ListenAndServe(":9090", nil))
}

type multiWeatherProvider []weatherProvider

type weatherUnderground struct {
  apiKey string
}

type weatherProvider interface {
    temperature(city string) (float64, error) // in Kelvin, naturally
}

type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
    if err != nil {
        return 0, err
    }

    defer resp.Body.Close()

    var d struct {
        Main struct {
            Kelvin float64 `json:"temp"`
        } `json:"main"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
        return 0, err
    }

    log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
    return d.Main.Kelvin, nil
}

func (w multiWeatherProvider) temperature(city string) (float64, error){
  // Make a channel for temperatures, and a channel for errors.
    // Each provider will push a value into only one.

  temps := make(chan float64, len(w))
  errs := make(chan error, len(w))

  // For each provider, spawn a goroutine with an anonymous function.
    // That function will invoke the temperature method, and forward the response.
  for _, provider := range w {
    go func(p weatherProvider) {  // 并发的进行调用第三方接口的工作， 将得到的值存入channel
      k, err := p.temperature(city)
      if err != nil {
        errs <- err
        return
      }
      temps <- k
    }(provider)
  }

  sum := 0.0

  for i := 0; i < len(w); i++ {  // 再在循环,使用select 监听从temps channel 中取出 前面并发goroutinue中得到的存入temps channel的temp值
    select{
    case temp :=<-temps:
      sum += temp
    case err := <-errrs:
      return 0, err
    }
  }

  return sum / float64(len(w)), nil
}



























