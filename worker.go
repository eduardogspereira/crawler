package main

import (
	"sync"
	"time"
)

type WorkerPool struct {
	tasks   chan interface{}
	workers []*Worker
}

type Worker struct {
	id        int
	isWorking bool
	m         sync.Mutex
}

func NewWorkerPool(numOfWorkers int) *WorkerPool {
	var workers []*Worker

	for workerID := 0; workerID < numOfWorkers; workerID++ {
		workers = append(workers, &Worker{id: workerID})
	}

	return &WorkerPool{workers: workers, tasks: make(chan interface{}, 1_000)}
}

func (p *WorkerPool) AddTask(task interface{}) {
	p.tasks <- task
}

func (p *WorkerPool) ProcessTasks(processTaskFunc func(interface{})) {
	wg := sync.WaitGroup{}

	for _, worker := range p.workers {
		wg.Add(1)
		go worker.Work(p.tasks, processTaskFunc, &wg)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		if p.AllTasksProcessed() {
			close(p.tasks)
			break
		}
	}

	wg.Wait()
}

func (p *WorkerPool) AllTasksProcessed() bool {
	for _, worker := range p.workers {
		if worker.IsWorking() {
			return false
		}
	}

	return true
}

func (w *Worker) Work(tasks chan interface{}, processTaskFunc func(interface{}), wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, ok := <-tasks
		if !ok {
			return
		}

		w.SetIsWorking(true)
		processTaskFunc(task)
		w.SetIsWorking(false)
	}
}

func (w *Worker) SetIsWorking(status bool) {
	w.m.Lock()
	defer w.m.Unlock()
	w.isWorking = status
}

func (w *Worker) IsWorking() bool {
	w.m.Lock()
	defer w.m.Unlock()
	return w.isWorking
}
