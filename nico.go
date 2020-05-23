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
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Nico struct {
	Handlers map[string]http.Handler
}

func NewNico() *Nico {
	return &Nico{
		Handlers: make(map[string]http.Handler),
	}
}

func (n *Nico) Add(domain, to string) error {
	if !strings.HasPrefix(to, "http://") && !strings.HasPrefix(to, "https://") {
		n.Handlers[domain] = http.FileServer(NewWebRoot(to))
	}
	if strings.HasPrefix(to, "http://") || strings.HasPrefix(to, "https://") {
		u, err := url.Parse(to)
		if err != nil {
			return err
		}
		n.Handlers[domain] = httputil.NewSingleHostReverseProxy(u)
	}
	return nil
}

func (n *Nico) Domains() []string {
	l := make([]string, 0)
	for k, _ := range n.Handlers {
		l = append(l, k)
	}
	return l
}

func (n *Nico) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := r.Host
	if strings.Contains(h, ":") {
		var err error
		h, _, err = net.SplitHostPort(h)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	f, ok := n.Handlers[h]
	if !ok {
		http.Error(w, "", 400)
		return
	}
	f.ServeHTTP(w, r)
}
