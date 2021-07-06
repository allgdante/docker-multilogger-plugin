package srslog

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const appNameMaxLength = 48 // limit to 48 chars as per RFC5424

// Formatter is a type of function that takes the constituent parts of a
// syslog message and returns a formatted string. A different Formatter is
// defined for each different syslog protocol we support.
type Formatter func(timestamp time.Time, p Priority, hostname, tag string, content []byte) []byte

// DefaultFormatter is the original format supported by the Go syslog package,
// and is a non-compliant amalgamation of 3164 and 5424 that is intended to
// maximize compatibility.
func DefaultFormatter(timestamp time.Time, p Priority, hostname, tag string, content []byte) []byte {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "<%d> %s %s %s[%d]: ",
		p, timestamp.Format(time.RFC3339), hostname, tag, os.Getpid())
	b.Write(content)
	return b.Bytes()
}

// UnixFormatter omits the hostname, because it is only used locally.
func UnixFormatter(timestamp time.Time, p Priority, _, tag string, content []byte) []byte {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "<%d>%s %s[%d]: ",
		p, timestamp.Format(time.Stamp), tag, os.Getpid())
	b.Write(content)
	return b.Bytes()
}

// RFC3164Formatter provides an RFC 3164 compliant message.
func RFC3164Formatter(timestamp time.Time, p Priority, hostname, tag string, content []byte) []byte {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "<%d>%s %s %s[%d]: ",
		p, timestamp.Format(time.Stamp), hostname, tag, os.Getpid())
	b.Write(content)
	return b.Bytes()
}

// if string's length is greater than max, then use the last part
func truncateStartStr(s string, max int) string {
	if len(s) > max {
		return s[len(s)-max:]
	}
	return s
}

// RFC5424Formatter provides an RFC 5424 compliant message.
func RFC5424Formatter(timestamp time.Time, p Priority, hostname, tag string, content []byte) []byte {
	var (
		b       = new(bytes.Buffer)
		pid     = os.Getpid()
		appName = truncateStartStr(os.Args[0], appNameMaxLength)
	)

	fmt.Fprintf(b, "<%d>%d %s %s %s %d %s - ",
		p, 1, timestamp.Format(time.RFC3339), hostname, appName, pid, tag)
	b.Write(content)
	return b.Bytes()
}
