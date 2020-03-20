package plugin

import (
	"errors"
	"io"
	"net/http"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/go-plugins-helpers/sdk"
)

type startLoggingRequest struct {
	File string
	Info logger.Info
}

type stopLoggingRequest struct {
	File string
}

type capabilitiesResponse struct {
	Err string
	Cap logger.Capability
}

type readLogsRequest struct {
	Info   logger.Info
	Config logger.ReadConfig
}

type response struct {
	Err string
}

// HTTPHandler is a helper struct that implements all the http handlers required
// for a logging driver
type HTTPHandler struct {
	Plugin Plugin
}

// StartLogging is an http handler for the /LogDriver.StartLogging endpoint
func (h *HTTPHandler) StartLogging(w http.ResponseWriter, r *http.Request) {
	var req startLoggingRequest
	if err := sdk.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if req.Info.ContainerID == "" {
		respond(errors.New("must provide container id in log context"), w)
		return
	}

	err := h.Plugin.StartLogging(req.File, req.Info)
	respond(err, w)
}

// StopLogging is an http handler for the /LogDriver.StopLogging endpoint
func (h *HTTPHandler) StopLogging(w http.ResponseWriter, r *http.Request) {
	var req stopLoggingRequest
	if err := sdk.DecodeRequest(w, r, &req); err != nil {
		return
	}

	err := h.Plugin.StopLogging(req.File)
	respond(err, w)
}

// Capabilities is an http handler for the /LogDriver.Capabilities endpoint
func (h *HTTPHandler) Capabilities(w http.ResponseWriter, r *http.Request) {
	res := &capabilitiesResponse{
		Cap: h.Plugin.Capabilities(),
	}
	sdk.EncodeResponse(w, res, false)
}

// ReadLogs is an http handler for the /LogDriver.ReadLogs endpoint
func (h *HTTPHandler) ReadLogs(w http.ResponseWriter, r *http.Request) {
	var req readLogsRequest

	if err := sdk.DecodeRequest(w, r, &req); err != nil {
		return
	}

	stream, err := h.Plugin.ReadLogs(req.Info, req.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", "application/x-json-stream")
	wf := ioutils.NewWriteFlusher(w)
	io.Copy(wf, stream)
}

// Initialize assign the handlers to the given sdk handler
func (h *HTTPHandler) Initialize(sh *sdk.Handler) {
	sh.HandleFunc("/LogDriver.StartLogging", h.StartLogging)
	sh.HandleFunc("/LogDriver.StopLogging", h.StopLogging)
	sh.HandleFunc("/LogDriver.Capabilities", h.Capabilities)
	sh.HandleFunc("/LogDriver.ReadLogs", h.ReadLogs)
}

func respond(err error, w http.ResponseWriter) {
	var res response
	if err != nil {
		res.Err = err.Error()
	}
	sdk.EncodeResponse(w, &res, err != nil)
}
