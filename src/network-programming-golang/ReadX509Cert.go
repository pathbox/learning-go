package main

import (
	"crypto/x509"

	"fmt"

	"os"
)

func main() {

	certCerFile, err := os.Open("cary.pathobx.cer")

	checkError(err)

	derBytes := make([]byte, 1000) // bigger than the file

	count, err := certCerFile.Read(derBytes)

	checkError(err)

	certCerFile.Close()

	// trim the bytes to actual length in call

	cert, err := x509.ParseCertificate(derBytes[0:count])

	checkError(err)

	fmt.Printf("Name %s\n", cert.Subject.CommonName)

	fmt.Printf("Not before %s\n", cert.NotBefore.String())

	fmt.Printf("Not after %s\n", cert.NotAfter.String())
	fmt.Println("Key content: ", cert.PublicKey)

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
