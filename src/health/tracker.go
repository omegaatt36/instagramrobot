package health

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// ProbeType defines probe type of tracker.
type ProbeType string

var serial int64

// ProbeType enums.
const (
	// ProbeNone won't influence alive or ready probe, used to report vars
	// in a context.
	ProbeNone ProbeType = "none"

	ProbeAlive ProbeType = "alive"
	ProbeReady ProbeType = "ready"
)

// Status is status, don't use it.
type Status string

const (
	statInit    Status = "init"
	statRunning Status = "running"
	statPause   Status = "pause"
	statExited  Status = "exited"
)

// Tracker defines tracker.
type Tracker struct {
	// for json marshalling only.
	Expired   bool
	Vars      map[string]interface{}
	Status    Status
	Age       string
	ProbeType ProbeType

	taskName string
	lastTime time.Time
	interval time.Duration
	deleted  bool
}

// NewTracker creates tracker, name must be unique.
func NewTracker(name string, d time.Duration, pt ProbeType) *Tracker {
	internalName := fmt.Sprintf("%s--%v", name, atomic.AddInt64(&serial, 1))
	t := &Tracker{
		taskName:  internalName,
		Status:    statInit,
		interval:  d,
		ProbeType: pt,
	}

	if _, loaded := trackers.LoadOrStore(internalName, *t); loaded {
		log.Panicln("name is duplicated:", name)
	}

	return t
}

// UpdateVars updates info.
func (i *Tracker) UpdateVars(vars map[string]interface{}) {
	i.lastTime = time.Now()
	i.Vars = vars
	trackers.Store(i.taskName, *i)
}

// Up marks this task status to up.
func (i *Tracker) Up() {
	i.lastTime = time.Now()
	i.Status = statRunning
	trackers.Store(i.taskName, *i)
}

// Down marks this task status to down and is unrecoverable.
func (i *Tracker) Down() {
	i.lastTime = time.Now()
	i.Status = statExited

	if !i.deleted {
		trackers.Store(i.taskName, *i)
	}

	log.Printf("task(%v) is exiting", i.taskName)
}

// Pause marks this task status to pause which is under self healing.
func (i *Tracker) Pause() {
	i.lastTime = time.Now()
	i.Status = statPause
	trackers.Store(i.taskName, *i)
}

// Unregister deletes the tracker from gathering table. it should be used for
// computation tasks.
func (i *Tracker) Unregister() {
	trackers.Delete(i.taskName)
	i.deleted = true
}

// TrackRoutine starts function f with recover to protect from panic,
// it sets alive probe to false when f returns or panics.
func (i *Tracker) TrackRoutine(f func()) {
	// Panic handling.
	defer func() {
		defer i.Down()
		if r := recover(); r != nil {
			// even though it's unregistered, we still need to report this issue
			// when getting panic.
			i.deleted = false
			fmt.Println("Recovered by health.Tracker: ", r)
			debug.PrintStack()
		}
	}()

	i.Up()

	f()
}
