package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/Yashasvi-webmob/gofleet/internal/job"
)

func main() {
	repo := job.NewInMemoryRepository()
	service := job.NewService(repo)
	ctx := context.Background()
	switch os.Args[1] {
	case "submit":
		targetArr := os.Args[2:]
		if !slices.Contains(targetArr, "--cmd") {
			fmt.Errorf("cmd flag not set for submit type")
			return
		}
		service.Submit(ctx, os.Args[3], 5, "high")
		break
	case "status":
		targetArr := os.Args[2:]
		if !slices.Contains(targetArr, "--job-id") {
			fmt.Errorf("cmd flag not set for submit type")
			return
		}
		got, err := service.GetJob(ctx)
		if err != nil {
			fmt.Errorf("Job with flag job-id not found: %v", os.Args[3])
		}
		return got
		break
	case "cancel":
		fs := flag.NewFlagSet("cancel", flag.ExitOnError)
		jobID := fs.String("job-id","", "job ID")
		fs.Parse()
		break
	}

}
