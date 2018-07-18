package main

import "database/sql"

func main() {
	cols, err := rows.Columns()
	if err != nil {
		// handle the error
	} else {
		dest := []interface{}{ // Standard MySQL columns
			new(uint64), // id
			new(string), // host
			new(string), // user
			new(string), // db
			new(string), // command
			new(uint32), // time
			new(string), // state
			new(string), // info
		}
		if len(cols) == 11 {
			// Percona Server
		} else if len(cols) > 8 {
			// Handle this case
		}
		err = rows.Scan(dest...)
		// Work with the values in dest
	}

	cols, err := rows.Columns() // Remember to check err afterwards
	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = new(sql.RawBytes)
	}
	for rows.Next() {
		err = rows.Scan(vals...)
		// Now you can check each element of vals for nil-ness,
		// and you can use type introspection and type assertions
		// to fetch the column into a typed variable.
	}
}
