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
	ENV_APOLLO_URL = "apollo_apollo_meta"
)

var (
	cfg     *zap.Config
	letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func SetConfig(logConfig *zap.Config) {
	cfg = logConfig
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
		path := cfg.OutputPaths[0]
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

		cfg.OutputPaths[0] = path
	}

	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	return logger
}

func Logger(ctx context.Context) *zap.Logger {
	return ctx.Value("logger").(*zap.Logger)
}
