/**
 * Logger Implementation
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
  log.SetOutput(os.Stdout)
  log.SetFormatter(&logrus.JSONFormatter{})
  log.SetLevel(logrus.InfoLevel)
}

func Debug(msg string) {
  log.Debug(msg)
}

func Info(msg string) {
  log.Info(msg)
}

func Warn(msg string) {
  log.Warn(msg)
}

func Error(msg string) {
  log.Error(msg)
}

func Errorf(format string, args ...interface{}) {
  log.Errorf(format, args...)
}
