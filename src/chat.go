package main

import (
	"fmt"
	"log"
	"net"
)

func server_mode(server_mode_confirmed chan bool, client_mode_confirmed chan bool, client_accepted chan net.Conn, server net.Listener) {
	for {
		log.Println("Waiting for a connection")
		conn, err := server.Accept()
		if <-client_mode_confirmed {
			return
		}
		if err != nil {
			log.Println("Could not accept from "+conn.RemoteAddr().String(), err)
			continue
		}
		log.Println("Accepted from " + conn.RemoteAddr().String())
		server_mode_confirmed <- true
		client_accepted <- conn
	}
}

func handle_client_as_server(client net.Conn) {
}

func main() {
	var username string
	fmt.Print("Username: ")
	fmt.Scan(&username)

	var server, err = net.Listen("tcp", ":7777")

	if err != nil {
		log.Panic("Could not open port 7777 to listen")
	}
	log.Println("Listening on port 7777")

	server_mode_confirmed := make(chan bool)
	client_mode_confirmed := make(chan bool)

	// used when in server mode to send to the main part the client object
	client_accepted := make(chan net.Conn)

	go server_mode(server_mode_confirmed, client_mode_confirmed, client_accepted, server)
	go func() {
		var input string
		fmt.Println("What to do ?")
		fmt.Scan(&input)
		client_mode_confirmed <- true
	}()

	select {
	case <-server_mode_confirmed:
		log.Println("Got server mode")
		handle_client_as_server(<-client_accepted)
	case <-client_mode_confirmed:
		log.Println("Got client mode")
		server.Close()
	}
}
