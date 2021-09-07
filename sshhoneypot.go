package main

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net"
	"time"
)

// ssh-keygen -P "" -f id_rsa

var config = &ssh.ServerConfig {
	PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		log.Printf("%s - %s:%s (%s)\n", c.RemoteAddr().String(), c.User(), string(pass), c.ClientVersion())
		if c.User() == "foo" && string(pass) == "bar" {
			return nil, nil
		}
		return nil, fmt.Errorf("password rejected for %q", c.User())
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
		log.Fatalf("Failed to listen on 2200 (%s)", err)
	}

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		go handleConnection(tcpConn)
	}
}

func handleConnection(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(5*time.Second))

	// Before use, a handshake must be performed on the incoming net.Conn.
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		//log.Printf("Failed to handshake (%s)", err)
		return
	}

	log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
	// Discard all global out-of-band Requests
	go ssh.DiscardRequests(reqs)
	// Accept all channels
	go handleChannels(chans)
}

func handleChannels(chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		go newChannel.Reject(ssh.ResourceShortage, "nope")
	}
}
