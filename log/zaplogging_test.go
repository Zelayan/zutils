package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_zapLogger_Info(t *testing.T) {
	check := assert.New(t)
	logger, err := NewZapLogger(Configuration{LogLevel: "info"})

	check.Nil(err)
	logger.Info("test")
}

func TestZapLogger_InfoF(t *testing.T) {
	check := assert.New(t)
	logger, err := NewZapLogger(Configuration{LogLevel: "info"})

	check.Nil(err)
	logger.InfoF("test:%s", "err")
}
