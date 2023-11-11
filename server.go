// Copyright (c) 2020-present Cloud <cloud@txthinking.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/quic-go/quic-go/http3"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

func Server(ll []string) (*http.Server, *http3.Server, error) {
	nico := NewNico()
	if strings.Contains(ll[0], " ") {
		for _, v := range ll {
			l := strings.Split(strings.TrimSpace(v), " ")
			if len(l) != 2 {
				return nil, nil, errors.New("Invalid format: " + v)
			}
			if err := nico.Add(l[0], l[1]); err != nil {
				return nil, nil, err
			}
		}
	}
	if !strings.Contains(ll[0], " ") {
		if len(ll)%2 != 0 {
			return nil, nil, errors.New("The number of parameters should be even")
		}
		for i := 0; i < len(ll); i = i + 2 {
			if err := nico.Add(ll[i], ll[i+1]); err != nil {
				return nil, nil, err
			}
		}
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	if os.Getenv("NICO_LOG") == "true" {
		n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			next(w, r)
			fmt.Printf(`{"from": "%s", "method": "%s", "hostname": "%s", "path": "%s", "status": "%d", "time": "%s"}`+"\n", r.RemoteAddr, r.Method, r.Host, r.URL.Path, w.(negroni.ResponseWriter).Status(), time.Now().Format(time.RFC3339))
		})
	}
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", niconame)
		next(w, r)
	})
	if nicohttp3 != "false" {
		n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			w.Header().Set("Alt-Svc", `h3=":`+strconv.FormatInt(port, 10)+`"; ma=2592000`)
			next(w, r)
		})
	}
	if maxbody != 0 {
		n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			r.Body = http.MaxBytesReader(w, r.Body, maxbody)
			next(w, r)
		})
	}
	lmt := tollbooth.NewLimiter(float64(rate), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr"})
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

	l := make([]tls.Certificate, 0)
	certs := make(map[string]*tls.Certificate)
	auto := make([]string, 0)
	for _, v := range nico.Domains() {
		c, err := os.ReadFile(filepath.Join(certpath, v+".cert.pem"))
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, err
		}
		k, err := os.ReadFile(filepath.Join(certpath, v+".key.pem"))
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, err
		}
		if c != nil && k != nil {
			ct, err := tls.X509KeyPair(c, k)
			if err != nil {
				return nil, nil, err
			}
			certs[v] = &ct
			if net.ParseIP(v) != nil {
				l = append(l, ct)
			}
			continue
		}
		if strings.Index(v, ".") != -1 {
			c, err := os.ReadFile(filepath.Join(certpath, v[strings.Index(v, "."):]+".cert.pem"))
			if err != nil && !os.IsNotExist(err) {
				return nil, nil, err
			}
			k, err := os.ReadFile(filepath.Join(certpath, v[strings.Index(v, "."):]+".key.pem"))
			if err != nil && !os.IsNotExist(err) {
				return nil, nil, err
			}
			if c != nil && k != nil {
				ct, err := tls.X509KeyPair(c, k)
				if err != nil {
					return nil, nil, err
				}
				certs[v] = &ct
				continue
			}
		}
		auto = append(auto, v)
	}

	var m autocert.Manager
	if len(auto) != 0 {
		m = autocert.Manager{
			Cache:      autocert.DirCache(".letsencrypt"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(auto...),
			Email:      "cloud+nico@txthinking.com",
		}
		go http.ListenAndServe(net.JoinHostPort(nicoip, "80"), m.HTTPHandler(nil))
	}
	tc := &tls.Config{
		Certificates: l,
		GetCertificate: func(c *tls.ClientHelloInfo) (*tls.Certificate, error) {
			v, ok := certs[c.ServerName]
			if ok {
				return v, nil
			}
			if len(auto) != 0 {
				return m.GetCertificate(c)
			}
			return nil, errors.New("Not found " + c.ServerName)
		},
	}
	return &http.Server{
			Addr:           net.JoinHostPort(nicoip, strconv.FormatInt(port, 10)),
			ReadTimeout:    time.Duration(timeout) * time.Second,
			WriteTimeout:   time.Duration(timeout) * time.Second,
			IdleTimeout:    time.Duration(timeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        n,
			ErrorLog:       log.New(&tlserr{}, "", log.LstdFlags),
			TLSConfig:      tc,
		}, &http3.Server{
			Addr:      net.JoinHostPort(nicoip, strconv.FormatInt(port, 10)),
			TLSConfig: tc,
			Handler:   n,
		}, nil
}
