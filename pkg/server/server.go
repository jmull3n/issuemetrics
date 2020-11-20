package server

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmull3n/issuemetrics/pkg/issuemetrics/github"
	"github.com/jmull3n/issuemetrics/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// Start starts the server
func Start(ctx context.Context, port string) {
	// setup routes and stuff
	log.Debug("setting up routes")
	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// wrap the response writer so we can get the response status
			respwriter := NewResponseWriter(w)

			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			ts := time.Now()
			// run the next now!
			next.ServeHTTP(respwriter, r)

			code := strconv.Itoa(respwriter.Status())
			metrics.RequestsTotal.WithLabelValues("issumetrics", path, code).Inc()
			metrics.RequestDurationMilliseconds.WithLabelValues("issumetrics", path, code).Observe(float64(time.Since(ts).Milliseconds()))
		})
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	issueManager := &github.IssueManager{}

	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/github", issueManager.IssueMetricsHandler).Methods(http.MethodPost, http.MethodOptions)

	go func() {
		log.Debug("starting server endpoint", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			select {
			case <-ctx.Done():
			default:
				log.Fatal("error starting Metrics server", "err", err)
			}
		}
	}()
}
