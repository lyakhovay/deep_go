package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	binaryHeap []Task
}

func NewScheduler() Scheduler {
	return Scheduler{}
}

func (s *Scheduler) len() int {
	return len(s.binaryHeap)
}

func (s *Scheduler) swap(a, b int) {
	s.binaryHeap[a], s.binaryHeap[b] = s.binaryHeap[b], s.binaryHeap[a]
}

func (s *Scheduler) isHigherPriority(a, b int) bool {
	return s.binaryHeap[a].Priority > s.binaryHeap[b].Priority
}

func (s *Scheduler) AddTask(task Task) {
	s.binaryHeap = append(s.binaryHeap, task)
	s.fixHeapAfterAdd(s.len() - 1)
}

func (s *Scheduler) fixHeapAfterAdd(index int) {
	if index == 0 || s.isHigherPriority((index-1)/2, index) {
		return
	}
	s.swap((index-1)/2, index)
	s.fixHeapAfterAdd((index - 1) / 2)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	for i := range s.binaryHeap {
		if s.binaryHeap[i].Identifier == taskID {
			s.binaryHeap[i].Priority = newPriority
			break
		}
	}
	s.heapify(s.len() / 2)
}

func (s *Scheduler) heapify(index int) {
	if index < 0 {
		return
	}
	maxChild, leftChild, rightChild := index, index*2+1, index*2+2
	if leftChild < s.len() && s.isHigherPriority(leftChild, maxChild) {
		maxChild = leftChild
	}
	if rightChild < s.len() && s.isHigherPriority(rightChild, maxChild) {
		maxChild = rightChild
	}
	if maxChild != index {
		s.swap(maxChild, index)
	}
	s.heapify(index - 1)
}

func (s *Scheduler) GetTask() Task {
	task := s.binaryHeap[0]
	s.binaryHeap = s.binaryHeap[1:]
	s.heapify(0)
	return task
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
