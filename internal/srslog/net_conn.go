package srslog

import (
	"bytes"
	"net"
	"time"
)

// netConn has an internal net.Conn and adheres to the serverConn interface,
// allowing us to send syslog messages over the network.
type netConn struct {
	conn net.Conn
}

// write formats syslog messages using time.RFC3339 and includes the
// hostname, and sends the message to the connection.
func (n *netConn) write(
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
		formatter = DefaultFormatter
	}

	fmsg := framer(formatter(timestamp, p, hostname, tag, msg))

	var err error
	switch conn := n.conn.(type) {
	case *net.TCPConn:
		nb := net.Buffers(fmsg)
		_, err = nb.WriteTo(conn)
	default:
		var wmsg []byte
		switch len(fmsg) {
		case 1:
			wmsg = fmsg[0]
		default:
			wmsg = bytes.Join(fmsg, []byte{})
		}
		_, err = conn.Write(wmsg)
	}

	return err
}

// close the network connection
func (n *netConn) close() error {
	return n.conn.Close()
}
