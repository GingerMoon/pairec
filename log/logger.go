package log

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ENV_APOLLO_URL   = "apollo_apollo_meta"
	ENV_PAIREC_DEBUG = "PAIREC_DEBUG"
)

var (
	prodCfg *zap.Config
	debug   bool
	letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func SetConfig(logConfig *zap.Config) {
	prodCfg = logConfig
	isDebug := os.Getenv(ENV_PAIREC_DEBUG)
	if isDebug != "" {
		debug = true
	} else {
		debug = false
	}
}

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewLogger() *zap.Logger {
	env := os.Getenv(ENV_APOLLO_URL)
	if env != "" {
		// required by devops team
		path := prodCfg.OutputPaths[0]
		segments := strings.Split(path, ".log")
		path = segments[0] + "_" + randStr(10) + ".log"
		for i := 1; i < len(segments); i++ {
			path += segments[i]
		}

		logPath := filepath.Dir(path)
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			panic("failed to make log dir " + logPath)
		}

		prodCfg.OutputPaths[0] = path
	}

	prodCfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	logger := zap.Must(prodCfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	return logger
}

func Logger(ctx context.Context) *zap.Logger {
	return ctx.Value("logger").(*zap.Logger)
}
