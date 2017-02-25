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

func temperature(city string, providers ...weatherProvider) (float64, error) {
    sum := 0.0

    for _, provider := range providers {
        k, err := provider.temperature(city)
        if err != nil {
            return 0, err
        }

        sum += k
    }

    return sum / float64(len(providers)), nil
}

