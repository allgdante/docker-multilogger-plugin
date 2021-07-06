package logassembler

import (
	"testing"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssembler(t *testing.T) {
	const (
		msgSize = 10
		maxSize = 20
	)

	var (
		assert  = assert.New(t)
		require = require.New(t)
		a       = New(maxSize)
	)

	for _, tc := range []struct {
		Source []byte
		Total  int
		Lines  [][]byte
	}{
		{
			[]byte("aaaa"),
			1,
			[][]byte{[]byte("aaaa")},
		},
		{
			[]byte("aaaaaaaabbbbbbbb"),
			1,
			[][]byte{[]byte("aaaaaaaabbbbbbbb")},
		},
		{
			[]byte("01234567890123456789"),
			1,
			[][]byte{[]byte("01234567890123456789")},
		},
		{
			[]byte("123456781234567812345678"),
			2,
			[][]byte{
				[]byte("12345678123456781234"),
				[]byte("5678"),
			},
		},
		{
			[]byte("01234567890123456789012345678901234567890"),
			3,
			[][]byte{
				[]byte("01234567890123456789"),
				[]byte("01234567890123456789"),
				[]byte("0"),
			},
		},
	} {
		t.Run(string(tc.Source), func(t *testing.T) {
			var msgs []*logger.Message

			for i, o := 0, 1; i < len(tc.Source); {
				var (
					n    int
					meta *backend.PartialLogMetaData
				)

				if len(tc.Source[i:]) > msgSize {
					n = msgSize
					meta = &backend.PartialLogMetaData{
						ID:      "meta",
						Ordinal: o,
						Last:    false,
					}
					o++
				} else {
					n = len(tc.Source[i:])
					if o != 1 {
						meta = &backend.PartialLogMetaData{
							ID:      "meta",
							Ordinal: o,
							Last:    true,
						}
					}
				}

				msgs = append(msgs, a.Assemble(newMessage(tc.Source[i:i+n], meta))...)
				i += n
			}

			require.Len(msgs, tc.Total)
			for i, msg := range msgs {
				assert.Equal(tc.Lines[i], msg.Line)
				logger.PutMessage(msg)
			}
		})
	}
}

func newMessage(message []byte, meta *backend.PartialLogMetaData) *logger.Message {
	msg := logger.NewMessage()
	msg.Line = append(msg.Line, []byte(message)...)
	msg.PLogMetaData = meta
	return msg
}
