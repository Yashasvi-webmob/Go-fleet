package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/Yashasvi-webmob/gofleet/gen"
	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/Yashasvi-webmob/gofleet/internal/scheduler"
	transport "github.com/Yashasvi-webmob/gofleet/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	repo := job.NewInMemoryRepository()
	service := job.NewService(repo)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	workerCount := 2
	sched := scheduler.New(service, 5, workerCount)
	sched.Start(ctx, workerCount)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, transport.NewJobServer(service, sched))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("Grpc never stopped")
		}
	}()
	// newJ, err := service.Submit(ctx, "redis-cli", 8, "high")
	// if err != nil {
	// 	fmt.Printf("Error service.Submit() server/main.go: %v \n", err)
	// 	return
	// }

	// workerCount := 2
	// sched := scheduler.New(service, 5, workerCount)
	// sched.Start(ctx, workerCount)
	// sched.Enqueue(newJ.ID)

	<-ctx.Done()
	fmt.Println("Shutdown SIG received")

	grpcServer.GracefulStop()
	sched.Close()
	sched.Wait()

	fmt.Print("Shutdown Complete")

}
