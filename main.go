package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			dumpAllDatabases,
		),
	)
	if err != nil {
		panic(err)
	}

	log.Printf("Job: %v\n", j)

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		panic(err)
	}
}

func getPostgresPath() string {
	return "C:\\Program Files\\PostgreSQL\\15\\bin"
}

func getPsqlExecutable() string {
	return filepath.Join(getPostgresPath(), "psql.exe")
}

func getPgDumpExecutable() string {
	return filepath.Join(getPostgresPath(), "pg_dump.exe")
}

func getBackupDir() string {
	return "C:\\backups"
}

func createBackupDir() {
	err := os.MkdirAll(getBackupDir(), 0750)
	if err != nil {
		log.Fatal(err)
	}
}

func createBackupFilename(database string) string {
	createBackupDir()
	current_date := time.Now().Format("2006-01-02") // YYYY-MM-DD
	return filepath.Join(getBackupDir(), database+"_"+current_date+".backup")
}

func getAllDatabases() []string {
	query := "SELECT datname FROM pg_database WHERE datistemplate = false;"
	cmd := exec.Command(getPsqlExecutable(), "--host", "127.0.0.1", "--port", "5432", "--username", "postgres", "--tuples-only", "--command", query)
	log.Printf("Running command: %v\n", cmd)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Output: %s\n", out)
	databases := strings.Fields(strings.TrimSpace(string(out)))
	log.Printf("Databases: %v\n", databases)
	return databases
}

func dumpDatabase(database string) {
	log.Printf("Dumping database %s\n", database)
	dump_file := createBackupFilename(database)
	log.Printf("Dumping database %s to %s\n", database, dump_file)
	cmd := exec.Command(getPgDumpExecutable(), "--host", "127.0.0.1", "--port", "5432", "--username", "postgres", "--format=custom", "-b", "--verbose", "--dbname", database, "--file", dump_file)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	} else {
		log.Printf("Output: %s\n", out)
		log.Printf("Database %s dumped to %s\n", database, dump_file)
	}
}

func dumpAllDatabases() {
	databases := getAllDatabases()
	for _, database := range databases {
		dumpDatabase(database)
	}
}
