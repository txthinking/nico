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
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/txthinking/brook"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
)

func main() {
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	if runtime.GOOS == "linux" {
		c := exec.Command("sysctl", "-w", "net.core.rmem_max=2500000")
		b, err := c.CombinedOutput()
		if err != nil {
			log.Println("Try to raise UDP Receive Buffer Size", "got", string(b))
		}
	}

	if len(os.Args) == 1 || (len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "help" || os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Print(`
Nico:

    A HTTP3 web server for reverse proxy and single page application, automatically apply for ssl certificate, zero-configuration.

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

Env variables or dotenv on $HOME/.nico.env:

    NICO_PORT:     default: 443
    NICO_IP:       default all
    NICO_MAX_BODY: maximum request body size(b), default: 0
    NICO_TIMEOUT:  read/write timeout(s), default: 0
    NICO_LOG:      default: false
    NICO_CERT:     default: $HOME/.nico/
    NICO_RATE:     DDoS mitigation, rate limit/second/IP, default: 30
    NICO_HTTP3:    default: true, this is temporary and will be removed in the future. If you want to disable http3, set it as false

Custom certificate in $NICO_CERT:

    - www.example.com
        - $NICO_CERT/www.example.com.cert.pem
        - $NICO_CERT/www.example.com.key.pem
    - *.example.com
        - $NICO_CERT/.example.com.cert.pem
        - $NICO_CERT/.example.com.key.pem
    - if nico does not find certificate for a domain name, then apply for a certificate automatically

Verson:
    v20240323

Copyright:
    https://github.com/txthinking/nico
`)
		return
	}

	l := os.Args
	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			if !filepath.IsAbs(os.Args[1]) {
				log.Println("It looks like you passed config file? It must be absolute path")
				os.Exit(1)
				return
			}
			l, err = brook.CAC(os.Args[1])
			if err != nil {
				log.Println(err)
				os.Exit(1)
				return
			}
			bin, err := os.Executable()
			if err != nil {
				log.Println(err)
				os.Exit(1)
				return
			}
			l = append([]string{bin}, l...)
		}
	}
	h2, h3, err := Server(l[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	g := runnergroup.New()
	g.Add(&runnergroup.Runner{
		Start: func() error {
			return h2.ListenAndServeTLS("", "")
		},
		Stop: func() error {
			return h2.Shutdown(context.Background())
		},
	})
	if nicohttp3 != "false" {
		g.Add(&runnergroup.Runner{
			Start: func() error {
				return h3.ListenAndServe()
			},
			Stop: func() error {
				return h3.Close()
			},
		})
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		g.Done()
	}()
	if err := g.Wait(); err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
}
