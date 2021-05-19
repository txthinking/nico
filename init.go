package main

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
