package srslog

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloseNonOpenWriter(t *testing.T) {
	w := Writer{}

	err := w.Close()
	assert.NoError(t, err, "should not fail to close if there is nothing to close")
}

func TestWriteAndRetryFails(t *testing.T) {
	w := Writer{network: "udp", raddr: "fakehost"}

	n, err := w.writeAndRetry(LOG_ERR, []byte("nope"))
	assert.Error(t, err, "should fail to write")
	assert.Equal(t, 0, n, "should not write any bytes")
}

func TestSetHostname(t *testing.T) {
	var (
		assert         = assert.New(t)
		customHostname = "kubernetesCluster"
		expected       = customHostname
	)

	done := make(chan string)
	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()
	assert.Equal("127.0.0.1", strings.Split(w.hostname, ":")[0], "hostname")

	w.SetHostname(customHostname)
	assert.Equal(expected, w.hostname, "hostname")
	<-done
}

func TestWriteFormatters(t *testing.T) {
	var (
		assert = assert.New(t)
		tests  = []struct {
			name string
			f    Formatter
		}{
			{"default", nil},
			{"unix", UnixFormatter},
			{"rfc 3164", RFC3164Formatter},
			{"rfc 5424", RFC5424Formatter},
			{"default", DefaultFormatter},
		}
	)

	for _, test := range tests {
		done := make(chan string)
		addr, sock, srvWG := startServer("udp", "", done)
		defer sock.Close()
		defer srvWG.Wait()

		w := Writer{
			priority: LOG_ERR,
			tag:      "tag",
			hostname: "hostname",
			network:  "udp",
			raddr:    addr,
		}

		_, err := w.connect()
		assert.NoError(err, "failed to connect")
		defer w.Close()

		w.SetFormatter(test.f)

		f := test.f
		if f == nil {
			f = DefaultFormatter
		}
		expected := string(f(time.Now(), LOG_ERR, "hostname", "tag", []byte("this is a test message")))

		_, err = w.Write([]byte("this is a test message"))
		assert.NoError(err, "failed to write")
		sent := strings.TrimSpace(<-done)
		assert.Equalf(expected, sent, "expected to use the %v formatter", test.name)
	}
}

func TestWriterFramers(t *testing.T) {
	var (
		assert = assert.New(t)
		tests  = []struct {
			name string
			f    Framer
		}{
			{"default", nil},
			{"rfc 5425", RFC5425MessageLengthFramer},
			{"default", DefaultFramer},
		}
	)

	for _, test := range tests {
		done := make(chan string)
		addr, sock, srvWG := startServer("udp", "", done)
		defer sock.Close()
		defer srvWG.Wait()

		w := Writer{
			priority: LOG_ERR,
			tag:      "tag",
			hostname: "hostname",
			network:  "udp",
			raddr:    addr,
		}

		_, err := w.connect()
		assert.NoError(err, "failed to connect")
		defer w.Close()

		w.SetFramer(test.f)

		f := test.f
		if f == nil {
			f = DefaultFramer
		}

		var (
			bb       = f(bytes.Join([][]byte{DefaultFormatter(time.Now(), LOG_ERR, "hostname", "tag", []byte("this is a test message")), {'\n'}}, []byte{}))
			expected = strings.TrimSpace(string(bytes.Join(bb, []byte{})))
		)

		_, err = w.Write([]byte("this is a test message"))
		assert.NoError(err, "failed to write")

		sent := strings.TrimSpace(<-done)
		assert.Equalf(expected, sent, "expected to use the %v framer", test.name)
	}
}

func TestWriteWithDefaultPriority(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	n, err := w.Write([]byte("this is a test message"))
	assert.NoError(err, "failed to write")
	assert.NotEqual(0, n, "zero bytes written")

	checkWithPriorityAndTag(t, LOG_ERR, "tag", "hostname", "this is a test message", <-done)
}

func TestWriteWithPriority(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	n, err := w.WriteWithPriority(LOG_DEBUG, []byte("this is a test message"))
	assert.NoError(err, "failed to write")
	assert.NotEqual(0, n, "zero bytes written")

	checkWithPriorityAndTag(t, LOG_DEBUG, "tag", "hostname", "this is a test message", <-done)
}

func TestWriteWithPriorityAndFacility(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	n, err := w.WriteWithPriority(LOG_DEBUG|LOG_LOCAL5, []byte("this is a test message"))
	assert.NoError(err, "failed to write")
	assert.NotEqual(0, n, "zero bytes written")

	checkWithPriorityAndTag(t, LOG_DEBUG|LOG_LOCAL5, "tag", "hostname", "this is a test message", <-done)
}

func TestDebug(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Debug("this is a test message")
	assert.NoError(err, "failed to debug")

	checkWithPriorityAndTag(t, LOG_DEBUG, "tag", "hostname", "this is a test message", <-done)
}

func TestInfo(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Info("this is a test message")
	assert.NoError(err, "failed to info")

	checkWithPriorityAndTag(t, LOG_INFO, "tag", "hostname", "this is a test message", <-done)
}

func TestNotice(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Notice("this is a test message")
	assert.NoError(err, "failed to notice")

	checkWithPriorityAndTag(t, LOG_NOTICE, "tag", "hostname", "this is a test message", <-done)
}

func TestWarning(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Warning("this is a test message")
	assert.NoError(err, "failed to warn")

	checkWithPriorityAndTag(t, LOG_WARNING, "tag", "hostname", "this is a test message", <-done)
}

func TestErr(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Err("this is a test message")
	assert.NoError(err, "failed to err")

	checkWithPriorityAndTag(t, LOG_ERR, "tag", "hostname", "this is a test message", <-done)
}

func TestCrit(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Crit("this is a test message")
	assert.NoError(err, "failed to crit")

	checkWithPriorityAndTag(t, LOG_CRIT, "tag", "hostname", "this is a test message", <-done)
}

func TestAlert(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Alert("this is a test message")
	assert.NoError(err, "failed to alert")

	checkWithPriorityAndTag(t, LOG_ALERT, "tag", "hostname", "this is a test message", <-done)
}

func TestEmerg(t *testing.T) {
	var (
		assert = assert.New(t)
		done   = make(chan string)
	)

	addr, sock, srvWG := startServer("udp", "", done)
	defer sock.Close()
	defer srvWG.Wait()

	w := Writer{
		priority: LOG_ERR,
		tag:      "tag",
		hostname: "hostname",
		network:  "udp",
		raddr:    addr,
	}

	_, err := w.connect()
	assert.NoError(err, "failed to connect")
	defer w.Close()

	err = w.Emerg("this is a test message")
	assert.NoError(err, "failed to emerg")

	checkWithPriorityAndTag(t, LOG_EMERG, "tag", "hostname", "this is a test message", <-done)
}
