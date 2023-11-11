package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

var instance Logger
var once sync.Once

func (s *Logger) ExtraFields(fields map[string]interface{}) *Logger {
	return &Logger{s.WithFields(fields)}
}

func GetLogger(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)

		if err != nil {
			log.Fatalln(err)
		}

		l := logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
			FullTimestamp: true, // ?
			DisableColors: false,
		}

		l.SetOutput(os.Stdout)

		l.SetLevel(logrusLevel)

		instance = Logger{logrus.NewEntry(l)}
	})
	return instance
}
