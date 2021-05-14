package Logging

import (
	"fmt"
	"time"
)

const format = "Mon Jan 2 15:04:05 MST"

type ILogger interface {
	Log(message string)
}

type Logger struct {
}

func (logger *Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", timeStamp(), message)
}

func timeStamp() string {
	return time.Now().Format(format)
}
