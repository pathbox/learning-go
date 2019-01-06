package main

import (
	"encoding/json"
	"os"
	"time"
)

type logEntry struct {
	Timestamp     string
	Username      string
	RequestTypes  []string
	Error         string
	KeysOffered   []string
	GitHub        string
	ClientVersion string
}

func main() {
	le := &logEntry{Timestamp: time.Now().Format(time.RFC3339)}
	defer json.NewEncoder(os.Stdout).Encode(le)
}
