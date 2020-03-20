package jsonfilelog

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLogOptEmpty(t *testing.T) {
	emptyConfig := make(map[string]string)
	err := ValidateLogOpt(emptyConfig)

	require.Nil(t, err)
}

func TestValidateLogDir(t *testing.T) {
	logDir, err := ioutil.TempDir("", "xxx")
	require.Nil(t, err)
	defer os.RemoveAll(logDir)

	err = ValidateLogOpt(map[string]string{
		"json-file-log-dir": logDir,
	})
	require.Nil(t, err)
}
