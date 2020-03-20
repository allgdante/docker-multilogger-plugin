package multilogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlueprint(t *testing.T) {
	var assert = assert.New(t)

	globalcfg := map[string]string{
		"labels":              "foo",
		"max-size":            "20",
		"awslogs-group":       "foo",
		"gcp-project":         "foobar",
		"gcp-labels":          "bar",
		"gelf-address":        "tcp://127.0.0.1:1234",
		"gelf-labels":         "foobar",
		"logentries-token":    "xxxxxxxx",
		"syslog-facility":     "local1",
		"syslog-labels":       "foobar",
		"syslog5424-facility": "local1",
		"syslog5424-address":  "tcp+tls://127.0.0.1:8443",
	}

	t.Run("JSONFile Blueprint", func(t *testing.T) {
		cfg, err := JSONFileLogBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"labels":   "foo",
			"max-size": "20",
		}, cfg)
	})

	t.Run("AWSLogs Blueprint Blueprint", func(t *testing.T) {
		cfg, err := AWSLogsBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"awslogs-group": "foo",
		}, cfg)
	})

	t.Run("GCPLogs Blueprint Blueprint", func(t *testing.T) {
		cfg, err := GCPLogsBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"gcp-project": "foobar",
			"labels":      "bar",
		}, cfg)
	})

	t.Run("Gelf Blueprint Blueprint", func(t *testing.T) {
		cfg, err := GelfBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"gelf-address": "tcp://127.0.0.1:1234",
			"labels":       "foobar",
		}, cfg)
	})

	t.Run("Journald Blueprint Blueprint", func(t *testing.T) {
		cfg, err := JournaldBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{}, cfg)
	})

	t.Run("Logentries Blueprint Blueprint", func(t *testing.T) {
		cfg, err := LogentriesBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"logentries-token": "xxxxxxxx",
		}, cfg)
	})

	t.Run("Splunk Blueprint Blueprint", func(t *testing.T) {
		cfg, err := SplunkBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{}, cfg)
	})

	t.Run("Syslog Blueprint Blueprint", func(t *testing.T) {
		cfg, err := SyslogBlueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"syslog-facility": "local1",
			"labels":          "foobar",
		}, cfg)
	})

	t.Run("Syslog5424 Blueprint Blueprint", func(t *testing.T) {
		cfg, err := Syslog5424Blueprint.Config(globalcfg)
		assert.Nil(err)
		assert.Equal(map[string]string{
			"syslog5424-facility": "local1",
			"syslog5424-address":  "tcp+tls://127.0.0.1:8443",
		}, cfg)
	})
}
