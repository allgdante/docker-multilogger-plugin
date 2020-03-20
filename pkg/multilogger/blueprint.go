package multilogger

import (
	"strings"

	"github.com/docker/docker/daemon/logger"
)

// Blueprint represents the specification that must be met by a log driver in order
// to be used by a multilogger
type Blueprint struct {
	Name     string
	Options  []string
	Create   logger.Creator
	Validate logger.LogOptValidator
}

// EnabledKey returns the config name option to check if the driver is enabled
func (b Blueprint) EnabledKey() string {
	return b.Name + "-enabled"
}

// Config returns the driver-specific config extracted from the global config after validating it
func (b Blueprint) Config(globalcfg map[string]string) (drivercfg map[string]string, err error) {
	drivercfg = make(map[string]string)

	for _, opt := range b.Options {
		if v, ok := globalcfg[opt]; ok {
			switch trimopt := strings.TrimPrefix(opt, b.Name+"-"); trimopt {
			case "labels", "labels-regex", "env", "env-regex", "tag":
				drivercfg[trimopt] = v
			default:
				drivercfg[opt] = v
			}
		}
	}

	err = b.Validate(drivercfg)
	return
}
