package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	buildStatusComplete = "complete"
	buildStatusQueued   = "queued"
	buildStatusStarted  = "started"
)

type Build struct {
	ID         string
	Repo       string
	Branch     string
	Commit     string
	Status     string
	TeamID     string
	QueueStart time.Time
	BuildStart time.Time
	BuildEnd   time.Time
}

func NewBuild(repo, branch, commit, teamID string) *Build {
	id := uuid.New()
	return &Build{
		ID:         id.String(),
		Repo:       repo,
		Branch:     branch,
		Commit:     commit,
		TeamID:     teamID,
		QueueStart: time.Now(),
		Status:     buildStatusQueued,
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
	return nil
}
