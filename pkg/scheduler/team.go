package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	builds           []Build
	id               string
	name             string
	totalQueueTime   time.Duration
	totalBuildTime   time.Duration
	plan             Plan
	runningBuilds    int
	dispatcher       *Dispatcher
	CompletionNotifs chan string
	ctx              context.Context
	cancel           context.CancelFunc
	mu               sync.RWMutex
}

func NewTeam(name string, plan Plan) *Team {
	id := uuid.New()
	ctx, cancel := context.WithCancel(context.Background())
	return &Team{
		id:               id.String(),
		name:             name,
		plan:             plan,
		CompletionNotifs: make(chan string, 100),
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (t *Team) SetDispatcher(dispatcher *Dispatcher) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.dispatcher = dispatcher
}

func (t *Team) SubmitBuild(build Build) {
	t.mu.Lock()
	defer t.mu.Unlock()
	build.CompletionNotif = t.CompletionNotifs
	t.builds = append(t.builds, build)
}

func (t *Team) Start() {
	go t.processQueue()
	go t.handleCompletions()
}

func (t *Team) Stop() {
	t.cancel()
}

func (t *Team) processQueue() {
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			t.tryForwardBuild()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *Team) tryForwardBuild() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.runningBuilds >= t.plan.Slots || len(t.builds) == 0 || t.dispatcher == nil {
		return
	}

	// Find the first queued build
	for i := range t.builds {
		if t.builds[i].Status == "queued" {
			// Mark as dispatched and increment running builds
			t.builds[i].Status = "dispatched"
			t.runningBuilds++

			// Submit to dispatcher
			go func(build Build) {
				t.dispatcher.SubmitBuildFromTeam(build, t)
			}(t.builds[i])
			return
		}
	}
}

func (t *Team) handleCompletions() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case buildID := <-t.CompletionNotifs:
			t.mu.Lock()
			for i, build := range t.builds {
				if build.ID == buildID {
					t.builds[i].Status = buildStatusComplete
					t.builds[i].BuildEnd = time.Now()
					fmt.Printf("Build completed: ID=%s, Status=%s, Team=%s, Duration=%v\n",
						build.ID, t.builds[i].Status, build.TeamID, t.builds[i].BuildEnd.Sub(build.BuildStart))
					t.builds = append(t.builds[:i], t.builds[i+1:]...)
					break
				}
			}
			if t.runningBuilds > 0 {
				t.runningBuilds--
			}
			t.mu.Unlock()
		}
	}
}

func (t *Team) OnBuildComplete() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.runningBuilds > 0 {
		t.runningBuilds--
	}
}

// Upsert adds a build if it does not exists and updates a build if it does exist
func (t *Team) Upsert(build Build) error {
	for i, b := range t.builds {
		if b.Branch == build.Branch {
			t.builds[i].Commit = build.Commit
			t.builds[i].QueueStart = build.QueueStart
			return nil
		}
	}
	t.builds = append(t.builds, build)
	return nil
}
