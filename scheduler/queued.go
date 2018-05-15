package scheduler

import "github.com/wind/engine"

//Request队列,Worker队列
type Queued struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request //chan套chan
}

func (s *Queued) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (s *Queued) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (s *Queued) Submit(r engine.Request) {
	s.requestChan <- r
}

func (s *Queued) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var (
			requestQ []engine.Request
			workerQ  []chan engine.Request
		)
		for {
			var (
				activeRequest engine.Request
				activeWorker  chan engine.Request
			)
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
