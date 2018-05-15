package engine

import (
	"log"
	"github.com/wind/fetcher"
)

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

//调度程序
type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

func (e *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan ParseResult)
	e.Scheduler.Run()
	//建worker
	for i := 0; i < e.WorkerCount; i++ {
		creatWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	//接收out
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item: %v", item)
		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

func creatWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			//TODO tell scheduler i'm ready
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

func worker(r Request) (ParseResult, error) {
	log.Printf("Fetching  %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching URL %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	return r.ParseFunc(body), err
}
