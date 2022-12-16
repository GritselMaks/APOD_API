package logger

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	logger := NewLogger("")
	if logger == nil {
		t.Fatal("logger to never be nil")
	}
}

func Test_levelToZapLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args string
		want zapcore.Level
	}{
		{args: levelDebug, want: zapcore.DebugLevel},
		{args: levelInfo, want: zapcore.InfoLevel},
		{args: levelWarning, want: zapcore.WarnLevel},
		{args: levelError, want: zapcore.ErrorLevel},
		{args: levelCritical, want: zapcore.DPanicLevel},
		{args: levelAlert, want: zapcore.PanicLevel},
		{args: levelEmergency, want: zapcore.FatalLevel},
		{args: "", want: zapcore.WarnLevel},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := levelToZapLevel(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("levelToZapLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
