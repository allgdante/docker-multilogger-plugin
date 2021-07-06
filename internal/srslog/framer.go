package srslog

import "strconv"

// Framer is a type of function that takes an input string (typically an
// already-formatted syslog message) and applies "message framing" to it. We
// have different framers because different versions of the syslog protocol
// and its transport requirements define different framing behavior.
type Framer func(in []byte) [][]byte

// DefaultFramer does nothing, since there is no framing to apply. This is
// the original behavior of the Go syslog package, and is also typically used
// for UDP syslog.
func DefaultFramer(in []byte) [][]byte {
	return [][]byte{in}
}

// RFC5425MessageLengthFramer prepends the message length to the front of the
// provided message, as defined in RFC 5425.
func RFC5425MessageLengthFramer(in []byte) [][]byte {
	return [][]byte{
		[]byte(strconv.Itoa(len(in)) + " "),
		in,
	}
}
