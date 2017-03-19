package httpway

import (
	"net/http"
)

type Session interface {
	Id() string
	IsAuth() bool
	Username() string

	Get(key string) interface{}
	Set(key string, val interface{})
	GetInt(key string) int
	GetString(key string) string
}

type SessionManger interface {
	Get(w http.ResponseWriter, r *http.Request, log Logger) Session
	Set(w http.ResponseWriter, r *http.Request, session Session, log Logger)
}
