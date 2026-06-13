package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GoProOrg/go-practice/internal/domain/order"
	"github.com/GoProOrg/go-practice/worker"
)

func main() {
	// fmt.Println("=== Order Processor - Level 2 (Context + Error Handling) ===")
	fmt.Println("=== Order Processor - Level 2 (Fan-In Pipeline & Poison Pill) ===")
	numWorkers := 3
	numOrders := 10

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ctrl+C handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("Shut down signal received...")
		cancel()
	}()
	// TODO: defines jobs & results chans
	jobs := make(chan order.Order, 20)
	results := make(chan string, 20)
	errCh := make(chan error, 10)
	var wg sync.WaitGroup

	// start workers
	for i := range numWorkers {
		wg.Add(1)
		go worker.OrderWorker(i+1, ctx, jobs, results, errCh, &wg)
	}
	// defines Producer (e.g. from Gin HTTP handler)
	// TODO: here we use loop to simulate
	for i := 1; i <= numOrders; i++ {
		ord := order.SeedOrder(i)
		select {
		case <-ctx.Done():
			goto shutdown
		case jobs <- ord:
			fmt.Printf("Produced: %s (Table: %d)\n", ord.ID, ord.Table)
		}
	}
	close(jobs) // no more orders

shutdown:
	go func() {
		wg.Wait()
		close(results)
		close(errCh)
	}()
	// collects results
	for {
		select {
		case res, ok := <-results:
			if !ok {
				goto done
			}
			fmt.Print(res)
		case err, ok := <-errCh:
			if ok {
				fmt.Printf("X %v\n", err)
			}
		case <-ctx.Done():
			goto done
		}
	}
done:
	fmt.Println("✅ Level 2 Practice Completed!")
}
