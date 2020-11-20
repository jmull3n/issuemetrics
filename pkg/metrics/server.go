package metrics

import (
	"github.com/jmull3n/issuemetrics/pkg/util"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartServer starts the metricserver
func StartServer(ctx context.Context, port string) {
	// run webserver
	http.Handle("/metrics", promhttp.Handler())

	// expose a `/ready` endpoint so we can standardize readiness checks & liveness checks
	http.HandleFunc("/ready", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// clients will be expected to implement custom /health checks based on a service's requirements

	server := &http.Server{
		Addr: ":" + port,
	}
	go func() {
		log.Debug("starting /metrics endpoint", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			select {
			case <-ctx.Done():
			default:
				log.Fatal("error starting Metrics server", "err", err)
			}
		}
	}()

	pprofPort := util.Getenv("PROFILE_PORT", "")
	if len(pprofPort) > 0 {
		log.Debug("pprof listening on port ", "pprofPort", pprofPort)

		go func() {
			fmt.Println(http.ListenAndServe(":"+pprofPort, nil))
		}()
	}
}
