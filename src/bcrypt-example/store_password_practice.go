package main 

import (
	"fmt"
	// "database/sql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// email := []byte("john.doe@gmail.com")
	password := []byte("1234567890")

	fmt.Println(bcrypt.DefaultCost)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	fmt.Println("hashed password: ", string(hashedPassword))

	// stmt, err := db.Prepare("INSERT INTO accounts SET hash=?, email=?")
	// if err != nil {
	// 	panic(err)
	// }

	// result, err := stmt.Exec(hashedPassword, email)
	// if err != nil {
	// 	panic(err)
	// }

	var expectedPassword string 
	expectedPassword = string(hashedPassword)

	if bcrypt.CompareHashAndPassword(password, []byte(expectedPassword)) != nil {
		fmt.Println("Password is not match")
	} else {
		fmt.Println("Password is right")
	}
}