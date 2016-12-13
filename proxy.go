package main

import (
	"crypto/tls"
	"io"
	"net"
)

// forward conn to another conn of a remote address.
func forward(conn net.Conn, remote string) error {
	bconn, err := net.Dial("tcp", remote)
	if err != nil {
		return err
	}
	defer bconn.Close()

	go func() {
		defer conn.Close()
		io.Copy(conn, bconn)
	}()
	io.Copy(bconn, conn)
	return nil
}

// listenAndServe accepts connections on addr and forwards them to another backend.
func listenAndServe(backend, addr string, conf *tls.Config) error {
	ln, err := tls.Listen("tcp", addr, conf)
	if err != nil {
		return err
	}

	connCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})
	defer close(doneCh)
	go func() {
		defer close(connCh)
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			select {
			case connCh <- conn:
			case <-doneCh:
				return
			}
		}
	}()

	for {
		select {
		case conn, ok := <-connCh:
			if !ok {
				return io.EOF
			}
			go func(conn net.Conn) {
				if err := forward(conn, backend); err != nil {
					select {
					case errCh <- err:
					default:
					}
				}
			}(conn)
		case err := <-errCh:
			return err
		}
	}
}
