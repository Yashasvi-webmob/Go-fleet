package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/Yashasvi-webmob/gofleet/internal/scheduler"
)

func main() {
	repo := job.NewInMemoryRepository()
	service := job.NewService(repo)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	newJ, err := service.Submit(ctx, "redis-cli", 8, "high")
	if err != nil {
		fmt.Printf("Error service.Submit() server/main.go: %v \n", err)
		return
	}

	workerCount := 2
	sched := scheduler.New(service, 5, workerCount)
	sched.Start(ctx, workerCount)
	sched.Enqueue(newJ.ID)

	<-ctx.Done()
	fmt.Println("Shutdown SIG received")

	sched.Close()
	sched.Wait()

	fmt.Print("Shutdown Complete")

}
