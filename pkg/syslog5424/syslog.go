// Package syslog5424 provides an specialized rfc5424 log driver for forwarding server logs to syslog endpoints.
package syslog5424

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	syslog "github.com/allgdante/srslog"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/daemon/logger/templates"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/docker/go-connections/tlsconfig"
)

// Driver name & available keys
const (
	DriverName       = "syslog5424"
	AddressKey       = DriverName + "-address"
	FacilityKey      = DriverName + "-facility"
	TimeFormatKey    = DriverName + "-time-format"
	TLSCACertKey     = DriverName + "-tls-ca-cert"
	TLSCertKey       = DriverName + "-tls-cert"
	TLSKeyKey        = DriverName + "-tls-key"
	TLSSkipVerifyKey = DriverName + "-tls-skip-verify"
	HostnameKey      = DriverName + "-hostname"
	MSGIDKey         = DriverName + "-msgid"
	DisableFramerKey = DriverName + "-disable-framer"
	EnvKey           = "env"
	EnvRegexKey      = "env-regex"
	LabelsKey        = "labels"
	LabelsRegexKey   = "labels-regex"
	TagKey           = "tag"
)

// Available default time formats
const (
	RFC3339TimeFormat      = "rfc3339"
	RFC3399MicroTimeFormat = "rfc3339micro"
)

const (
	secureProto = "tcp+tls"
)

var facilities = map[string]syslog.Priority{
	"kern":     syslog.LOG_KERN,
	"user":     syslog.LOG_USER,
	"mail":     syslog.LOG_MAIL,
	"daemon":   syslog.LOG_DAEMON,
	"auth":     syslog.LOG_AUTH,
	"syslog":   syslog.LOG_SYSLOG,
	"lpr":      syslog.LOG_LPR,
	"news":     syslog.LOG_NEWS,
	"uucp":     syslog.LOG_UUCP,
	"cron":     syslog.LOG_CRON,
	"authpriv": syslog.LOG_AUTHPRIV,
	"ftp":      syslog.LOG_FTP,
	"local0":   syslog.LOG_LOCAL0,
	"local1":   syslog.LOG_LOCAL1,
	"local2":   syslog.LOG_LOCAL2,
	"local3":   syslog.LOG_LOCAL3,
	"local4":   syslog.LOG_LOCAL4,
	"local5":   syslog.LOG_LOCAL5,
	"local6":   syslog.LOG_LOCAL6,
	"local7":   syslog.LOG_LOCAL7,
}

type syslogger struct {
	writer *syslog.Writer
}

// New creates a syslog logger using the configuration passed in on
// the context.
func New(info logger.Info) (logger.Logger, error) {
	tag, err := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if err != nil {
		return nil, err
	}

	msgid, err := parseOptAsTemplate(info, MSGIDKey)
	if err != nil {
		return nil, err
	}
	if msgid == "" {
		msgid = tag
	}

	proto, address, err := parseAddress(info.Config[AddressKey])
	if err != nil {
		return nil, err
	}

	facility, err := parseFacility(info.Config[FacilityKey])
	if err != nil {
		return nil, err
	}

	timeFormat, err := parseTimeFormat(info.Config[TimeFormatKey])
	if err != nil {
		return nil, err
	}

	extra, err := info.ExtraAttributes(nil)
	if err != nil {
		return nil, errdefs.InvalidParameter(err)
	}

	hostname, err := parseOptAsTemplate(info, HostnameKey)
	if err != nil {
		return nil, err
	}

	var disableFramer bool
	if df, ok := info.Config[DisableFramerKey]; ok {
		if disableFramer, err = strconv.ParseBool(df); err != nil {
			return nil, errdefs.InvalidParameter(err)
		}
	}

	var log *syslog.Writer
	if proto == secureProto {
		tlsConfig, tlsErr := parseTLSConfig(info.Config)
		if tlsErr != nil {
			return nil, tlsErr
		}
		log, err = syslog.DialWithTLSConfig(proto, address, facility, tag, tlsConfig)
	} else {
		log, err = syslog.Dial(proto, address, facility, tag)
	}
	if err != nil {
		return nil, err
	}

	if hostname != "" {
		log.SetHostname(hostname)
	}

	log.SetFormatter(rfc5424Formatter(timeFormat, msgid, extra))
	if !disableFramer {
		log.SetFramer(syslog.RFC5425MessageLengthFramer)
	}

	return &syslogger{
		writer: log,
	}, nil
}

func (s *syslogger) Log(msg *logger.Message) error {
	if len(msg.Line) == 0 {
		logger.PutMessage(msg)
		return nil
	}

	if msg.Source == "stderr" {
		_, err := s.writer.WriteWithTimestampAndPriority(msg.Timestamp, syslog.LOG_ERR, msg.Line)
		logger.PutMessage(msg)
		return err
	}
	_, err := s.writer.WriteWithTimestampAndPriority(msg.Timestamp, syslog.LOG_INFO, msg.Line)
	logger.PutMessage(msg)
	return err
}

