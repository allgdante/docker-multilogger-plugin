package srslog

import (
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

// unixSyslog opens a connection to the syslog daemon running on the
// local machine using a Unix domain socket. This function exists because of
// Solaris support as implemented by gccgo.  On Solaris you can not
// simply open a TCP connection to the syslog daemon.  The gccgo
// sources have a syslog_solaris.go file that implements unixSyslog to
// return a type that satisfies the serverConn interface and simply calls the C
// library syslog function.
func unixSyslog() (conn serverConn, err error) {
	logTypes := []string{"unixgram", "unix"}
	logPaths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}
	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.Dial(network, path)
			if err != nil {
				continue
			} else {
				return &localConn{conn: conn}, nil
			}
		}
	}
	return nil, errors.New("Unix syslog delivery error")
}

// localConn adheres to the serverConn interface, allowing us to send syslog
// messages to the local syslog daemon over a Unix domain socket.
type localConn struct {
	conn io.WriteCloser
}

// write formats syslog messages using time.Stamp instead of time.RFC3339,
// and omits the hostname (because it is expected to be used locally).
func (n *localConn) write(
	framer Framer,
	formatter Formatter,
	timestamp time.Time,
	p Priority,
	hostname, tag string,
	msg []byte) error {
	if framer == nil {
		framer = DefaultFramer
	}

	if formatter == nil {
		formatter = UnixFormatter
	}

	fmsg := framer(formatter(timestamp, p, hostname, tag, msg))

	var wmsg []byte
	switch len(fmsg) {
	case 1:
		wmsg = fmsg[0]
	default:
		wmsg = bytes.Join(fmsg, []byte{})
	}
	_, err := n.conn.Write(wmsg)
	return err
}

// close the (local) network connection
func (n *localConn) close() error {
	return n.conn.Close()
}
