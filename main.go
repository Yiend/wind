package main

import (
	"github.com/wind/engine"
	"github.com/wind/scheduler"
	"github.com/wind/fiction_17k/parser"
)

func main() {
	//打开一个小说网页
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.Queued{},
		WorkerCount: 20,
	}
	e.Run(engine.Request{
		Url:       "http://all.17k.com/lib/book/2_0_0_0_0_0_1_0_1.html",
		ParseFunc: parser.ParseFreeBookAll,
	})
}