func (s *syslogger) Close() error {
	return s.writer.Close()
}

func (s *syslogger) Name() string {
	return DriverName
}

// ValidateLogOpt looks for syslog specific log options
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case EnvKey:
		case EnvRegexKey:
		case LabelsKey:
		case LabelsRegexKey:
		case AddressKey:
		case FacilityKey:
		case TimeFormatKey:
		case TLSCACertKey:
		case TLSCertKey:
		case TLSKeyKey:
		case TLSSkipVerifyKey:
		case HostnameKey:
		case MSGIDKey:
		case DisableFramerKey:
		case TagKey:
		default:
			return fmt.Errorf("unknown log opt '%s' for syslog5424 log driver", key)
		}
	}
	if _, _, err := parseAddress(cfg[AddressKey]); err != nil {
		return err
	}
	if _, err := parseFacility(cfg[FacilityKey]); err != nil {
		return err
	}
	if _, err := parseTimeFormat(cfg[TimeFormatKey]); err != nil {
		return err
	}
	return nil
}

func parseAddress(address string) (string, string, error) {
	if address == "" {
		return "", "", nil
	}
	if !urlutil.IsTransportURL(address) {
		return "", "", fmt.Errorf("address should be in form proto://address, got %v", address)
	}
	url, err := url.Parse(address)
	if err != nil {
		return "", "", err
	}

	// unix and unixgram socket validation
	if url.Scheme == "unix" || url.Scheme == "unixgram" {
		if _, err := os.Stat(url.Path); err != nil {
			return "", "", err
		}
		return url.Scheme, url.Path, nil
	}

	// here we process tcp|udp
	host := url.Host
	if _, _, err := net.SplitHostPort(host); err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			return "", "", err
		}
		host = host + ":514"
	}

	return url.Scheme, host, nil
}

func parseFacility(facility string) (syslog.Priority, error) {
	if facility == "" {
		return syslog.LOG_DAEMON, nil
	}

	if syslogFacility, valid := facilities[facility]; valid {
		return syslogFacility, nil
	}

	fInt, err := strconv.Atoi(facility)
	if err == nil && 0 <= fInt && fInt <= 23 {
		return syslog.Priority(fInt << 3), nil
	}

	return syslog.Priority(0), errors.New("invalid syslog facility")
}

func parseTimeFormat(timeFormat string) (string, error) {
	switch timeFormat {
	case "", RFC3339TimeFormat:
		return time.RFC3339, nil
	case RFC3399MicroTimeFormat:
		return "2006-01-02T15:04:05.000000Z07:00", nil
	default:
		return "", errors.New("invalid syslog time format")
	}
}

func parseTLSConfig(cfg map[string]string) (*tls.Config, error) {
	_, skipVerify := cfg[TLSSkipVerifyKey]

	opts := tlsconfig.Options{
		CAFile:             cfg[TLSCACertKey],
		CertFile:           cfg[TLSCertKey],
		KeyFile:            cfg[TLSKeyKey],
		InsecureSkipVerify: skipVerify,
	}

	return tlsconfig.Client(opts)
}

func parseOptAsTemplate(info logger.Info, key string) (string, error) {
	optTemplate := info.Config[key]
	if optTemplate == "" {
		return "", nil
	}

	tmpl, err := templates.NewParse("opt", optTemplate)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, &info); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// rfc5424Formatter provides an RFC 5424 compliant message formatter.
func rfc5424Formatter(timeFormat, msgid string, extra map[string]string) syslog.Formatter {
	pid := os.Getpid()
	return func(timestamp time.Time, p syslog.Priority, hostname, tag, content string) string {
		var b strings.Builder
		fmt.Fprintf(&b, "<%d>1 %s %s %s %d %s ",
			p, timestamp.Format(timeFormat), hostname, tag, pid, msgid)
		if len(extra) > 0 {
			b.WriteString("[docker@3071")
			for k, v := range extra {
				fmt.Fprintf(&b, " %s=\"%s\"", k, escapeSDParam(v))
			}
			b.WriteString("]")
		} else {
			b.WriteString("-")
		}
		b.WriteString(" ")
		b.WriteString(content)
		return b.String()
	}
}

// escapeSDParam escapes the sd-param according to rfc5424.
// Taken from https://github.com/crewjam/rfc5424/blob/master/marshal.go#L39
func escapeSDParam(s string) string {
	escapeCount := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\', '"', ']':
			escapeCount++
		}
	}
	if escapeCount == 0 {
		return s
	}

	t := make([]byte, len(s)+escapeCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\', '"', ']':
			t[j] = '\\'
			t[j+1] = c
			j += 2
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
