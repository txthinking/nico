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
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/txthinking/brook/limits"
)

var maxbody int64 = 0
var timeout int64 = 0
var port int64 = 443
var niconame string = "github.com/txthinking/nico"
var certpath string = "/root/.nico/"

func main() {
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	if err := godotenv.Load("/root/.nico.env"); err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if os.Getenv("NICO_PORT") != "" {
		var err error
		port, err = strconv.ParseInt(os.Getenv("NICO_PORT"), 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	maxbody, _ = strconv.ParseInt(os.Getenv("NICO_MAX_BODY"), 10, 64)
	timeout, _ = strconv.ParseInt(os.Getenv("NICO_TIMEOUT"), 10, 64)
	if s := os.Getenv("NICO_NAME"); s != "" {
		niconame = s
	}
	if s := os.Getenv("NICO_CERT"); s != "" {
		certpath = s
	}

	if len(os.Args) == 1 || (len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "help" || os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Print(`
Nico:

	A HTTP2 web server for reverse proxy and single page application, automatically apply for ssl certificate, zero-configuration.

Make sure your domains are already resolved to your server IP and open 80/443 port

Static server, can be used for Single Page Application:

	$ nico domain.com /path/to/web/root

Reverse proxy:

	$ nico domain.com http://127.0.0.1:2020

Reverse proxy https website:

	$ nico domain.com https://reactjs.org

Dispatch according to path, such as, exact match: domain.com/ws; prefix match when / is suffix: domain.com/api/; default match: domain.com; A special one: domain.com/ is exact match.

	$ nico domain.com /path/to/web/root domain.com/ws http://127.0.0.1:9999 domain.com/api/ http://127.0.0.1:2020

Multiple domains:

	$ nico domain0.com /path/to/web/root domain1.com /another/web/root domain1.com/ws http://127.0.0.1:9999 domain1.com/api/ http://127.0.0.1:2020

Custom certificate:

	Put your certificate to NICO_CERT, default: /root/.nico/
	- File name format
		- DOMAIN.cert.pem
		- DOMAIN.key.pem
	- Simple domain certificate
		- Example: domain.com
			- domain.com.cert.pem
			- domain.com.key.pem
		- Example: a.domain.com
			- a.domain.com.cert.pem
			- a.domain.com.key.pem
	- Wildcard domain certificate
		- Example: *.domain.com
			- .domain.com.cert.pem
			- .domain.com.key.pem
		- Example: *.a.domain.com
			- .a.domain.com.cert.pem
			- .a.domain.com.key.pem
	If nico does not find certificate for a domain name, then apply for a certificate automatically

Env variables:

	NICO_PORT: default 443
	NICO_CERT: default /root/.nico/
	NICO_MAX_BODY: Maximum body size(b)
	NICO_TIMEOUT: Read/write timeout(s)

Also supoorts dotenv on /root/.nico.env

Verson:
	v20220420

Copyright:
	https://github.com/txthinking/nico
`)
		return
	}

	ss, err := Server(os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		ss.Shutdown(context.Background())
	}()
	if err := ss.ListenAndServeTLS("", ""); err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
}
