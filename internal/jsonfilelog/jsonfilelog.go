package jsonfilelog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog"
)

const (
	defaultLogDir = "/var/log/docker"
	logDirKey     = "json-file-log-dir"
)

// New returns the generic jsonfilelog log driver after parsing
// our custom options
func New(info logger.Info) (logger.Logger, error) {
	logDir := removeLogDirOption(info.Config)
	if logDir == "" {
		logDir = defaultLogDir
	}
	info.LogPath = filepath.Join(logDir, info.ContainerID)

	if err := os.MkdirAll(filepath.Dir(info.LogPath), 0755); err != nil {
		return nil, fmt.Errorf("error setting up logger dir: %v", err)
	}

	return jsonfilelog.New(info)
}

// ValidateLogOpt takes care of removing our custom jsonfilelog options
// before executing the real validator
func ValidateLogOpt(cfg map[string]string) error {
	logDir := removeLogDirOption(cfg)

	if err := jsonfilelog.ValidateLogOpt(cfg); err != nil {
		return err
	}

	if logDir != "" {
		cfg[logDirKey] = logDir
	}

	return nil
}

func removeLogDirOption(cfg map[string]string) string {
	if logDir, ok := cfg[logDirKey]; ok {
		delete(cfg, logDirKey)
		return logDir
	}
	return ""
}
