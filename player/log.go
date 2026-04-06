package player

import (
	"fmt"
	"strings"
)

type Log struct {
  messages []string
}

func CreateLog() *Log {
  return &Log { messages: []string {} }
}

func (log *Log) Logf(format string, a ...any) {
  var stringBuilder strings.Builder
  fmt.Fprintf(&stringBuilder, format, a...)
  log.messages = append(log.messages, stringBuilder.String())
}

func (log *Log) Clear() {
  log.messages = []string {}
}

func (log *Log) GetMessages() []string {
  return log.messages
}
