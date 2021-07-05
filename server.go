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
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

func Server(ll []string) (*http.Server, error) {
	nico := NewNico()
	if strings.Contains(ll[0], " ") {
		for _, v := range ll {
			l := strings.Split(strings.TrimSpace(v), " ")
			if len(l) != 2 {
				return nil, errors.New("Invalid format: " + v)
			}
			if err := nico.Add(l[0], l[1]); err != nil {
				return nil, err
			}
		}
	}
	if !strings.Contains(ll[0], " ") {
		if len(ll)%2 != 0 {
			return nil, errors.New("The number of parameters should be even")
		}
		for i := 0; i < len(ll); i = i + 2 {
			if err := nico.Add(ll[i], ll[i+1]); err != nil {
				return nil, err
			}
		}
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", "github.com/txthinking/nico")
		next(w, r)
	})
	if maxbody != 0 {
		n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			r.Body = http.MaxBytesReader(w, r.Body, maxbody)
			next(w, r)
		})
	}
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
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		IdleTimeout:    time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        n,
		ErrorLog:       log.New(&tlserr{}, "", log.LstdFlags),
		TLSConfig:      &tls.Config{GetCertificate: m.GetCertificate},
	}, nil
}
