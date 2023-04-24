package health

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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

func aliveHandler(ctx *gin.Context) {
	defer ctx.Abort()

	tmpResult.Lock()
	defer tmpResult.Unlock()

	if tmpResult.alive {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusServiceUnavailable)
	}
}

func readyHandler(ctx *gin.Context) {
	defer ctx.Abort()

	tmpResult.Lock()
	defer tmpResult.Unlock()

	if tmpResult.ready {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusServiceUnavailable)
	}
}

func varHandler(ctx *gin.Context) {
	defer ctx.Abort()

	tmpResult.Lock()
	defer tmpResult.Unlock()

	ctx.JSON(http.StatusOK, tmpResult.vars)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RegisterToGinEngine registers health check endpoint to an existing engine.
func RegisterToGinEngine(engine *gin.Engine) {
	engine.GET("/alive", aliveHandler)
	engine.GET("/ready", readyHandler)
	engine.GET("/vars", varHandler)
	engine.GET("/metrics", prometheusHandler())
	engine.GET("/dump", func(c *gin.Context) {
		gatherTrackers(true)
	})
}

var engine *gin.Engine

// Engine returns engine.
func Engine() *gin.Engine {
	return engine
}

func init() {
	engine = gin.New()
	engine.RedirectTrailingSlash = true
	RegisterToGinEngine(engine)
}

// StartServer starts health server and blocks.
func StartServer() {
	go StartCollector()

	srv := &http.Server{
		Addr:    ":7001",
		Handler: engine,
	}

	log.Info("starts serving health server at", srv.Addr)
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
