package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/Yashasvi-webmob/gofleet/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient(
		"localhost:5000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewJobServiceClient(conn)

	ctx := context.Background()
	if len(os.Args) < 2 {
		fmt.Printf("Insufficient Args: gofleet <submit|status|cancel> [flags] \n")
		return
	}
	switch os.Args[1] {
	case "submit":
		fs := flag.NewFlagSet("submit", flag.ExitOnError)
		cmd := fs.String("cmd", "", "cmd")
		timeout := fs.Int("timeout", 5, "timeout in seconds")
		priority := fs.String("priority", "high", "job priority")
		fs.Parse(os.Args[2:])
		if *cmd == "" {
			fmt.Println("cmd flag not set for submit type")
			return
		}
		job, err := client.CreateJob(ctx, &pb.SubmitJobRequest{Command: *cmd, TimeoutSeconds: int32(*timeout), Priority: *priority})
		if err != nil {
			fmt.Printf("Error while submitting job : %v \n", err)
			return
		}
		fmt.Printf("Created Job Successfully , ID: %v \n", job.Id)

	case "status":
		fs := flag.NewFlagSet("status", flag.ExitOnError)
		jobID := fs.String("job-id", "", "job ID")
		fs.Parse(os.Args[2:])
		jobUUID, err := uuid.Parse(*jobID)
		if err != nil {
			fmt.Printf("Error in converting string to uuid: %v \n", err)
			return
		}
		job, err := client.GetJob(ctx, &pb.GetJobRequest{Id: *jobID})
		if err != nil {
			fmt.Printf("Job with flag job-id not found: %v \n", jobUUID)
			return
		}
		fmt.Printf("job with JobID found: ID- %v, Status- %v", job.Id, job.Status)

	case "cancel":
		fs := flag.NewFlagSet("cancel", flag.ExitOnError)
		jobID := fs.String("job-id", "", "job ID")
		fs.Parse(os.Args[2:])

		if _, err := client.CancelJob(ctx, &pb.CancelJobRequest{Id: *jobID}); err != nil {
			fmt.Printf("cancel err : %v \n", err)
		}

	default:
		fmt.Println("Invalid command or command format")

	}

}
