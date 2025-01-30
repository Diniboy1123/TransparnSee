package internal

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

func DisplayProgress(lastProcessed *int64, totalEntries int64) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("\rProcessed: %d/%d entries (%.2f%%)",
			atomic.LoadInt64(lastProcessed), totalEntries,
			(float64(atomic.LoadInt64(lastProcessed))/float64(totalEntries))*100)
	}
}

func WriteResultsToFile(outputFile string, resultCh <-chan string) {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	for domain := range resultCh {
		_, err := fmt.Fprintln(file, domain)
		if err != nil {
			log.Printf("Error writing to output file: %v", err)
		}
	}
}
