func HasPrefix(s, prefix string) bool {
  return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
  return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func Contains(s, substr string) bool{
  for i := 0; i < len(s); i++{
    if HasPrefix(s[i:], substr){
      return true
    }
  }
  return false
}