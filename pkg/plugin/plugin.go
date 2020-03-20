// Package plugin provides several helpers to build a docker log driver
package plugin

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"syscall"

	"github.com/containerd/fifo"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	protoio "github.com/gogo/protobuf/io"
	"github.com/sirupsen/logrus"
)

// Plugin represents the minimal functionality that must be implemented
// by a logging plugin
type Plugin interface {
	StartLogging(file string, info logger.Info) error
	StopLogging(file string) error
	Capabilities() logger.Capability
	ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error)
}

// loggingPlugin is a base implementation for a Plugin
type loggingPlugin struct {
	validator logger.LogOptValidator
	creator   logger.Creator
	logs      map[string]*adapter
	mu        sync.Mutex
}

// New returns a new Plugin
func New(validator logger.LogOptValidator, creator logger.Creator) Plugin {
	return &loggingPlugin{
		validator: validator,
		creator:   creator,
		logs:      make(map[string]*adapter),
	}
}

// StartLogging implements the Plugin interface
func (p *loggingPlugin) StartLogging(file string, info logger.Info) error {
	p.mu.Lock()
	if _, exists := p.logs[file]; exists {
		p.mu.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	}
	p.mu.Unlock()

	if err := p.validator(info.Config); err != nil {
		return fmt.Errorf("error validating logging drivers: %w", err)
	}

	logger, err := p.creator(info)
	if err != nil {
		return fmt.Errorf("error loading logging drivers: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"id":   info.ContainerID,
		"file": file,
	}).Debugf("Start logging")
	logFifo, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return fmt.Errorf("error opening logger fifo: %q: %w", file, err)
	}

	p.mu.Lock()
	a := &adapter{
		id:     info.ContainerID,
		logger: logger,
		stream: logFifo,
	}
	p.logs[file] = a
	p.logs[info.ContainerID] = a
	p.mu.Unlock()

	go a.Start()
	return nil
}

// StopLogging implements the Plugin interface
func (p *loggingPlugin) StopLogging(file string) error {
	logrus.WithField("file", file).Debugf("Stop logging")

	p.mu.Lock()
	a, ok := p.logs[file]
	if ok {
		//	Stop logger and remove reference from driver state
		a.logger.Close()
		delete(p.logs, file)
		delete(p.logs, a.id)
	}
	p.mu.Unlock()

	return nil
}

// Capabilities implements the Plugin interface
func (p *loggingPlugin) Capabilities() logger.Capability {
	return logger.Capability{ReadLogs: true}
}

// ReadLogs implements the Plugin interface
func (p *loggingPlugin) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	p.mu.Lock()
	a, exists := p.logs[info.ContainerID]
	p.mu.Unlock()
	if !exists {
		return nil, fmt.Errorf("logger does not exist for %s", info.ContainerID)
	}

	r, w := io.Pipe()
	lr, ok := a.logger.(logger.LogReader)
	if !ok {
		return nil, logger.ErrReadLogsNotSupported{}
	}

	watcher := lr.ReadLogs(config)
	if watcher == nil {
		return nil, logger.ErrReadLogsNotSupported{}
	}

	go func() {
		enc := protoio.NewUint32DelimitedWriter(w, binary.BigEndian)
		defer enc.Close()
		defer watcher.ConsumerGone()

		var buf logdriver.LogEntry
		for {
			select {
			case msg, ok := <-watcher.Msg:
				if !ok {
					w.Close()
					return
				}

				buf.Line = msg.Line
				buf.Partial = msg.PLogMetaData != nil
				buf.TimeNano = msg.Timestamp.UnixNano()
				buf.Source = msg.Source

				if err := enc.WriteMsg(&buf); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			case err := <-watcher.Err:
				_ = w.CloseWithError(err)
				return
			}

			buf.Reset()
		}
	}()

	return r, nil
}
