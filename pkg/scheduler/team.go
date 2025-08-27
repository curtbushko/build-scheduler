package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	builds         []Build
	id             string
	totalQueueTime time.Duration
	totalBuildTime time.Duration
	plan           Plan
}

func NewTeam(name string, plan Plan) *Team {
	id := uuid.New()
	return &Team{
		id:   id.String(),
		plan: plan,
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
