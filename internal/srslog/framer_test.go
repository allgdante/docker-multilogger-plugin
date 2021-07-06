package srslog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFramer(t *testing.T) {
	msg := []byte("input message")
	out := DefaultFramer(msg)
	assert.Equal(t, [][]byte{msg}, out, "should match the input message")
}

func TestRFC5425MessageLengthFramer(t *testing.T) {
	var (
		msg      = []byte("input message")
		expected = [][]byte{
			[]byte("13 "),
			msg,
		}
	)

	out := RFC5425MessageLengthFramer([]byte("input message"))
	assert.Equal(t, expected, out, "should prepend the input message length")
}
