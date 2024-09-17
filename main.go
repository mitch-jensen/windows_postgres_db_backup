package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			getAllDatabases,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}

func getAllDatabases() []string {
	query := "SELECT datname FROM pg_database WHERE datistemplate = false;"
	cmd := exec.Command("psql", "--host", "127.0.0.1", "--port", "5432", "--username", "postgres", "--tuples-only", "--command", query)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	databases := strings.Fields(strings.TrimSpace(string(out)))
	return databases
}
