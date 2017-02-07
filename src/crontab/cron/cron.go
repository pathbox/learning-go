package cron

import (
	"log"
	"runtime"
	"sort"
	"time"
)

type Cron struct {
	entries  []*Entry
	stop     chan struct{}
	add      chan *Entry
	snapshot chan []*Entry
	running  bool
	ErrorLog *log.Logger
	location *time.Location
}

type Job interface {
	Run()
}

type Schedule interface {
	Next(time.Time) time.Time
}
