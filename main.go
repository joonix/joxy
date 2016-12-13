package main

import (
	"crypto/tls"
	"flag"
	"log"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/context"
)

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
