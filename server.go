package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

func Server(ll []string) (*http.Server, error) {
	nico := NewNico()
	for _, v := range ll {
		l := strings.Split(strings.TrimSpace(v), " ")
		if len(l) != 2 {
			return nil, errors.New("Invalid format: " + v)
		}
		if err := nico.Add(l[0], l[1]); err != nil {
			return nil, err
		}
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(secure.New(secure.Options{
		AllowedHosts:            nico.Domains(),
		SSLRedirect:             false,
		STSSeconds:              315360000,
		STSIncludeSubdomains:    true,
		STSPreload:              true,
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
	}).HandlerFuncWithNext))

	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", "github.com/txthinking/nico")
		next(w, r)
	})
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024) // 3M
		next(w, r)
	})

	lmt := tollbooth.NewLimiter(30, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		httpError := tollbooth.LimitByRequest(lmt, w, r)
		if httpError != nil {
			w.Header().Add("Content-Type", lmt.GetMessageContentType())
			w.WriteHeader(httpError.StatusCode)
			w.Write([]byte(httpError.Message))
			return
		}
		next(w, r)
	})

	n.Use(gzip.Gzip(gzip.DefaultCompression))

	n.UseHandler(nico)

	m := autocert.Manager{
		Cache:      autocert.DirCache(".letsencrypt"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(nico.Domains()...),
		Email:      "cloud+nico@txthinking.com",
	}
	go http.ListenAndServe(":80", m.HTTPHandler(nil))
	return &http.Server{
		Addr:           ":443",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        n,
		ErrorLog:       log.New(&tlserr{}, "", log.LstdFlags),
		TLSConfig:      &tls.Config{GetCertificate: m.GetCertificate},
	}, nil
}
