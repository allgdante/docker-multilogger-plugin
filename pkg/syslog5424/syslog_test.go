package syslog5424

import (
	"net"
	"testing"

	"github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLogOptEmpty(t *testing.T) {
	emptyConfig := make(map[string]string)

	err := ValidateLogOpt(emptyConfig)
	assert.Nil(t, err, "Failed to parse empty config")
}

func TestValidateSyslogAddress(t *testing.T) {
	assert := assert.New(t)

	err := ValidateLogOpt(map[string]string{
		AddressKey: "this is not an uri",
	})
	assert.NotNil(err, "Expected error with invalid uri")

	// File exists
	err = ValidateLogOpt(map[string]string{
		AddressKey: "unix:///",
	})
	assert.Nil(err)

	// File does not exist
	err = ValidateLogOpt(map[string]string{
		AddressKey: "unix:///does_not_exist",
	})
	assert.NotNil(err, "Expected error when address is non existing file")

	// accepts udp and tcp URIs
	err = ValidateLogOpt(map[string]string{
		AddressKey: "udp://1.2.3.4",
	})
	assert.Nil(err)

	err = ValidateLogOpt(map[string]string{
		AddressKey: "tcp://1.2.3.4",
	})
	assert.Nil(err)
}

func TestParseAddressDefaultPort(t *testing.T) {
	_, address, err := parseAddress("tcp://1.2.3.4")
	require.Nil(t, err)

	_, port, _ := net.SplitHostPort(address)
	assert.Equal(t, "514", port, "Expected to default to port 514. It used port ", port)
}

func TestValidateSyslogFacility(t *testing.T) {
	err := ValidateLogOpt(map[string]string{
		FacilityKey: "Invalid facility",
	})
	assert.NotNil(t, err, "Expected error if facility level is invalid")
}

func TestValidateLogOptSyslogTimeFormat(t *testing.T) {
	err := ValidateLogOpt(map[string]string{
		TimeFormatKey: "Invalid format",
	})
	assert.NotNil(t, err, "Expected error if time format is invalid")
}

func TestParseOptAsTemplate(t *testing.T) {
	assert := assert.New(t)

	info := logger.Info{
		Config: map[string]string{
			HostnameKey: "sample",
			MSGIDKey:    `{{index .ContainerLabels "foo"}}`,
		},
		ContainerLabels: map[string]string{
			"foo": "bar",
		},
	}

	t.Run("check literal values", func(t *testing.T) {
		value, err := parseOptAsTemplate(info, HostnameKey)
		assert.Nil(err)
		assert.Equal("sample", value)
	})

	t.Run("check template values", func(t *testing.T) {
		value, err := parseOptAsTemplate(info, MSGIDKey)
		assert.Nil(err)
		assert.Equal("bar", value)

	})
}

func TestValidateLogOpt(t *testing.T) {
	err := ValidateLogOpt(map[string]string{
		EnvKey:           "http://127.0.0.1",
		EnvRegexKey:      "abc",
		LabelsKey:        "labelA",
		LabelsRegexKey:   "def",
		AddressKey:       "udp://1.2.3.4:1111",
		FacilityKey:      "daemon",
		TLSCACertKey:     "/etc/ca-certificates/custom/ca.pem",
		TLSCertKey:       "/etc/ca-certificates/custom/cert.pem",
		TLSKeyKey:        "/etc/ca-certificates/custom/key.pem",
		TLSSkipVerifyKey: "true",
		TagKey:           "true",
		TimeFormatKey:    "rfc3339",
	})
	require.Nil(t, err)

	err = ValidateLogOpt(map[string]string{
		"not-supported-option": "a",
	})
	assert.NotNil(t, err, "Expecting error on unsupported options")
}
