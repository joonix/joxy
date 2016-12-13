package main

import (
	"fmt"
	"net"
	"testing"
)

func TestForward(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	backend := ln.Addr().String()
	t.Log(backend)

	acceptCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	backendDone := make(chan struct{})
	go func() {
		defer close(backendDone)
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
		}
		acceptCh <- conn
	}()

	frontend, conn := net.Pipe()
	frontendDone := make(chan struct{})
	go func() {
		defer close(frontendDone)
		errCh <- forward(conn, backend)
	}()

	select {
	case backendConn := <-acceptCh:
		fmt.Fprint(frontend, "hello world")
		b := make([]byte, 64)
		n, err := backendConn.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		if string(b[:n]) != "hello world" {
			t.Error("forwarded bytes mismatch")
		}
	case err := <-errCh:
		t.Fatal(err)
	}
	frontend.Close()
	<-frontendDone
	<-backendDone
}
