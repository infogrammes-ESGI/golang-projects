package main

import (
	"fmt"
	"net"
	"os"
)

func shell( /*host string*/ ) {

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("erreur de connexion a l'hote")
		os.Exit(1)
	}
	fmt.Fprintf(conn, "HelloWorld")
}

func main() {

	shell()

}
