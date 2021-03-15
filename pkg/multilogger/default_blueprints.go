package multilogger

import (
	"github.com/allgdante/docker-multilogger-plugin/internal/jsonfilelog"
	"github.com/allgdante/docker-multilogger-plugin/pkg/syslog5424"

	"github.com/docker/docker/daemon/logger/awslogs"
	"github.com/docker/docker/daemon/logger/fluentd"
	"github.com/docker/docker/daemon/logger/gcplogs"
	"github.com/docker/docker/daemon/logger/gelf"
	"github.com/docker/docker/daemon/logger/journald"
	"github.com/docker/docker/daemon/logger/logentries"
	"github.com/docker/docker/daemon/logger/splunk"
	"github.com/docker/docker/daemon/logger/syslog"
)

var (
	// JSONFileLogBlueprint is the blueprint for our customized jsonfilelog driver
	JSONFileLogBlueprint = Blueprint{
		"json-file",
		[]string{
			"max-size",
			"max-file",
			"labels",
			"labels-regex",
			"env",
			"env-regex",
			"compress",
			"tag",
			"json-file-log-dir",
		},
		jsonfilelog.New,
		jsonfilelog.ValidateLogOpt,
	}

	// AWSLogsBlueprint is the blueprint for the original awslogs docker driver
	AWSLogsBlueprint = Blueprint{
		"awslogs",
		[]string{
			"awslogs-region",
			"awslogs-endpoint",
			"awslogs-group",
			"awslogs-create-group",
			"awslogs-datetime-format",
			"awslogs-multiline-pattern",
			"awslogs-credentials-endpoint",
			"awslogs-force-flush-interval-seconds",
			"awslogs-max-buffered-events",
			"awslogs-tag",
		},
		awslogs.New,
		awslogs.ValidateLogOpt,
	}

	// FluentdBlueprint is the blueprint for the original fluentd docker driver
	FluentdBlueprint = Blueprint{
		"fluentd",
		[]string{
			"fluentd-address",
			"fluentd-async",
			"fluentd-async-connect",
			"fluentd-buffer-limit",
			"fluentd-max-retries",
			"fluentd-request-ack",
			"fluentd-retry-wait",
			"fluentd-sub-second-precision",
			"fluentd-labels",
			"fluentd-labels-regex",
			"fluentd-env",
			"fluentd-env-regex",
			"fluentd-tag",
		},
		fluentd.New,
		fluentd.ValidateLogOpt,
	}

	// GCPLogsBlueprint is the blueprint for the original gcplogs docker driver
	GCPLogsBlueprint = Blueprint{
		"gcp",
		[]string{
			"gcp-project",
			"gcp-log-cmd",
			"gcp-meta-zone",
			"gcp-meta-name",
			"gcp-meta-id",
			"gcp-labels",
			"gcp-labels-regex",
			"gcp-env",
			"gcp-env-regex",
		},
		gcplogs.New,
		gcplogs.ValidateLogOpts,
	}

	// GelfBlueprint is the blueprint for the original gelf docker driver
	GelfBlueprint = Blueprint{
		"gelf",
		[]string{
			"gelf-address",
			"gelf-compression-level",
			"gelf-compression-type",
			"gelf-tcp-max-reconnect",
			"gelf-labels",
			"gelf-labels-regex",
			"gelf-env",
			"gelf-env-regex",
			"gelf-tag",
		},
		gelf.New,
		gelf.ValidateLogOpt,
	}

	// JournaldBlueprint is the blueprint for the original journald docker driver
	JournaldBlueprint = Blueprint{
		"journald",
		[]string{},
		journald.New,
		func(_ map[string]string) error {
			return nil
		},
	}

	// LogentriesBlueprint is the blueprint for the original logentries docker driver
	LogentriesBlueprint = Blueprint{
		"logentries",
		[]string{
			"logentries-token",
			"line-only",
			"logentries-labels",
			"logentries-labels-regex",
			"logentries-env",
			"logentries-env-regex",
			"logentries-tag",
		},
		logentries.New,
		logentries.ValidateLogOpt,
	}

	// SplunkBlueprint is the blueprint for the original splunk docker driver
	SplunkBlueprint = Blueprint{
		"splunk",
		[]string{
			"splunk-url",
			"splunk-token",
			"splunk-source",
			"splunk-sourcetype",
			"splunk-index",
			"splunk-capath",
			"splunk-caname",
			"splunk-insecureskipverify",
			"splunk-format",
			"splunk-verify-connection",
			"splunk-gzip",
			"splunk-gzip-devel",
			"splunk-index-acknowledgment",
			"splunk-labels",
			"splunk-labels-regex",
			"splunk-env",
			"splunk-env-regex",
			"splunk-tag",
		},
		splunk.New,
		splunk.ValidateLogOpt,
	}

	// SyslogBlueprint is the blueprint for the original syslog docker driver
	SyslogBlueprint = Blueprint{
		"syslog",
		[]string{
			"syslog-address",
			"syslog-facility",
			"syslog-tls-ca-cert",
			"syslog-tls-cert",
			"syslog-tls-key",
			"syslog-tls-skip-verify",
			"syslog-format",
			"syslog-labels",
			"syslog-labels-regex",
			"syslog-env",
			"syslog-env-regex",
			"syslog-tag",
		},
		syslog.New,
		syslog.ValidateLogOpt,
	}

	// Syslog5424Blueprint is the blueprint for our own syslog driver
	Syslog5424Blueprint = Blueprint{
		syslog5424.DriverName,
		[]string{
			syslog5424.AddressKey,
			syslog5424.FacilityKey,
			syslog5424.TimeFormatKey,
			syslog5424.TLSCACertKey,
			syslog5424.TLSCertKey,
			syslog5424.TLSKeyKey,
			syslog5424.TLSSkipVerifyKey,
			syslog5424.HostnameKey,
			syslog5424.MSGIDKey,
			syslog5424.DisableFramerKey,
			syslog5424.DriverName + "-" + syslog5424.LabelsKey,
			syslog5424.DriverName + "-" + syslog5424.LabelsRegexKey,
			syslog5424.DriverName + "-" + syslog5424.EnvKey,
			syslog5424.DriverName + "-" + syslog5424.EnvRegexKey,
			syslog5424.DriverName + "-" + syslog5424.TagKey,
		},
		syslog5424.New,
		syslog5424.ValidateLogOpt,
	}

	// DefaultBlueprints represents the builtin docker log drivers with our
	// custom syslog5424 driver
	DefaultBlueprints = []Blueprint{
		JSONFileLogBlueprint,
		AWSLogsBlueprint,
		GCPLogsBlueprint,
		GelfBlueprint,
		JournaldBlueprint,
		LogentriesBlueprint,
		SplunkBlueprint,
		SyslogBlueprint,
		Syslog5424Blueprint,
	}
)
