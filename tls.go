package main

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func TLS(domains []string) *tls.Config {
	m := autocert.Manager{
		Cache:      autocert.DirCache(".letsencrypt"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Email:      "cloud+nico@txthinking.com",
	}
	go http.ListenAndServe(":80", m.HTTPHandler(nil))
	return &tls.Config{GetCertificate: m.GetCertificate}
}
