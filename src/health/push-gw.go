package health

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// StartVarsPoller starts vars poller and triggers callback function.
func StartVarsPoller(ctx context.Context, fn func([]byte)) {
	go func() {
		for ctx.Err() == nil {
			time.Sleep(pollPeriod)

			tmpResult.Lock()
			vars := tmpResult.vars
			tmpResult.Unlock()

			b, err := json.Marshal(vars)
			if err != nil {
				log.Error(err)
				continue
			}

			fn(b)
		}
	}()
}
