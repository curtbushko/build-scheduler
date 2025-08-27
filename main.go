package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/curtbushko/build-scheduler/pkg/scheduler"
)

func main() {
	workerCount := 4
	queueSize := 10
	repo := "foo.com"
	branch := "main"
	commit := "c1398f"
	teamID := "1234"

	fmt.Println("Starting dispatcher")
	dispatcher := scheduler.NewDispatcher(workerCount, queueSize)
	dispatcher.Start()

	go func() {
		for i := 1; i <= 20; i++ {
			build := scheduler.NewBuild(repo, branch, commit, teamID)
			fmt.Println("Submitting build", build.ID)
			dispatcher.SubmitBuild(*build)
		}
	}()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // wait for interrupt
	fmt.Println("Shutting down gracefully...")
	dispatcher.Stop()
	fmt.Println("All workers stopped. Exiting.")
}
