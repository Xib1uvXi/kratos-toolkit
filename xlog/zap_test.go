package xlog

import (
	"go.uber.org/zap"
	"testing"
)

func TestInitLogger(t *testing.T) {
	log := InitLogger(zap.InfoLevel)

	log.Debug("debug")
	log.Error("error")

	tmp := t.TempDir()
	log2 := InitLoggerToFile(zap.InfoLevel, tmp, "test.log", false)

	log2.Debug("debug")
	log2.Error("error")

}
