package srslog

import (
	"bytes"
	"crypto/tls"
	"sync"
	"time"
)

// A Writer is a connection to a syslog server.
type Writer struct {
	priority  Priority
	tag       string
	hostname  string
	network   string
	raddr     string
	tlsConfig *tls.Config
	framer    Framer
	formatter Formatter

	//non-nil if custom dialer set, used in getDialer
	customDial DialFunc

	mu   sync.RWMutex // guards conn
	conn serverConn
}

// getConn provides access to the internal conn, protected by a mutex. The
// conn is threadsafe, so it can be used while unlocked, but we want to avoid
// race conditions on grabbing a reference to it.
func (w *Writer) getConn() serverConn {
	w.mu.RLock()
	conn := w.conn
	w.mu.RUnlock()
	return conn
}

// setConn updates the internal conn, protected by a mutex.
func (w *Writer) setConn(c serverConn) {
	w.mu.Lock()
	w.conn = c
	w.mu.Unlock()
}

// connect makes a connection to the syslog server.
func (w *Writer) connect() (serverConn, error) {
	conn := w.getConn()
	if conn != nil {
		// ignore err from close, it makes sense to continue anyway
		conn.close()
		w.setConn(nil)
	}

	var hostname string
	var err error
	dialer := w.getDialer()
	conn, hostname, err = dialer.Call()
	if err == nil {
		w.setConn(conn)
		w.hostname = hostname

		return conn, nil
	}
	return nil, err
}

// SetFormatter changes the formatter function for subsequent messages.
func (w *Writer) SetFormatter(f Formatter) {
	w.formatter = f
}

// SetFramer changes the framer function for subsequent messages.
func (w *Writer) SetFramer(f Framer) {
	w.framer = f
}

// SetHostname changes the hostname for syslog messages if needed.
func (w *Writer) SetHostname(hostname string) {
	w.hostname = hostname
}

// Write sends a log message to the syslog daemon using the default priority
// passed into `srslog.New` or the `srslog.Dial*` functions.
func (w *Writer) Write(b []byte) (int, error) {
	return w.writeAndRetry(w.priority, b)
}

// WriteWithTimestamp sends a log message with a custom timestamp.
func (w *Writer) WriteWithTimestamp(timestamp time.Time, b []byte) (int, error) {
	return w.writeAndRetryWithTimestamp(timestamp, w.priority, b)
}

// WriteWithPriority sends a log message with a custom priority.
func (w *Writer) WriteWithPriority(p Priority, b []byte) (int, error) {
	return w.writeAndRetryWithPriority(p, b)
}

// WriteWithTimestampAndPriority sends a log message with custom timestamp and priority.
func (w *Writer) WriteWithTimestampAndPriority(timestamp time.Time, p Priority, b []byte) (int, error) {
	return w.writeAndRetryWithTimestampAndPriority(timestamp, p, b)
}

// Close closes a connection to the syslog daemon.
func (w *Writer) Close() error {
	conn := w.getConn()
	if conn != nil {
		err := conn.close()
		w.setConn(nil)
		return err
	}
	return nil
}

// Emerg logs a message with severity LOG_EMERG; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Emerg(m string) (err error) {
	_, err = w.writeAndRetry(LOG_EMERG, []byte(m))
	return err
}

// Alert logs a message with severity LOG_ALERT; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Alert(m string) (err error) {
	_, err = w.writeAndRetry(LOG_ALERT, []byte(m))
	return err
}

// Crit logs a message with severity LOG_CRIT; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Crit(m string) (err error) {
	_, err = w.writeAndRetry(LOG_CRIT, []byte(m))
	return err
}

// Err logs a message with severity LOG_ERR; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Err(m string) (err error) {
	_, err = w.writeAndRetry(LOG_ERR, []byte(m))
	return err
}

// Warning logs a message with severity LOG_WARNING; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Warning(m string) (err error) {
	_, err = w.writeAndRetry(LOG_WARNING, []byte(m))
	return err
}

// Notice logs a message with severity LOG_NOTICE; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Notice(m string) (err error) {
	_, err = w.writeAndRetry(LOG_NOTICE, []byte(m))
	return err
}

// Info logs a message with severity LOG_INFO; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Info(m string) (err error) {
	_, err = w.writeAndRetry(LOG_INFO, []byte(m))
	return err
}

// Debug logs a message with severity LOG_DEBUG; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Debug(m string) (err error) {
	_, err = w.writeAndRetry(LOG_DEBUG, []byte(m))
	return err
}

// writeAndRetry takes a severity and the string to write. Any facility passed to
// it as part of the severity Priority will be ignored.
func (w *Writer) writeAndRetry(severity Priority, b []byte) (int, error) {
	return w.writeAndRetryWithTimestamp(time.Now(), severity, b)
}

// writeAndRetryWithTimestamp differs from writeAndRetry in that it allows setting
// the timestamp
func (w *Writer) writeAndRetryWithTimestamp(timestamp time.Time, severity Priority, b []byte) (int, error) {
	pr := (w.priority & FacilityMask) | (severity & SeverityMask)

	return w.writeAndRetryWithTimestampAndPriority(timestamp, pr, b)
}

// writeAndRetryWithPriority differs from writeAndRetry in that it allows setting
// of both the facility and the severity.
func (w *Writer) writeAndRetryWithPriority(p Priority, b []byte) (int, error) {
	return w.writeAndRetryWithTimestampAndPriority(time.Now(), p, b)
}

// writeAndRetryWithTimestampAndPriority differs from writeAndRetryWithPriority in that it
// allows setting the timestap
func (w *Writer) writeAndRetryWithTimestampAndPriority(
	timestamp time.Time,
	p Priority,
	b []byte) (int, error) {
	conn := w.getConn()
	if conn != nil {
		if n, err := w.write(conn, timestamp, p, b); err == nil {
			return n, err
		}
	}

	var err error
	if conn, err = w.connect(); err != nil {
		return 0, err
	}
	return w.write(conn, timestamp, p, b)
}

// write generates and writes a syslog formatted string. It formats the
// message based on the current Formatter and Framer.
func (w *Writer) write(conn serverConn, timestamp time.Time, p Priority, msg []byte) (int, error) {
	// ensure it ends in a \n
	if !bytes.HasSuffix(msg, []byte("\n")) {
		msg = append(msg, byte('\n'))
	}

	err := conn.write(w.framer, w.formatter, timestamp, p, w.hostname, w.tag, msg)
	if err != nil {
		return 0, err
	}

	return len(msg), nil
}
