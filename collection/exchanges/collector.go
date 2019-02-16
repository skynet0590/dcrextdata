package exchanges

import (
	"fmt"
	"sync"
	"time"
)

type Collector struct {
	Retrievers []Retriever
}

func (collector Collector) CollectAtInterval(interval time.Duration, wg *sync.WaitGroup, quit chan struct{}) (chan []DataTick, chan error) {
	resultChan := make(chan []DataTick, 1)
	errChan := make(chan error, 1)

	wg.Add(1)
	go func(resultChan chan []DataTick, errChan chan error) {
		ticker := time.NewTicker(interval)
		for _, v := range collector.Retrievers {
			wg.Add(1)
			go retrieve(wg, v, resultChan, errChan)
		}
	loop:
		for {
			select {
			case <-ticker.C:
				for _, v := range collector.Retrievers {
					wg.Add(1)
					go retrieve(wg, v, resultChan, errChan)
				}
			case <-quit:
				fmt.Print("Stopping collector")
				ticker.Stop()
				break loop
			}
		}
		wg.Done()
	}(resultChan, errChan)

	return resultChan, errChan
}
func retrieve(wg *sync.WaitGroup, retriever Retriever, resultChan chan []DataTick, errChan chan error) {
	result, err := retriever.Retrieve()
	if result != nil {
		resultChan <- result
	}
	if err != nil {
		errChan <- err
	}
	wg.Done()
}

type Retriever interface {
	Retrieve() ([]DataTick, error)
}
