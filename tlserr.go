package main

import (
	"fmt"
	"strings"
)

type tlserr struct {
}

func (l *tlserr) Write(p []byte) (int, error) {
	if strings.Contains(string(p), "TLS handshake error") {
		return 0, nil
	}
	fmt.Printf("%s\n", p)
	return 0, nil
}
