package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/google/uuid"
)

func main() {
	repo := job.NewInMemoryRepository()
	service := job.NewService(repo)
	ctx := context.Background()
	if len(os.Args) < 2 {
		fmt.Printf("Insufficient Args: gofleet <submit|status|cancel> [flags] \n")
		return
	}
	switch os.Args[1] {
	case "submit":
		fs := flag.NewFlagSet("submit", flag.ExitOnError)
		cmd := fs.String("cmd", "", "cmd")
		fs.Parse(os.Args[2:])
		if *cmd == "" {
			fmt.Println("cmd flag not set for submit type")
			return
		}
		job, err := service.Submit(ctx, *cmd, 5, "high")
		if err != nil {
			fmt.Printf("Error while submitting job : %v \n", err)
			return
		}
		fmt.Printf("Created Job Successfully , ID: %v \n", job.ID)

	case "status":
		fs := flag.NewFlagSet("status", flag.ExitOnError)
		jobID := fs.String("job-id", "", "job ID")
		fs.Parse(os.Args[2:])
		jobUUID, err := uuid.Parse(*jobID)
		if err != nil {
			fmt.Printf("Error in converting string to uuid: %v \n", err)
			return
		}
		job, err := service.GetJob(ctx, jobUUID)
		if err != nil {
			fmt.Printf("Job with flag job-id not found: %v \n", jobUUID)
			return
		}
		fmt.Printf("job with JobID found: ID- %v, Status- %v", job.ID, job.Status)

	case "cancel":
		fs := flag.NewFlagSet("cancel", flag.ExitOnError)
		jobID := fs.String("job-id", "", "job ID")
		fs.Parse(os.Args[2:])
		jobUUID, err := uuid.Parse(*jobID)
		if err != nil {
			fmt.Printf("Error in converting string to uuid: %v \n", err)
			return
		}
		if err := service.Transition(ctx, jobUUID, job.StatusCancelled); err != nil {
			fmt.Printf("cancel err : %v \n", err)
		}

	default:
		fmt.Println("Invalid command or command format")

	}

}
