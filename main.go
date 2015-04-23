package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: wscat <url>")
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("Fail to parse:", err)
	}

	if _, _, err := net.SplitHostPort(u.Host); err != nil {
		// no port specified
		switch u.Scheme {
		case "wss", "https", "":
			u.Host += ":443"
		case "ws", "http":
			u.Host += ":80"
		default:
			log.Fatal("Unsupported URL scheme: %q", u.Scheme)
		}
	}

	netConn, err := tls.Dial("tcp", u.Host, &tls.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer netConn.Close()

	socket, _, err := websocket.NewClient(netConn, u, http.Header{}, 1024, 1024)
	if err != nil {
		log.Fatal("WS request failed:", err)
	}

	for {
		_, reader, err := socket.NextReader()
		if err != nil {
			log.Fatal("Error:", err)
		}
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			log.Fatal("Error:", err)
		}
		fmt.Fprintln(os.Stdout, "")
	}
}