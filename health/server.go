package health

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/omegaatt36/instagramrobot/logging"
)

const (
	collectPeriod = 125 * time.Millisecond
	pollPeriod    = 3 * time.Second
)

var trackers sync.Map

var tmpResult result

type result struct {
	sync.Mutex
	alive bool
	ready bool
	vars  map[string]interface{}
}

// StartCollector starts collector.
func StartCollector() {
	for {
		r := gatherTrackers(false)

		tmpResult.Lock()
		tmpResult.alive = r.alive
		tmpResult.ready = r.ready
		tmpResult.vars = r.vars
		tmpResult.Unlock()

		time.Sleep(collectPeriod)
	}
}

func gatherTrackers(detail bool) result {
	var (
		alive = true
		ready = true
		vars  = make(map[string]interface{})
	)

	trackers.Range(func(key, value interface{}) bool {
		tracker := value.(Tracker)

		age := time.Since(tracker.lastTime)
		expired := tracker.interval > 0 && age > tracker.interval

		if tracker.ProbeType != ProbeNone {
			if expired || tracker.Status != statRunning {
				switch tracker.ProbeType {
				case ProbeAlive:
					alive = false
				case ProbeReady:
					ready = false
				}

				if detail {
					log.Printf("bad task(%v) expired(%v,%v) status(%v) type(%v)",
						tracker.taskName, expired, age, tracker.Status, tracker.ProbeType)
				}
			}

			if tracker.Status == statExited || tracker.Status == statInit {
				alive = false

				if detail {
					log.Printf("dead task(%v) status(%v)",
						tracker.taskName, tracker.ProbeType)
				}
			}
		}

		tracker.Age = age.String()
		tracker.Expired = expired
		vars[tracker.taskName] = tracker

		return true
	})

	vars["CollectedAt"] = time.Now()

	return result{
		alive: alive,
		ready: ready,
		vars:  vars,
	}
}

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	tmpResult.Lock()
	defer tmpResult.Unlock()

	if tmpResult.alive {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	tmpResult.Lock()
	defer tmpResult.Unlock()

	if tmpResult.ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func varHandler(w http.ResponseWriter, r *http.Request) {
	tmpResult.Lock()
	defer tmpResult.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tmpResult.vars); err != nil {
		logging.Error(err)
	}
}

// StartServer starts health server and blocks.
func StartServer() {
	go StartCollector()

	r := http.NewServeMux()
	r.HandleFunc("/alive", aliveHandler)
	r.HandleFunc("/ready", readyHandler)
	r.HandleFunc("/vars", varHandler)
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	r.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {
		gatherTrackers(true)
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    ":7001",
		Handler: r,
	}

	logging.Info("starts serving health server at", srv.Addr)
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		logging.Fatalf("listen: %s\n", err)
	}
}
