package worker

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Truongle68/go-practice/internal/domain/order"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// TODO: write the order-worker logic
func OrderWorker(id int, ctx context.Context, jobs <-chan order.Order, results chan<- string, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down gracefully...\n", id)
			return
		case ord, ok := <-jobs:
			if !ok {
				return
			}

			time.Sleep(time.Duration(200+id*300) * time.Millisecond)

			if random.Float64() < 0.1 {
				errCh <- fmt.Errorf("worker %d: inventory check failed for %s", id, ord.ID)
				continue
			}

			ord.CalculateTotal()
			results <- fmt.Sprintf("Worker %d Order %s for Table %d | Total: %.0f VND\n", id, ord.ID, ord.Table, ord.Total)
		}
	}
}
