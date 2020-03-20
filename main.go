package main // import "github.com/allgdante/docker-multilogger-plugin"

import (
	"fmt"
	"os"

	"github.com/allgdante/docker-multilogger-plugin/pkg/multilogger"
	"github.com/allgdante/docker-multilogger-plugin/pkg/plugin"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/sirupsen/logrus"
)

const socketAddress = "/run/docker/plugins/multilogger.sock"

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	levelVal := os.Getenv("LOG_LEVEL")
	if levelVal == "" {
		levelVal = "info"
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		fmt.Fprintln(os.Stderr, "invalid log level: ", levelVal)
		os.Exit(1)
	}

	var (
		handler       = sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
		blueprints    = multilogger.DefaultBlueprints
		pluginHandler = &plugin.HTTPHandler{
			Plugin: plugin.New(
				multilogger.Validator(blueprints),
				multilogger.Creator(blueprints),
			),
		}
	)

	pluginHandler.Initialize(&handler)
	if err := handler.ServeUnix(socketAddress, 0); err != nil {
		logrus.Fatal(err)
	}
}
