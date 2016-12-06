package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://akagi201:@localhost:5432/bookstore?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandlerFunc("/books", booksIndex)
	http.HandlerFunc("/books/show", booksShow)
	http.HandlerFunc("/books/create", booksCreate)
	http.ListenAndServe(":3000", nil)
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, http.StatusText(405), 405)
    return
  }
  rows, err := db.Query("SELECT * FROM books")
  if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }
  defer rows.Close()

  bks := make([]*Book, 0)
  for rows.Next() {
    bk := new(Book)
    err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }
    bks = append(bks, bk)
  }
  if err = rows.Err(); err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }
  for _, bk := range bks {
    fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
  }
}

func booksShow(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, http.StatusText(405), 405)
    return
  }

  isbn := r.FormValue("isbn")
  if isbn == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }

  row := db.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)

  bk := new(Book)
  err := row.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
  if err == sql.ErrNoRows {
    http.NotFound(w, r)
    return
  } else if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }

  fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
}

func booksCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Error(w, http.StatusText(405), 405)
    return
  }
  isbn := r.FormValue("isbn")
  title := r.FormValue("title")
  author := r.FormValue("author")
  if isbn == "" || title == "" || author == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
  price, err := strconv.ParseFloat(r.FormValue("price"), 32)
  if err != nil {
    http.Error(w, http.StatusText(400), 400)
    return
  }

  result, err := db.Exec("INSERT INTO books VALUES($1, $2, $3, $4)", isbn, title, author, price)
  if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }
  rowsAffected, err := result.RowsAffected()
  if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }

  fmt.Fprintf(w, "Book %s created successfully (%d row affected)\n", isbn, rowsAffected)
}
