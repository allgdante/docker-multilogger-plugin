package multilogger

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorEmpty(t *testing.T) {
	emptyConfig := make(map[string]string)

	err := Validator(DefaultBlueprints)(emptyConfig)
	assert.Nil(t, err)
}

func TestAddedJSONFileDriverIfNeeded(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
	)

	// Note that we need to create a temporal directory as the default
	// directory for the jsonfile log driver.
	// Otherwise, the driver will try to use /var/log/docker directory
	// and it could be possible that we don't have permission on this dir
	logDir, err := ioutil.TempDir("", "xxxx")
	require.Nil(err)
	defer os.RemoveAll(logDir)

	var info logger.Info
	info.ContainerID = "7f0ebc7d0b9a756b16dc6c1c4df31050e6a76fc7b013761df97b79c07bc0336e"
	info.Config = map[string]string{
		"json-file-log-dir": logDir,
	}

	logger, err := Creator(DefaultBlueprints)(info)
	require.Nil(err)

	ml, ok := logger.(*multiLogger)
	require.True(ok)
	assert.Equal(1, len(ml.loggers))
	require.Equal("json-file", ml.loggers[0].Name())
}
