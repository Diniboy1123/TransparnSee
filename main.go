package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/Diniboy1123/transparnsee/config"
	"github.com/Diniboy1123/transparnsee/internal"
	ct "github.com/google/certificate-transparency-go/client"
	ctjson "github.com/google/certificate-transparency-go/jsonclient"
)

func main() {
	config.LoadConfig()

	startIndex := flag.Int64("start-index", 0, "Starting index for fetching log entries")
	workerCount := flag.Int("workers", 5, "Number of workers to process log entries")
	flag.Parse()

	client, err := ct.New(config.AppConfig.CtLogURL, nil, ctjson.Options{})
	if err != nil {
		log.Fatalf("Failed to create CT log client: %v", err)
	}

	entryCh := make(chan [2]int64, *workerCount)
	resultCh := make(chan string, 100)
	var wg sync.WaitGroup
	var lastProcessed = *startIndex

	for i := 0; i < *workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for indices := range entryCh {
				processed := internal.ProcessEntries(client, indices[0], indices[1], resultCh)
				atomic.AddInt64(&lastProcessed, int64(processed))
			}
		}()
	}

	go internal.WriteResultsToFile(config.AppConfig.OutputFile, resultCh)

	logInfo, err := client.GetSTH(context.Background())
	if err != nil {
		log.Fatalf("Failed to get log info: %v", err)
	}
	totalEntries := int64(logInfo.TreeSize)

	go internal.DisplayProgress(&lastProcessed, totalEntries)

	fmt.Println("Fetching certificates...")
	for start := *startIndex; start < totalEntries; start += config.AppConfig.BatchSize {
		end := start + config.AppConfig.BatchSize - 1
		if end >= totalEntries {
			end = totalEntries - 1
		}
		entryCh <- [2]int64{start, end}
	}
	close(entryCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nGracefully shutting down...")
		close(resultCh)
		wg.Wait()
		os.Exit(0)
	}()

	wg.Wait()
	close(resultCh)
	fmt.Println("Processing complete.")
}
