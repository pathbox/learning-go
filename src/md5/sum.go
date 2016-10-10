package main

import(
  "crypto/md5"
  "fmt"
)

func main() {
  m := map[[16]byte]bool{}
  data := []byte("These pretzels are making me thirsty.")
  fmt.Printf("%x", md5.Sum(data))
  _, ok := m[md5.Sum(data)]
	fmt.Println(ok)
	m[md5.Sum(data)] = true
	fmt.Println(m)
	_, ok = m[md5.Sum(data)]
	fmt.Println(ok)
}
