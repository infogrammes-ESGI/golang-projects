package main

import (
	"fmt"
	"log"
	"net"
)

func handle_client_logic(client net.Conn) {
	defer client.Close()
}

func connect_as_client(client_string string) net.Conn {
	conn, err := net.Dial("tcp", client_string) // appell TCP sur l'ip:port
	if err != nil {                             // verification des erreurs
		if conn != nil {
			conn.Close() // ferme la connexion si err == nil et que la connexion est initialisee, permet de fermer proprement
		}
		log.Panic("Could not connect to " + client_string)
	}
	return conn
}

func wait_in_server_mode(server_mode_confirmed chan bool, client_mode_confirmed chan bool, client_accepted chan net.Conn, server net.Listener) {
	for {
		conn, err := server.Accept()
		if err != nil {
			// if the server throw an error, it might be because of the connection being closed because the user
			// wants to go in client mode, so we need to check if client_mode_confirmed has been set to true
			// or if it is just a normal socket error
			switch {
			case <-client_mode_confirmed:
				return
			default:
				log.Println("Could not accept from "+conn.RemoteAddr().String(), err)
				continue
			}
		}
		server_mode_confirmed <- true
		client_accepted <- conn
	}
}

func wait_in_client_mode(client_mode_confirmed chan bool, client_string chan string) {
	var input string
	fmt.Print("Who to connect to (host:port) ? ")
	fmt.Scan(&input)
	client_string <- input
	client_mode_confirmed <- true
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

	// for the server
	server_mode_confirmed := make(chan bool)
	// used when in server mode to send to the main part the client's connection object
	client_accepted := make(chan net.Conn)

	// for the client
	client_mode_confirmed := make(chan bool)
	// used when in client mode to send to the main part the client's host:port
	client_string := make(chan string)

	go wait_in_server_mode(server_mode_confirmed, client_mode_confirmed, client_accepted, server)
	go wait_in_client_mode(client_mode_confirmed, client_string)

	select {
	case <-server_mode_confirmed:
		handle_client_logic(<-client_accepted)
	case <-client_mode_confirmed:
		server.Close()
		handle_client_logic(connect_as_client(<-client_string))
	}
}
