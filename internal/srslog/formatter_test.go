package srslog

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFormatter(t *testing.T) {
	out := DefaultFormatter(time.Now(), LOG_ERR, "hostname", "tag", []byte("content"))
	expected := fmt.Sprintf("<%d> %s %s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.RFC3339), "hostname", "tag", os.Getpid(), "content")
	assert.Equal(t, expected, string(out))
}

func TestUnixFormatter(t *testing.T) {
	out := UnixFormatter(time.Now(), LOG_ERR, "hostname", "tag", []byte("content"))
	expected := fmt.Sprintf("<%d>%s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.Stamp), "tag", os.Getpid(), "content")
	assert.Equal(t, expected, string(out))
}

func TestRFC3164Formatter(t *testing.T) {
	out := RFC3164Formatter(time.Now(), LOG_ERR, "hostname", "tag", []byte("content"))
	expected := fmt.Sprintf("<%d>%s %s %s[%d]: %s",
		LOG_ERR, time.Now().Format(time.Stamp), "hostname", "tag", os.Getpid(), "content")
	assert.Equal(t, expected, string(out))
}

func TestRFC5424Formatter(t *testing.T) {
	out := RFC5424Formatter(time.Now(), LOG_ERR, "hostname", "tag", []byte("content"))
	expected := fmt.Sprintf("<%d>%d %s %s %s %d %s - %s",
		LOG_ERR, 1, time.Now().Format(time.RFC3339), "hostname", truncateStartStr(os.Args[0], appNameMaxLength),
		os.Getpid(), "tag", "content")
	assert.Equal(t, expected, string(out))
}

func TestTruncateStartStr(t *testing.T) {
	assert := assert.New(t)

	out := truncateStartStr("abcde", 3)
	assert.Equal("cde", out)

	out = truncateStartStr("abcde", 5)
	assert.Equal("abcde", out)
}
