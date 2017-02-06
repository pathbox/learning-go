type HelloHandler struct {
  db *sql.DB
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var name string

  row := h.db.QueryRow("SELECT myname FROM mytable")
  if err := row.Scan(&name); err != nil {
    http.Error(w, err.Error(), 500)
    return
  }
  fmt.Fprintf(w, "hi %s!\n", name)
}

func main() {
  db, err := sql.Open("postgres", "...")
  if err != nil {
    log.Fatal(err)
  }
  http.Handle("/hello", &HelloHandler{db: db})
  http.ListenAndServe(":9090", nil)
}

func TestHelloHandler_ServeHTTP(t *testing.T) {
    // Open our connection and setup our handler.
    db, _ := sql.Open("postgres", "...")
    defer db.Close()
    h := HelloHandler{db: db}
    // Execute our handler with a simple buffer.
    rec := httptest.NewRecorder()
    rec.Body = bytes.NewBuffer()
    h.ServeHTTP(rec, nil)
    if rec.Body.String() != "hi bob!\n" {
        t.Errorf("unexpected response: %s", rec.Body.String())
    }
}

package myapp

import (
  "database/sql"
)

type DB struct {
  *sql.DB
}

type Tx struct {
  *sql.Tx
}

// We then wrap the initialization function for our database and transaction:
// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
  db, err := sql.Open("postgres", dataSourceName)
  if err != nil {
    return nil, err
  }
  return &DB{db}, nil
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
    tx, err := db.DB.Begin()
    if err != nil {
        return nil, err
    }
    return &Tx{tx}, nil
}

// CreateUser creates a new user.
// Returns an error if user is invalid or the tx fails.
func (tx *Tx) CreateUser(u *User) error {
    // Validate the input.
    if u == nil {
        return errors.New("user required")
    } else if u.Name == "" {
        return errors.New("name required")
    }

    // Perform the actual insert and return any errors.
    return tx.Exec(`INSERT INTO users (...) VALUES`, ...)
}

tx, _ := db.Begin()
tx.CreateUser(&User{Name:"susy"})
tx.Commit()

tx, _ := db.Begin()
for _, u := range users {
    tx.CreateUser(u)
}
tx.Commit()