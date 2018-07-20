package main

import (
	"github.com/alexflint/go-arg"
	"github.com/skiloop/tcpproxy/server"
	"time"
	"syscall"
	"fmt"
	"os"
	"sync"
)

type Args struct {
	PortMaps    []string      `arg:"-m,required,separate",help:"port maps, example: 0:9998:192.168.1.99:8991"`
	Limit       uint64        `arg:"-l",help:"soft open files limit"`
	KeepAlive   bool          `arg:"-k",help:"make connection keep alive"`
	AlivePeriod time.Duration `arg:"-a",help:"keepalive period in second"`
}

func main() {
	var args Args
	args.AlivePeriod = 180

	arg.MustParse(&args)

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
	var wg sync.WaitGroup
	fmt.Println(args.PortMaps)
	for _, m := range args.PortMaps {
		localHost, localPort, remoteHost, remotePort, err := server.ParseMapString(m)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		fmt.Println("start server: " + localHost)
		wg.Add(1)
		go func() {
			defer func(s string) {
				fmt.Println(s + " done")
				wg.Done()
			}(m)
			srv := server.NewServer(localHost, localPort, remoteHost, remotePort, args.KeepAlive, args.AlivePeriod, nil)
			if err := srv.Serve(); err != nil {
				fmt.Println(err.Error())
			}
		}()
	}
	time.Sleep(time.Second * 5)
	wg.Wait()
	fmt.Println("All done")
}
