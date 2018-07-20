package server

import (
	"net"
	"github.com/skiloop/goutils"
	"fmt"
	"time"
	"github.com/op/go-logging"
	"strings"
	"errors"
	"regexp"
	"strconv"
)

// Server which listen on a specific port and accept tcp connection
// then send to target server
type Server struct {
	LocalHost  string
	LocalPort  int64
	remoteAddr string

	enableKeepAlive bool
	keepAlivePeriod time.Duration

	log logging.Logger

	listener net.Listener
}

func ParseMapString(m string) (localHost string, localPort int64, remoteHost string, remotePort int64, err error) {
	vars := strings.Split(m, ":")
	regex, _ := regexp.Compile("^\\d+$")
	if len(vars) < 4 || !regex.MatchString(vars[1]) || !regex.MatchString(vars[3]) {
		return "", 0, "", 0, errors.New("format error")
	}
	localPort, err = strconv.ParseInt(vars[1], 10, 64)
	if err != nil {
		return "", 0, "", 0, err
	}
	remotePort, err = strconv.ParseInt(vars[1], 10, 64)
	if err != nil {
		return "", 0, "", 0, err
	}
	return vars[0], localPort, vars[2], remotePort, nil
}

func NewServer(localHost string, localPort int64, remoteHost string, remotePort int64, enableKeepAlive bool, keepAlivePeriod time.Duration, logger *logging.Logger) Server {
	if logger == nil {
		logger = logging.MustGetLogger("tcpproxy")
	}
	return Server{LocalHost: localHost, LocalPort: localPort,
		remoteAddr: fmt.Sprintf("%s:%d", remoteHost, remotePort),
		enableKeepAlive: enableKeepAlive, keepAlivePeriod: keepAlivePeriod,
		log: *logger}
}

func (srv *Server) Serve() error {
	err := srv.init()
	if err != nil {
		return err
	}
	var delay time.Duration
	for {
		con, err := srv.listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				if delay == 0 {
					delay = 5 * time.Millisecond
				} else {
					delay *= 2
				}
				if max := time.Second; delay > max {
					delay = max
				}

				srv.log.Debugf("tcpproxy: temporary error on accept: %v", err)
				time.Sleep(delay)
				continue
			}

			srv.log.Errorf("tcpproxy: failed to accept: %v", err)
			return err
		}
		if srv.enableKeepAlive {
			if tconn, ok := con.(*net.TCPConn); ok {
				tconn.SetKeepAlive(true)
				tconn.SetKeepAlivePeriod(srv.keepAlivePeriod)
			}
		}
		go srv.work(con)
	}
	return nil
}

func (srv *Server) init() error {
	srv.log.Info("init listener")
	addr := fmt.Sprintf("%s:%d", srv.LocalHost, srv.LocalPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv.listener = l
	return nil
}

func (srv *Server) work(local net.Conn) {
	defer local.Close()
	remote, err := srv.connect()
	if err != nil {
		return
	}
	defer remote.Close()
	goutils.Join(local, remote)
}

func (srv *Server) connect() (con net.Conn, err error) {
	return net.Dial("tcp", srv.remoteAddr)
}
