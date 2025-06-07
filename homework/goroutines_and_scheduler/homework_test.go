package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type taskHeap struct {
	tasks []*Task
}

func (th *taskHeap) len() int {
	return len(th.tasks)
}

func (th *taskHeap) swap(a, b int) {
	th.tasks[a], th.tasks[b] = th.tasks[b], th.tasks[a]
}

func (th *taskHeap) pop() *Task {
	task := th.tasks[0]
	th.tasks = th.tasks[1:]
	return task
}

func (th *taskHeap) isHigherPriority(a, b int) bool {
	return th.tasks[a].Priority > th.tasks[b].Priority
}

func (th *taskHeap) fixHeapAfterAdd(index int) {
	if index == 0 || th.isHigherPriority((index-1)/2, index) {
		return
	}
	th.swap((index-1)/2, index)
	th.fixHeapAfterAdd((index - 1) / 2)
}

func (th *taskHeap) heapify(index int) {
	if index < 0 {
		return
	}
	maxChild, leftChild, rightChild := index, index*2+1, index*2+2
	if leftChild < th.len() && th.isHigherPriority(leftChild, maxChild) {
		maxChild = leftChild
	}
	if rightChild < th.len() && th.isHigherPriority(rightChild, maxChild) {
		maxChild = rightChild
	}
	if maxChild != index {
		th.swap(maxChild, index)
	}
	th.heapify(index - 1)
}

type Scheduler struct {
	taskHeap
	taskMap map[int]*Task
}

func NewScheduler() Scheduler {
	return Scheduler{
		taskMap: make(map[int]*Task),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.taskHeap.tasks = append(s.tasks, &task)
	s.taskMap[task.Identifier] = &task
	s.fixHeapAfterAdd(s.len() - 1)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	task, ok := s.taskMap[taskID]
	if !ok {
		return
	}
	task.Priority = newPriority
	s.heapify(s.len() / 2)
}

func (s *Scheduler) GetTask() Task {
	task := s.taskHeap.pop()
	delete(s.taskMap, task.Identifier)
	s.heapify(0)
	return *task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100 // change priority to correct comparison

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
