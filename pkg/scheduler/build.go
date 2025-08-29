package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	buildStatusComplete  = "complete"
	buildStatusQueued    = "queued"
	buildStatusDispatched = "dispatched"
	buildStatusStarted   = "started"
)

type Build struct {
	ID              string
	Repo            string
	Branch          string
	Commit          string
	Status          string
	TeamID          string
	QueueStart      time.Time
	BuildStart      time.Time
	BuildEnd        time.Time
	CompletionNotif chan<- string
}

func NewBuild(repo, branch, commit, teamID string, completionChan chan<- string) *Build {
	id := uuid.New()
	return &Build{
		ID:              id.String(),
		Repo:            repo,
		Branch:          branch,
		Commit:          commit,
		TeamID:          teamID,
		QueueStart:      time.Now(),
		Status:          buildStatusQueued,
		CompletionNotif: completionChan,
	}
}

func (b *Build) Process() error {
	b.Status = buildStatusStarted
	b.BuildStart = time.Now()
	fmt.Println("Running build", b.ID)
	time.Sleep(1 * time.Second)
	b.BuildEnd = time.Now()
	b.Status = buildStatusComplete
	fmt.Println("Build complete", b.ID)
	
	if b.CompletionNotif != nil {
		b.CompletionNotif <- b.ID
	}
	
	return nil
}
