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
	id             int
	processingTask bool
	m              sync.Mutex
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
		if p.AllWorkersAreDone() {
			close(p.tasks)
			break
		}
	}

	wg.Wait()
}

func (p *WorkerPool) AllWorkersAreDone() bool {
	for _, worker := range p.workers {
		if worker.IsProcessingTask() {
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

		w.SetProcessingTask(true)
		processTaskFunc(task)
		w.SetProcessingTask(false)
	}
}

func (w *Worker) SetProcessingTask(status bool) {
	w.m.Lock()
	defer w.m.Unlock()
	w.processingTask = status
}

func (w *Worker) IsProcessingTask() bool {
	w.m.Lock()
	defer w.m.Unlock()
	return w.processingTask
}
