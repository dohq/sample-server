package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap/zapcore"
)

var addr = flag.String("addr", ":8080", "listen address")

func main() {
	flag.Parse()

	h := NewHandler(os.Stdout, zapcore.InfoLevel)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", h.GetStatus)
	mux.HandleFunc("/time", h.GetTime)
	mux.HandleFunc("/sleep", h.GetSleep)
	mux.HandleFunc("/ip", h.GetRemoteIP)

	log.Fatal(http.ListenAndServe(*addr, mux))
}
