package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/context"
)

func forward(conn net.Conn, raddr *net.TCPAddr) error {
	bconn, err := net.DialTCP("tcp", nil, raddr)
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
	raddr, err := net.ResolveTCPAddr("tcp", backend)
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
				if err := forward(conn, raddr); err != nil {
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

func main() {
	certpath := flag.String("certpath", "certs", "path of where to cache certificates")
	domain := flag.String("domain", "", "the domain name to request certificate for")
	backend := flag.String("backend", "", "backend service address to forward to")
	flag.Parse()

	if *domain == "" {
		log.Fatalln("no domain set")
	}
	if *backend == "" {
		log.Fatalln("no backend set")
	}
	var cache autocert.Cache
	if *certpath != "" {
		cache = autocert.DirCache(*certpath)
	}

	whitelist := autocert.HostWhitelist(*domain)
	m := autocert.Manager{
		Cache:  cache,
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			log.Println("checking whitelist for", host)
			return whitelist(ctx, host)
		},
	}
	conf := tls.Config{GetCertificate: m.GetCertificate}

	log.Fatalln(listenAndServe(*backend, ":https", &conf))
}
