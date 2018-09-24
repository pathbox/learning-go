package main

import (
	"encoding/asn1"

	"fmt"

	"os"
)

func main() {

	mdata, err := asn1.Marshal(13)

	checkError(err)

	var n int

	_, err1 := asn1.Unmarshal(mdata, &n)

	checkError(err1)

	fmt.Println("After marshal/unmarshal: ", n)

}

func checkError(err error) {

	if err != nil {

		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())

		os.Exit(1)

	}

}
