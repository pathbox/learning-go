// 退出程序

func main() {
  _, err := os.OpenFile("path/to/file")
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error Opening the file ")
    os.Exit(1)
  }
}