package logassembler

import (
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/daemon/logger"
)

// Assembler represents an object capable of assemble logger messages
type Assembler interface {
	Assemble(msg *logger.Message) []*logger.Message
}

// New returns a new Assembler
func New(size int64) Assembler {
	return &assembler{
		line: make([]byte, 0, size),
	}
}

type assembler struct {
	line    []byte
	id      string
	ordinal int
	last    *logger.Message
}

// Assemble implements the LogAssembler interface
func (a *assembler) Assemble(msg *logger.Message) (msgs []*logger.Message) {
	if msg.PLogMetaData == nil {
		if a.id != "" {
			msgs = append(msgs, a.generate(true))
			a.reset()
		}
		msgs = append(msgs, msg)
		return
	}

	if a.id != "" && a.id != msg.PLogMetaData.ID {
		msgs = append(msgs, a.generate(true))
		a.reset()
	}

	if a.id == "" {
		a.init(msg)
	}

	if len(a.line)+len(msg.Line) > cap(a.line) {
		n := cap(a.line) - len(a.line)
		a.line = append(a.line, msg.Line[:n]...)
		msgs = append(msgs, a.generate(false))
		a.ordinal++
		a.line = append(a.line[:0], msg.Line[n:]...)
	} else {
		a.line = append(a.line, msg.Line...)
	}

	if msg.PLogMetaData.Last {
		msgs = append(msgs, a.generate(true))
		a.reset()
	}

	logger.PutMessage(msg)
	return
}

func (a *assembler) init(msg *logger.Message) {
	a.id = msg.PLogMetaData.ID
	a.ordinal = 1
	a.last = logger.NewMessage()
	copyMessage(a.last, msg)
}

func (a *assembler) reset() {
	a.line = a.line[:0]
	a.id = ""
	a.ordinal = 1
	logger.PutMessage(a.last)
	a.last = nil
}

func (a *assembler) generate(isLastPartial bool) *logger.Message {
	msg := logger.NewMessage()
	msg.Line = append(msg.Line[:0], a.line...)
	copyMessage(msg, a.last)
	if !isLastPartial || (isLastPartial && a.ordinal != 1) {
		msg.PLogMetaData = &backend.PartialLogMetaData{
			ID:      a.id,
			Ordinal: a.ordinal,
			Last:    isLastPartial,
		}
	}
	return msg
}

func copyMessage(dst, src *logger.Message) {
	dst.Source = src.Source
	dst.Timestamp = src.Timestamp
	dst.Err = src.Err
	dst.Attrs = append(dst.Attrs, src.Attrs...)
}
