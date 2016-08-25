package singularity

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	assert "github.com/stretchr/testify/require"
)

func TestCreateOutputDirs(t *testing.T) {
	err := CreateOutputDirs(TargetDir)
	assert.NoError(t, err)
}

func TestPlainFormatter(t *testing.T) {
	formatter := &plainFormatter{}

	var entry *log.Entry
	var bytes []byte
	var err error

	entry = &log.Entry{Message: "hello."}
	bytes, err = formatter.Format(entry)
	assert.NoError(t, err)
	assert.Equal(t, "hello.\n", string(bytes))

	entry = &log.Entry{Message: "debug.", Level: log.DebugLevel}
	bytes, err = formatter.Format(entry)
	assert.NoError(t, err)
	assert.Equal(t, "DEBUG: debug.\n", string(bytes))
}
