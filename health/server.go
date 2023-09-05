package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type jsonResponse struct {
	Pid     int    `json:"pid"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func StartHTTPServer(
	ctx context.Context,
	addr string,
	port int,
	pid int,
	appVersion string,
	commit string,
) {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pid := os.Getpid()
		if err := json.NewEncoder(w).Encode(jsonResponse{
			Pid:     pid,
			Version: appVersion,
			Commit:  commit,
		}); err != nil {
			log.WithError(err).Error("responding to health-check")
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: m,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Info("Starting healthcheck server")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Error("error starting healtcheck server")
		panic(err)
	}
}
