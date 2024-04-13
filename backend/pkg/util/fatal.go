package util

import (
	"fmt"
	"log/slog"
)

func Fatal(msg string, arg ...any) {
	err := fmt.Errorf(msg, arg...)
	slog.Error(err.Error())
	panic(err.Error())
}
