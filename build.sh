#!/bin/bash

if [ $# -ne 1 ]; then
    echo "./build.sh version"
    exit
fi

mkdir _

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -trimpath --ldflags="-w -s" -o _/nico_darwin_amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go  build -trimpath --ldflags="-w -s" -o _/nico_darwin_arm64
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go   build -trimpath --ldflags="-w -s" -o _/nico_freebsd_386
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath --ldflags="-w -s" -o _/nico_freebsd_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go     build -trimpath --ldflags="-w -s" -o _/nico_linux_386
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go   build -trimpath --ldflags="-w -s" -o _/nico_linux_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go   build -trimpath --ldflags="-w -s" -o _/nico_linux_arm64
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go    build -trimpath --ldflags="-w -s" -o _/nico_netbsd_386
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go  build -trimpath --ldflags="-w -s" -o _/nico_netbsd_amd64
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go   build -trimpath --ldflags="-w -s" -o _/nico_openbsd_386
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -trimpath --ldflags="-w -s" -o _/nico_openbsd_amd64
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -trimpath --ldflags="-w -s" -o _/nico_openbsd_arm64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath --ldflags="-w -s" -o _/nico_windows_amd64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go   build -trimpath --ldflags="-w -s" -o _/nico_windows_386.exe

nami release github.com/txthinking/nico $1 _

rm -rf _
