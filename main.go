package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	domain := flag.String("domain", "", "the domain name to request certificate for")
	flag.Parse()

	if *domain == "" {
		log.Fatalln("no domain set")
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(*domain),
	}

	s := http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}
	log.Println(s.ListenAndServeTLS("", ""))
}
