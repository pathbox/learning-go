package main  
  
import (  
    "fmt"  
    "math/rand"  
    "time"  
)  
  
func main() {  
    rand.Seed(time.Now().UnixNano())  
    for i := 0; i < 10; i++ {  
        x := rand.Intn(100)  
        fmt.Println(x)  
    }  
}  