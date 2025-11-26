package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int, ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Duration(id) * time.Second)
				}
				fmt.Printf("%d: hello world\n", id)
			}
		}(i, ctx)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	<-sigCh   // получили сигнал
	cancel()  // отменяем context
	wg.Wait() // ждём завершения воркера
	fmt.Println("shutdown complete")
}
