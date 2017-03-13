func Shuffle(vals []int) []int {
  r := rand.New(rand.NewSource(time.Now().Unix()))
  ret := make([]int, len(vals))
  n := len(vals)
  for i := 0; i < n; i++ {
    randIndex := r.Intn(len(vals))
    ret[i] = vals[randIndex]
    vals = append(vals[:randIndex], vals[randIndex+1:]...)
  }
}

func main() {
  vals := []int{10, 12, 14, 16, 18, 20}
  r := rand.New(rand.NewSource(time.Now().Unix()))
  for _, i := range r.Perm(len(vals)) {
    val := vals[i]
    fmt.Println(val)
  }
}