package main

import (
	"github.com/alexflint/go-arg"
	"github.com/skiloop/tcpproxy/server"
	"time"
)

type Args struct {
	LocalHost   string `arg:"required,positional",help:"local host"`
	LocalPort   int64  `arg:"required,positional",help:"local port"`
	RemoteHost  string `arg:"required,positional",help:"remote host"`
	RemotePort  int64  `arg:"required,positional",help:"remote port"`
	KeepAlive   bool   `arg:"-k",help:"make connection keep alive"`
	AlivePeriod int    `arg:"-a",help:"keepalive period in second"`
}

func main() {
	var args Args
	args.AlivePeriod = 180

	arg.MustParse(&args)
	alivePeriod := time.Duration(time.Second.Nanoseconds() * int64(args.AlivePeriod))
	srv := server.NewServer(args.LocalHost, args.LocalPort, args.RemoteHost, args.RemotePort, args.KeepAlive, alivePeriod, nil)
	srv.Serve()
}
