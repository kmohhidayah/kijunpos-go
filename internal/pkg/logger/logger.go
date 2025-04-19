package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel/trace"
)

var (
	instance *logrus.Logger
	once     sync.Once
)

// GetLogger mengembalikan instance singleton dari logger
func GetLogger() *logrus.Logger {
	once.Do(func() {
		instance = &logrus.Logger{
			Out:          os.Stderr,
			Level:        logrus.InfoLevel,
			ReportCaller: true,
			Hooks:        make(logrus.LevelHooks),
		}
		// Set custom log formatter yang mengekstrak trace id dan lainnya dari context
		instance.SetFormatter(customLogger{
			formatter: logrus.JSONFormatter{
				PrettyPrint: false,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "timestamp",
				},
			},
		})
		// Tambahkan hook OpenTelemetry
		instance.Hooks.Add(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		)))
	})
	return instance
}

// custom struct to override logrus.Formatter
type customLogger struct {
	formatter logrus.JSONFormatter
}

// implement Format() to satisfy logrus.Formatter interface
func (l customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	// Ekstrak span data dari context
	span := trace.SpanFromContext(entry.Context)
	// Inject span data ke dalam log
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	return l.formatter.Format(entry)
}
