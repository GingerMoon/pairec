package log

import (
	"fmt"

	"go.uber.org/zap"
)

var (
	zapLogger *zap.Logger
)

func SetLogger() {
	zapLogger = NewLogger()
}

func Debug(msg string) {
	zapLogger.Debug(msg)
}
func Info(msg string) {
	zapLogger.Info(msg)
}
func Warning(msg string) {
	zapLogger.Warn(msg)
}
func Error(msg string) {
	zapLogger.Error(msg)
}
func Flush() {

}

type ABTestLogger struct{}

func (l ABTestLogger) Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}
func (l ABTestLogger) Errorf(msg string, args ...interface{}) {
	Error(fmt.Sprintf(msg, args...))
}

type FeatureStoreLogger struct{}

func (l FeatureStoreLogger) Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}
func (l FeatureStoreLogger) Errorf(msg string, args ...interface{}) {
	Error(fmt.Sprintf(msg, args...))
}
