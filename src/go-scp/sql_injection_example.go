customerId := r.URL.Query().Get("id")

query := "SELECT number, FROM creditcards WHERE customerId = " + customerId

row, _ := db.Query(query)

// customerId == 1 OR 1=1

// The query is:
// SELECT number FROM creditcards WHERE customerId = 1 OR 1=1
// you will dump all table records(yes, 1=1 will be true for any records)

// right way

customerId := r.URL.Query().Get("id")

query := "SELECT number, FROM creditcards WHERE customerId = ?"

smt, _ := db.Query(query, customerId)