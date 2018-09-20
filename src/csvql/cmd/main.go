package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/campoy/csvql"
	sqle "gopkg.in/src-d/go-mysql-server.v0"
	"gopkg.in/src-d/go-mysql-server.v0/server"
	"gopkg.in/src-d/go-vitess.v0/mysql"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("could not find path: %v", err)
	}
	db, err := csvql.NewDatabase(path)
	if err != nil {
		log.Fatalf("could not create database: %v", err)
	}

	engine := sqle.NewDefault()
	engine.AddDatabase(db)

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     new(mysql.AuthServerNone),
	}
	server, err := server.NewDefaultServer(config, engine)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server on %s", config.Address)
	log.Fatal(server.Start())
}
