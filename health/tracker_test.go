package health

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGatherTrackers(t *testing.T) {
	s := assert.New(t)

	type probeResult struct {
		alive bool
		ready bool
	}

	type testdata struct {
		d  time.Duration
		pt ProbeType

		whenInit       probeResult
		whenRun        probeResult
		whenTimeout    probeResult
		whenDown       probeResult
		whenUnregister probeResult
	}

	datas := []testdata{
		// ProbeNone won't make any probe unhappy.
		{
			d:              time.Second,
			pt:             ProbeNone,
			whenInit:       probeResult{alive: true, ready: true},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: true, ready: true},
			whenDown:       probeResult{alive: true, ready: true},
			whenUnregister: probeResult{alive: true, ready: true},
		},
		{
			d:              0,
			pt:             ProbeNone,
			whenInit:       probeResult{alive: true, ready: true},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: true, ready: true},
			whenDown:       probeResult{alive: true, ready: true},
			whenUnregister: probeResult{alive: true, ready: true},
		},
		// ProbeReady influences ready probe.
		{
			d:              time.Second,
			pt:             ProbeReady,
			whenInit:       probeResult{alive: false, ready: false},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: true, ready: false},
			whenDown:       probeResult{alive: false, ready: false},
			whenUnregister: probeResult{alive: true, ready: true},
		},
		// ProbeReady without d only turn it to false when tracker is down.
		{
			d:              0,
			pt:             ProbeReady,
			whenInit:       probeResult{alive: false, ready: false},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: true, ready: true},
			whenDown:       probeResult{alive: false, ready: false},
			whenUnregister: probeResult{alive: true, ready: true},
		},
		// ProbeAlive influences alive probe.
		{
			d:              time.Second,
			pt:             ProbeAlive,
			whenInit:       probeResult{alive: false, ready: true},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: false, ready: true},
			whenDown:       probeResult{alive: false, ready: true},
			whenUnregister: probeResult{alive: true, ready: true},
		},
		// ProbeReady without d only turn it to false when tracker is down.
		{
			d:              0,
			pt:             ProbeAlive,
			whenInit:       probeResult{alive: false, ready: true},
			whenRun:        probeResult{alive: true, ready: true},
			whenTimeout:    probeResult{alive: true, ready: true},
			whenDown:       probeResult{alive: false, ready: true},
			whenUnregister: probeResult{alive: true, ready: true},
		},
	}

	for i, data := range datas {
		// clear map
		trackers = sync.Map{}

		// no trackers, should be good.
		res := gatherTrackers(false)
		s.True(res.alive)
		s.True(res.ready)

		// create tracker.
		t1 := NewTracker(
			fmt.Sprintf("test%v", i),
			data.d,
			data.pt)

		// at status init
		res = gatherTrackers(false)
		s.Equal(data.whenInit.alive, res.alive, i)
		s.Equal(data.whenInit.ready, res.ready, i)

		// at status running
		t1.Up()
		res = gatherTrackers(false)
		s.Equal(data.whenRun.alive, res.alive, i)
		s.Equal(data.whenRun.ready, res.ready, i)

		// timeout.
		time.Sleep(data.d + time.Second)
		res = gatherTrackers(false)
		s.Equal(data.whenTimeout.alive, res.alive, i)
		s.Equal(data.whenTimeout.ready, res.ready, i)

		// at status exited.
		t1.Down()
		res = gatherTrackers(false)
		s.Equal(data.whenDown.alive, res.alive, i)
		s.Equal(data.whenDown.ready, res.ready, i)

		t1.Unregister()
		res = gatherTrackers(false)
		s.Equal(data.whenUnregister.alive, res.alive, i)
		s.Equal(data.whenUnregister.ready, res.ready, i)
	}

	// clear map
	trackers = sync.Map{}
}

func TestProtectPanic(t *testing.T) {
	s := assert.New(t)

	// clear map
	trackers = sync.Map{}

	tracker := NewTracker("test1", 0, ProbeAlive)
	s.NotPanics(func() {
		tracker.TrackRoutine(func() {
			s.Equal(statRunning, tracker.Status)

			res := gatherTrackers(false)
			s.Equal(true, res.alive)
			s.Equal(true, res.ready)
		})
	})
	s.Equal(statExited, tracker.Status)

	res := gatherTrackers(false)
	s.Equal(false, res.alive)
	s.Equal(true, res.ready)

	// clear map
	trackers = sync.Map{}

	// add panic.
	tracker = NewTracker("test1", 0, ProbeAlive)
	s.NotPanics(func() {
		tracker.TrackRoutine(func() {
			s.Equal(statRunning, tracker.Status)

			res := gatherTrackers(false)
			s.Equal(true, res.alive)
			s.Equal(true, res.ready)

			panic("failed")
		})
	})
	s.Equal(statExited, tracker.Status)

	res = gatherTrackers(false)
	s.Equal(false, res.alive)
	s.Equal(true, res.ready)

}

func TestGatherVars(t *testing.T) {
	s := assert.New(t)

	// clear map
	trackers = sync.Map{}
	serial = 0

	res := gatherTrackers(false)
	s.Equal(1, len(res.vars))

	tracker := NewTracker("var1", 0, ProbeNone)
	tracker.UpdateVars(map[string]interface{}{
		"test": "123",
	})

	res = gatherTrackers(false)
	s.Equal(2, len(res.vars))
	s.Contains(res.vars, "var1--1")

	tracker.Unregister()

	res = gatherTrackers(false)
	s.Equal(1, len(res.vars))
}

func TestProtectPanicWithUnregistered(t *testing.T) {
	s := assert.New(t)

	// clear map
	trackers = sync.Map{}

	tracker2 := NewTracker("test2", 0, ProbeAlive)
	s.NotPanics(func() {
		tracker2.TrackRoutine(func() {
			defer tracker2.Unregister()

			s.Equal(statRunning, tracker2.Status)
			tracker2.UpdateVars(map[string]interface{}{
				"info_in_tracker2": 123,
			})
		})
	})
	s.Equal(statExited, tracker2.Status)

	res := gatherTrackers(false)
	s.Equal(true, res.alive)
	s.Equal(true, res.ready)
	s.Equal(1, len(res.vars))

	// clear map
	trackers = sync.Map{}

	// with panic this time.
	tracker3 := NewTracker("test3", 0, ProbeAlive)
	s.NotPanics(func() {
		tracker3.TrackRoutine(func() {
			defer tracker3.Unregister()

			s.Equal(statRunning, tracker3.Status)
			panic("gg")
		})
	})
	s.Equal(statExited, tracker3.Status)

	res = gatherTrackers(false)
	s.Equal(false, res.alive) // because tracker3's routine was panic.
	s.Equal(true, res.ready)
}
