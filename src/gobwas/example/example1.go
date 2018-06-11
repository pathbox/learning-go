package main

import (
	"io"
	"net/http"
)

func main() {
	http.ListenAndServe(":9090", http.HandlerFunc(index))
}
func index(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {

	}

	go func() {
		defer conn.Close()
		var (
			state  = ws.StateServerSide
			reader = wsutil.NewReader(conn, state)
			writer = wsutil.NewWriter(conn, state, ws.OpText)
		)

		for {
			header, err := reader.NextFrame()
			if err != nil {

			}

			writer.Reset(conn, state, header.OpCode)
			if _, err = io.Copy(writer, reader); err != nil {

			}
			if err = writer.Flush(); err != nil {

			}
		}
	}()
}
