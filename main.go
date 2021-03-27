package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap/zapcore"
)

var addr = flag.String("addr", ":8080", "listen address(default :8080)")
var username = flag.String("username", "", "Basic Auth Username")
var password = flag.String("password", "", "Basic Auth Passowrd")

func main() {
	// parse environment variable args.
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	h := NewHandler(os.Stdout, zapcore.InfoLevel)
	defer h.Logger.Sync()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", h.GetStatus)
	mux.HandleFunc("/time", h.GetTime)
	mux.HandleFunc("/sleep", h.GetSleep)

	// Required BasicAuth Handler
	if *username != "" || *password != "" {
		mux.HandleFunc("/ip", h.BasicAuthMiddleware(h.GetRemoteIP, *username, *password))
		mux.HandleFunc("/env", h.BasicAuthMiddleware(h.GetEnv, *username, *password))
	}

	log.Fatal(http.ListenAndServe(*addr, mux))
}
