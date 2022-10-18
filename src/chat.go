package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const SERVER_LISTEN = 7778

func read_all_from(client net.Conn) (string, error) {
	var res string = ""

	for {
		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		res += string(buf)
		if n < 1024 {
			// if we have read everything we could, meaning there is nothing left
			break
		}
	}
	return res, nil
}

func send_and_get_usernames(client net.Conn, username string) string {
	_, err := client.Write([]byte(username))
	if err != nil {
		client.Close()
		log.Panic("Error while sending username: " + err.Error())
	}

	res, err := read_all_from(client)
	if err != nil {
		client.Close()
		log.Panic("Empty username from peer")
	}
	return res
}

func handle_client_logic(client net.Conn, username string) {
	defer client.Close()

	var peer_username string = send_and_get_usernames(client, username)
	fmt.Println("Speaking with " + peer_username + " at " + client.RemoteAddr().String())

	go func() {
		// async read from connection
		for {
			from_peer, err := read_all_from(client)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			fmt.Println(peer_username + "> " + from_peer)
		}
	}()

	// async write from connection
	scanner := bufio.NewScanner(os.Stdin)
	var to_peer string
	for {
		fmt.Print(username + "> ")
		scanner.Scan()
		to_peer = scanner.Text()
		if len(to_peer) == 0 {
			// cannot send empty messages
			continue
		}

		n, err := client.Write([]byte(to_peer))
		if err != nil || n == 0 {
			if err != nil {
				fmt.Printf("Write error: " + err.Error())
			}
			break
		}
	}
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
	client_mode_confirmed <- true
	client_string <- input
}

func main() {
	var username string
	fmt.Print("Username: ")
	fmt.Scan(&username)

	var server, err = net.Listen("tcp", fmt.Sprintf(":%d", SERVER_LISTEN))

	if err != nil {
		log.Panic(fmt.Sprintf("Could not open port %d to listen", SERVER_LISTEN))
	}
	log.Println(fmt.Sprintf("Listening on port %d", SERVER_LISTEN))

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
		handle_client_logic(<-client_accepted, username)
	case <-client_mode_confirmed:
		server.Close()
		handle_client_logic(connect_as_client(<-client_string), username)
	}
}
