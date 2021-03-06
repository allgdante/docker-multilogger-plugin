package multilogger

import (
	"fmt"
	"strconv"

	"github.com/allgdante/docker-multilogger-plugin/internal/jsonfilelog"
	"github.com/allgdante/docker-multilogger-plugin/pkg/logassembler"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/go-units"
	"github.com/hashicorp/go-multierror"
)

// Driver name & available keys
const (
	DriverName = "multilogger"
	MaxSizeKey = DriverName + "-max-size"
)

const (
	minSize     = 16 * 1024
	defaultSize = 2 >> 20
)

type multiLogger struct {
	loggers   []logger.Logger
	assembler logassembler.Assembler
}

// Name implements the logger.Logger interface
func (ml *multiLogger) Name() string {
	return "multilogger"
}

// Log implements the logger.Logger interface
func (ml *multiLogger) Log(origmsg *logger.Message) (err error) {
	for _, cmsg := range ml.assembler.Assemble(origmsg) {
		for i, l := range ml.loggers {
			// Every builtin docker log driver resets the log message after writing it,
			// so we must clone the message before passing it to the native driver.
			// If it's the last driver, we will pass the original message
			var msg *logger.Message
			if i+1 != len(ml.loggers) {
				msg = logger.NewMessage()
				dumbCopyMessage(msg, cmsg)
			} else {
				msg = cmsg
			}
			if lerr := l.Log(msg); lerr != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", l.Name(), lerr))
				fmt.Println(lerr.Error())
			}
		}
	}

	return
}

// Close implements the logger.Logger interface
func (ml *multiLogger) Close() (err error) {
	for _, l := range ml.loggers {
		if lerr := l.Close(); lerr != nil {
			err = multierror.Append(err, fmt.Errorf("%s: %w", l.Name(), lerr))
		}
	}
	return
}

// ReadLogs implements the LogReader interface
// Note that it returns the first listed logger that implements the LogReader
// interface or nil if we can't find one
func (ml *multiLogger) ReadLogs(config logger.ReadConfig) *logger.LogWatcher {
	for _, l := range ml.loggers {
		if lr, ok := l.(logger.LogReader); ok {
			return lr.ReadLogs(config)
		}
	}
	return nil
}

// Logger creates a logger.Logger that writes the log messages to the
// provided loggers, similar to the Unix tee(1) command.
//
// Each log is written to each listed logger, one at a time.
// If a listed logger returns an error, the Log operation continue down the list.
func Logger(size int64, loggers ...logger.Logger) logger.Logger {
	allLoggers := make([]logger.Logger, 0, len(loggers))
	for _, l := range loggers {
		if ml, ok := l.(*multiLogger); ok {
			allLoggers = append(allLoggers, ml.loggers...)
		} else {
			allLoggers = append(allLoggers, l)
		}
	}
	return &multiLogger{
		loggers:   allLoggers,
		assembler: logassembler.New(size),
	}
}

// Validator returns a logger.LogOptValidator which will validate the config
// for all the enabled logging drivers
func Validator(blueprints []Blueprint) logger.LogOptValidator {
	return func(cfg map[string]string) (err error) {
		if _, serr := parseMaxSize(cfg[MaxSizeKey]); serr != nil {
			err = multierror.Append(err, serr)
		}

		for _, blp := range blueprints {
			if parseLogOptBoolean(cfg, blp.EnabledKey()) {
				logcfg, serr := blp.Config(cfg)
				if serr != nil {
					err = multierror.Append(err, fmt.Errorf("%s: %w", blp.Name, serr))
					continue
				}

				if lerr := blp.Validate(logcfg); lerr != nil {
					err = multierror.Append(err, fmt.Errorf("%s: %w", blp.Name, lerr))
					continue
				}
			}
		}

		return
	}
}

// Creator returns a logger.Creator which will take care of create
// a logger.Logger which will have enabled all the requested logging drivers
func Creator(blueprints []Blueprint) logger.Creator {
	return func(info logger.Info) (logger.Logger, error) {
		var (
			loggers []logger.Logger
			err     error
		)

		size, serr := parseMaxSize(info.Config[MaxSizeKey])
		if serr != nil {
			err = multierror.Append(err, serr)
		}

		for _, blp := range blueprints {
			if parseLogOptBoolean(info.Config, blp.EnabledKey()) {
				logcfg, serr := blp.Config(info.Config)
				if serr != nil {
					err = multierror.Append(err, fmt.Errorf("%s: %w", blp.Name, serr))
					continue
				}

				newinfo := info
				newinfo.Config = logcfg
				logdrv, lerr := blp.Create(newinfo)
				if lerr != nil {
					err = multierror.Append(err, fmt.Errorf("%s: %w", blp.Name, lerr))
					continue
				}
				loggers = append(loggers, logdrv)
			}
		}

		if err != nil {
			return nil, err
		}

		// If there is no logger enabled, we always add the jsonfile driver by default.
		if len(loggers) == 0 {
			logger, err := jsonfilelog.New(info)
			if err != nil {
				return nil, fmt.Errorf("failure while adding default jsonfile driver: %v", err)
			}
			loggers = append(loggers, logger)
		}

		return Logger(size, loggers...), nil
	}
}

// dumbCopyMessage is a bit of a fake copy but avoids extra allocations which
// are not necessary for this use case.
// XXX: extracted from https://github.com/moby/moby/pull/40543
func dumbCopyMessage(dst, src *logger.Message) {
	dst.Source = src.Source
	dst.Timestamp = src.Timestamp
	dst.PLogMetaData = src.PLogMetaData
	dst.Err = src.Err
	dst.Attrs = src.Attrs
	dst.Line = append(dst.Line[:0], src.Line...)
}

// parseMaxSize parses the size for the underlying logassembler
func parseMaxSize(humanSize string) (int64, error) {
	if humanSize == "" {
		return defaultSize, nil
	}

	if size, err := units.FromHumanSize(humanSize); err != nil {
		return 0, fmt.Errorf("invalid value for max size: %v", err)
	} else if size < minSize {
		return 0, fmt.Errorf("max size must be at least 16KB")
	} else {
		return size, nil
	}
}

// parseLogOptBoolean parses an option as a boolean value
func parseLogOptBoolean(config map[string]string, logOptKey string) bool {
	if input, exists := config[logOptKey]; exists {
		// If we can't parse the provided value as a bool, we just return
		// true here and let the validator fail
		inputValue, err := strconv.ParseBool(input)
		if err != nil {
			return true
		}
		return inputValue
	}
	return false
}
