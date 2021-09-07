package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// ssh-keygen -P "" -f id_rsa

var config = &ssh.ServerConfig{
	PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		log.Printf("%s - %s:%s (%s)\n", c.RemoteAddr().String(), c.User(), string(pass), c.ClientVersion())
		return nil, fmt.Errorf("nope")
	},
	NoClientAuth: false,
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <rsa key file> <listen address>, %s id_rsa 0.0.0.0:2222, ssh-keygen -P \"\" -f id_rsa", os.Args[0], os.Args[0])
	}

	privateBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Failed to load private key (./id_rsa)")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", os.Args[2])
	if err != nil {
		log.Fatalf("Failed to listen on %s (%s)", os.Args[2], err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		go ssh.NewServerConn(conn, config)
	}
}
