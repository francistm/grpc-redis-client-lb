package logger

import (
	"google.golang.org/grpc/grpclog"
)

var _ grpclog.LoggerV2 = (*NopLogger)(nil)

type NopLogger struct{}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (logPtr *NopLogger) Info(args ...interface{}) {
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (logPtr *NopLogger) Infoln(args ...interface{}) {
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (logPtr *NopLogger) Infof(format string, args ...interface{}) {
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (logPtr *NopLogger) Warning(args ...interface{}) {
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (logPtr *NopLogger) Warningln(args ...interface{}) {
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (logPtr *NopLogger) Warningf(format string, args ...interface{}) {
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (logPtr *NopLogger) Error(args ...interface{}) {
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (logPtr *NopLogger) Errorln(args ...interface{}) {
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (logPtr *NopLogger) Errorf(format string, args ...interface{}) {
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logPtr *NopLogger) Fatal(args ...interface{}) {
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logPtr *NopLogger) Fatalln(args ...interface{}) {
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (logPtr *NopLogger) Fatalf(format string, args ...interface{}) {
}

// V reports whether verbosity level l is at least the requested verbose level.
func (logPtr *NopLogger) V(l int) bool {
	return false
}
