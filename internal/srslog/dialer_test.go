package srslog

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDialer(t *testing.T) {
	var (
		assert = assert.New(t)
		w      = Writer{
			priority: LOG_ERR,
			tag:      "tag",
			hostname: "",
			network:  "",
			raddr:    "",
		}
	)

	dialer := w.getDialer()
	assert.Equal("unixDialer", dialer.Name)

	for _, tc := range []struct {
		Network string
		Name    string
	}{
		{"tcp+tls", "tlsDialer"},
		{"tcp", "basicDialer"},
		{"udp", "basicDialer"},
		{"something else entirely", "basicDialer"},
	} {
		w.network = tc.Network
		dialer = w.getDialer()
		assert.Equal(tc.Name, dialer.Name)
	}

	w.network = "custom"
	w.customDial = func(string, string) (net.Conn, error) { return nil, nil }
	dialer = w.getDialer()
	assert.Equal("customDialer", dialer.Name)
}

func TestUnixDialer(t *testing.T) {
	var (
		assert = assert.New(t)
		w      = Writer{
			priority: LOG_ERR,
			tag:      "tag",
			hostname: "",
			network:  "",
			raddr:    "",
		}
	)

	_, hostname, err := w.unixDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("localhost", hostname, "should set blank hostname")

	w.hostname = "my other hostname"
	_, hostname, err = w.unixDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("my other hostname", hostname, "should not interfere with hostname")
}

func TestTLSDialer(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, _ := startServer("tcp+tls", "", done)
	defer sock.Close()

	pool := x509.NewCertPool()
	serverCert, err := ioutil.ReadFile("test/cert.pem")
	assert.NoError(err, "failed to read file")
	pool.AppendCertsFromPEM(serverCert)
	config := tls.Config{
		RootCAs: pool,
	}

	w := Writer{
		priority:  LOG_ERR,
		tag:       "tag",
		hostname:  "",
		network:   "tcp+tls",
		raddr:     addr,
		tlsConfig: &config,
	}

	_, hostname, err := w.tlsDialer()
	assert.NoError(err, "failed to dial")
	assert.NotEqual("", hostname, "should set default hostname")

	w.hostname = "my other hostname"
	_, hostname, err = w.tlsDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("my other hostname", hostname, "should not interfere with hostname")
}

func TestTCPDialer(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, _ := startServer("tcp", "", done)
	defer sock.Close()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "",
		network:  "tcp",
		raddr:    addr,
	}

	_, hostname, err := w.basicDialer()
	assert.NoError(err, "failed to dial")
	assert.NotEqual("", hostname, "should set default hostname")

	w.hostname = "my other hostname"
	_, hostname, err = w.basicDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("my other hostname", hostname, "should not interfere with hostname")
}

func TestUDPDialer(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, _ := startServer("udp", "", done)
	defer sock.Close()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "",
		network:  "udp",
		raddr:    addr,
	}

	_, hostname, err := w.basicDialer()
	assert.NoError(err, "failed to dial")
	assert.NotEqual("", hostname, "should set default hostname")

	w.hostname = "my other hostname"
	_, hostname, err = w.basicDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("my other hostname", hostname, "should not interfere with hostname")
}

func TestCustomDialer(t *testing.T) {
	// A custom dialer can really be anything, so we don't test an actual connection
	// instead we test the behavior of this code path
	var (
		assert      = assert.New(t)
		nwork, addr = "custom", "custom_addr_to_pass"
		w           = Writer{
			priority: LOG_ERR,
			tag:      "tag",
			hostname: "",
			network:  nwork,
			raddr:    addr,
			customDial: func(n string, a string) (net.Conn, error) {
				if n != nwork || a != addr {
					return nil, errors.New("Unexpected network or address, expected: (" +
						nwork + ":" + addr + ") but received (" + n + ":" + a + ")")
				}
				return fakeConn{addr: &fakeAddr{nwork, addr}}, nil
			},
		}
	)

	_, hostname, err := w.customDialer()
	assert.NoError(err, "failed to dial")
	assert.NotEqual("", hostname, "should set default hostname")

	w.hostname = "my other hostname"
	_, hostname, err = w.customDialer()
	assert.NoError(err, "failed to dial")
	assert.Equal("my other hostname", hostname, "should not interfere with hostname")
}

type fakeConn struct {
	net.Conn
	addr net.Addr
}

func (fc fakeConn) Close() error {
	return nil
}

func (fc fakeConn) Write(p []byte) (int, error) {
	return len(p), nil
}

func (fc fakeConn) LocalAddr() net.Addr {
	return fc.addr
}

type fakeAddr struct {
	nwork, addr string
}

func (fa *fakeAddr) Network() string {
	return fa.nwork
}

func (fa *fakeAddr) String() string {
	return fa.addr
}
