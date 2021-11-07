package main
 
import (
    "sync"
    "net"
    "strconv"
    "fmt"
    "log"
 
)
 
const (
    MAX_CONCURRENCY = 10000 
    CHANNEL_CACHE = 200
)
 
var tmpChan = make(chan struct{}, MAX_CONCURRENCY)
var waitGroup sync.WaitGroup
 
func main(){
    concurrency()
    waitGroup.Wait()
}
 
//进行网络io
func request(currentCount int){
    fmt.Println("request" + strconv.Itoa(currentCount) + "\r")
    tmpChan <- struct{}{}
    conn, err := net.Dial("tcp",":8080")
    <- tmpChan
    if err != nil { log.Fatal(err) }
    defer conn.Close()
    defer waitGroup.Done()
}
 
//并发
func concurrency(){
    for i := 0;i < MAX_CONCURRENCY;i++ {
        waitGroup.Add(1)
        go request(i)
    }
}