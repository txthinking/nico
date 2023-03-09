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
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var maxbody int64 = 0
var timeout int64 = 0
var port int64 = 443
var rate int64 = 30
var niconame string = "github.com/txthinking/nico"
var certpath string = ""

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	if err := godotenv.Load(filepath.Join(home, ".nico.env")); err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if os.Getenv("NICO_MAX_BODY") != "" {
		maxbody, err = strconv.ParseInt(os.Getenv("NICO_MAX_BODY"), 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if os.Getenv("NICO_TIMEOUT") != "" {
		timeout, err = strconv.ParseInt(os.Getenv("NICO_TIMEOUT"), 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if os.Getenv("NICO_PORT") != "" {
		port, err = strconv.ParseInt(os.Getenv("NICO_PORT"), 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if os.Getenv("NICO_RATE") != "" {
		rate, err = strconv.ParseInt(os.Getenv("NICO_RATE"), 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
	}
	if s := os.Getenv("NICO_NAME"); s != "" {
		niconame = s
	}
	certpath = filepath.Join(home, ".nico")
	if s := os.Getenv("NICO_CERT"); s != "" {
		certpath = s
	}
}
