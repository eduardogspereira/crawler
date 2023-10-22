package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestWorkerPool_Success(t *testing.T) {
	workerPool := NewWorkerPool(10)

	tasks := []string{"a", "b", "c", "d", "e", "f"}
	for _, task := range tasks {
		workerPool.AddTask(task)
	}

	var results []string
	var m sync.Mutex
	processTaskFunc := func(letter interface{}) {
		result := letter.(string) + letter.(string)
		m.Lock()
		results = append(results, result)
		m.Unlock()
		if result == "ee" {
			workerPool.AddTask(result)
		}
	}

	workerPool.ProcessTasks(processTaskFunc)
	assert.Contains(t, results, "aa")
	assert.Contains(t, results, "bb")
	assert.Contains(t, results, "cc")
	assert.Contains(t, results, "dd")
	assert.Contains(t, results, "ee")
	assert.Contains(t, results, "eeee")
}
