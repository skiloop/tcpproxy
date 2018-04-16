package main

import (
	"github.com/alexflint/go-arg"
	"github.com/skiloop/tcpproxy/server"
	"time"
	"syscall"
	"fmt"
	"os"
)

type Args struct {
	LocalHost   string `arg:"required,positional",help:"local host"`
	LocalPort   int64  `arg:"required,positional",help:"local port"`
	Limit       uint64 `arg:"-l",help:"soft open files limit"`
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
	if args.Limit > 0 {
		rlimit := syscall.Rlimit{}
		err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
		if err == nil {
			rlimit.Cur = args.Limit
			if rlimit.Max < rlimit.Cur {
				rlimit.Max = rlimit.Cur
			}
			err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
		}
		if err != nil {
			fmt.Fprint(os.Stderr, "failed to set open files limit")
		}
	}
	srv := server.NewServer(args.LocalHost, args.LocalPort, args.RemoteHost, args.RemotePort, args.KeepAlive, alivePeriod, nil)
	srv.Serve()
}
