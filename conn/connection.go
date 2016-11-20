package conn

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Conn interface {
	io.Closer
	ReadContinuously(chan string)
	Send(string)
}

type tcpConn struct {
	Conn net.Conn
}

func New(serverName string) *tcpConn {
	conn, err := net.Dial("tcp", serverName)
	if err != nil {
		log.Fatalf("ERROR - %v\n", err)
	}

	return &tcpConn{Conn: conn}
}

func (conn *tcpConn) ReadContinuously(respCh chan string) {
	reader := bufio.NewReader(conn.Conn)
	for {
		serverResponse, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("ERROR - Reading from connection - %v\n", err)
			return
		}
		respCh <- serverResponse
	}
}

func (conn *tcpConn) Send(msg string) {
	_, err := fmt.Fprintf(conn.Conn, msg)
	if err != nil {
		log.Fatalf("Error while sending message - %v\n", err)
	}
}

func (conn *tcpConn) Close() error {
	return conn.Conn.Close()
}
