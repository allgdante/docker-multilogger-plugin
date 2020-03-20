package plugin

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	protoio "github.com/gogo/protobuf/io"
	"github.com/sirupsen/logrus"
)

type adapter struct {
	id     string
	logger logger.Logger
	stream io.ReadCloser
}

func (a *adapter) Start() {
	dec := protoio.NewUint32DelimitedReader(a.stream, binary.BigEndian, 1e6)
	defer dec.Close()
	defer a.Close()

	var buf logdriver.LogEntry

	for {
		if err := dec.ReadMsg(&buf); err != nil {
			if err == io.EOF || err == os.ErrClosed || strings.Contains(err.Error(), "file already closed") {
				logrus.WithField("id", a.id).WithError(err).Debug("shutting down log consumer")
				return
			}

			logrus.WithField("id", a.id).WithError(err).Error("received unexpected error. retrying...")
			dec = protoio.NewUint32DelimitedReader(a.stream, binary.BigEndian, 1e6)
		}

		var msg logger.Message
		msg.Line = buf.Line
		msg.Source = buf.Source
		if buf.PartialLogMetadata != nil {
			if msg.PLogMetaData == nil {
				msg.PLogMetaData = &backend.PartialLogMetaData{}
			}
			msg.PLogMetaData.ID = buf.PartialLogMetadata.Id
			msg.PLogMetaData.Last = buf.PartialLogMetadata.Last
			msg.PLogMetaData.Ordinal = int(buf.PartialLogMetadata.Ordinal)
		}
		msg.Timestamp = time.Unix(0, buf.TimeNano)
		if err := a.logger.Log(&msg); err != nil {
			logrus.WithFields(logrus.Fields{
				"source":  buf.Source,
				"message": msg,
			}).WithError(err).Error("Error writing log message")
		}

		buf.Reset()
	}
}

func (a *adapter) Close() {
	a.stream.Close()
}
