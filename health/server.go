package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	palomaWasOnline = false
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Pid     int    `json:"pid"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func StartHTTPServer(
	ctx context.Context,
	port int,
	pid int,
	appVersion string,
	commit string,
) {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pid := os.Getpid()
		if err := json.NewEncoder(w).Encode(jsonResponse{
			OK:      palomaWasOnline,
			Pid:     pid,
			Version: appVersion,
			Commit:  commit,
		}); err != nil {
			log.WithError(err).Error("responding to health-check")
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: m,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	log.Info("Starting healthcheck server")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Error("error starting healtcheck server")
		panic(err)
	}

}
