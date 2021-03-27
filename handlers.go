package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Handler with a simple logger.
type Handler struct {
	Logger *zap.Logger
}

// NewHandler is make new handler.
func NewHandler(ws *os.File, loglevel zapcore.LevelEnabler) *Handler {
	ec := ecszap.EncoderConfig{
		EnableName:       true,
		EnableStackTrace: true,
		EnableCaller:     false,
		EncodeName:       zapcore.FullNameEncoder,
		// upper case log level output.
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   ecszap.ShortCallerEncoder,
	}
	core := ecszap.NewCore(ec, ws, loglevel)
	logger := zap.New(core, zap.AddCaller()).Named("sample-server")

	return &Handler{
		Logger: logger,
	}
}

// GetStatus is A simple endpoint that just returns 200.
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	result := Result{Result: "ok"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	h.Logger.Info("ok", zap.Int("status", http.StatusOK))
}

// GetTime is simple response server time.
func (h *Handler) GetTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()

	result := Result{Result: now}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	h.Logger.Info(now, zap.Int("status", http.StatusOK))
}

// GetSleep is Sleeps for the time passed in the query and returns the response.
func (h *Handler) GetSleep(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("d")

	if q == "" {
		http.Error(w, "duration is blank", http.StatusInternalServerError)
		h.Logger.Warn("duration is blank", zap.Int("status", http.StatusInternalServerError))
		return
	}

	dur, err := time.ParseDuration(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Warn(q, zap.Int("status", http.StatusInternalServerError))
		return
	}

	time.Sleep(dur)

	result := Result{Result: q}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	h.Logger.Info(q, zap.Int("status", http.StatusOK))
}

// GetRemoteIP is return remote ip include X-Forwarded-For
func (h *Handler) GetRemoteIP(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	xf := w.Header().Get("X-Forwarded-For")

	if xf != "" {
		addr = xf
	}

	result := Result{Result: addr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	h.Logger.Info(addr, zap.Int("status", http.StatusOK))
}

// GetEnv is return remote ip include X-Forwarded-For
func (h *Handler) GetEnv(w http.ResponseWriter, r *http.Request) {
	envs := os.Environ()
	m := make(map[string]string)
	for _, i := range envs {
		e := strings.Split(i, "=")
		m[e[0]] = e[1]
	}

	result := Result{Result: "ok", Envs: m}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	h.Logger.Info("ok", zap.Int("count", len(envs)))
}
