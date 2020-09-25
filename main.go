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
)

var maxbody int64 = 0
var timeout int64 = 0

func main() {
	maxbody, _ = strconv.ParseInt(os.Getenv("NICO_MAX_BODY"), 10, 64)
	timeout, _ = strconv.ParseInt(os.Getenv("NICO_TIMEOUT"), 10, 64)

	if len(os.Args) == 1 || (len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "help" || os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Print(`
Nico:

	A HTTP2 web server for reverse proxy and single page application, automatically apply for ssl certificate, zero-configuration.

Make sure your domains are already resolved to your server IP and open 80/443 port

Reverse proxy:

	$ nico 'domain.com http://127.0.0.1:2020'

Static server, can be used for Single Page Application:

	$ nico 'domain.com /path/to/web/root'

Support multiple domain in one command:

	$ nico 'domain1.com http://127.0.0.1:2020' 'domain2.com /path/to/web/root'

Verson:
	v20200925

Copyright:
	https://github.com/txthinking/nico
`)
		return
	}

	ss, err := Server(os.Args[1:])
	if err != nil {
		log.Println(err)
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
		return
	}
}
