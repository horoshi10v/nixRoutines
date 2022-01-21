package main

import (
	"NIXRutine/models"
	"log"
	"sync"
)

type WorkerPool struct {
	Count  int
	Sender chan models.Restaurant
	Ender  chan bool
}

func NewWorkerPool(count int) *WorkerPool {
	return &WorkerPool{
		Count:  count,
		Sender: make(chan models.Restaurant, count),
		Ender:  make(chan bool),
	}
}

func (p *WorkerPool) Run(wg *sync.WaitGroup, handler func(author models.Restaurant)) {
	defer wg.Done()
	var rest models.Restaurant
	for {
		select {
		case rest = <-p.Sender:
			handler(rest)
		case <-p.Ender:
			log.Println("finish")
			return
		}
	}
}

func (p *WorkerPool) Stop() {
	for i := 0; i < p.Count; i++ {
		p.Ender <- false
	}
	close(p.Sender)
	close(p.Ender)
}
