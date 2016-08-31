package main

import "github.com/nareix/curl"
import "time"
import "log"
import "net/http"
import "fmt"

func main() {


req := curl.New("https://kernel.org/pub/linux/kernel/v4.x/linux-4.0.4.tar.xz")

req.Method("POST") // can be  PUT POST DELETE.....

req.Header("MyHeader", "value")  // custom header

req.Headers = http.Header{
  "User-Agent": {"mycurl/1.0"},
}

ctrl := req.ControlDownload()
go func(){
  // control functions are thread safe
    ctrl.Stop()   // Stop download
    ctrl.Pause()  // Pause download
    ctrl.Resume() // Resume download
}()

req.DialTimeout(time.Second * 10) // TCP connection Timeout
req.Timeout(time.Second * 30)

req.Progress(func (p curl.ProgressStatus){
  log.Println(
    "Stat", p.Stat,
    "speed", curl.PrettySpeedString(p.Speed),
    "len", curl.PrettySizeString(p.ContentLength),
    "got", curl.PrettySizeString(p.Size),
    "percent", p.Percent,
    "paused", p.Paused,
  )
}, time.Second)

res, _ := req.Do()
fmt.Println(res.HttpResponse)                             // related *http.Response struct
log.Println(res.Body)                        // Body in string
log.Println(res.StatusCode)                  // HTTP Status Code: 200,404,302 etc
// log.Println(res.Hearders)                    // Reponse headers
log.Println(res.DownloadStatus.AverageSpeed) // Average speed
}
