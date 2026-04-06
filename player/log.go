package player

import (
	"fmt"
	"strings"
)

type Log interface {
  Logf(format string, a ...any)
}

type MessageLog struct {
  messages []string
}

func CreateLog() *MessageLog {
  return &MessageLog { messages: []string {} }
}

func (log *MessageLog) Logf(format string, a ...any) {
  var stringBuilder strings.Builder
  fmt.Fprintf(&stringBuilder, format, a...)
  log.messages = append(log.messages, stringBuilder.String())
}

func (log *MessageLog) Clear() {
  log.messages = []string {}
}

func (log *MessageLog) GetMessages() []string {
  return log.messages
}

type nilLog struct {}

var NilLog Log = &nilLog{}

func (*nilLog) Logf(format string, a ...any) {}
