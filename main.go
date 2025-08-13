package main

import (
	"errors"
	"time"
)

type Task struct {
	Name		string
	Completed	bool
	CreatedAt	time.Time
}

type TaskSlice []Task

func (t *TaskSlice) Create(taskName string) {
	newTask := Task{
		Name: taskName,
		Completed: false,
		CreatedAt: time.Now(),
	}
	*t = append(*t, newTask)
	// return nil
}

func (t *TaskSlice) Complete(taskIndex int) error {
	if taskIndex <= 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	}
	(*t)[taskIndex].Completed = true
	return nil
}

func (t *TaskSlice) Edit(taskIndex int, newName string) error {
	if taskIndex <= 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	} else if newName == "" {
		return errors.New("invalid task name")
	}
	(*t)[taskIndex].Name = newName
	return nil
}

func (t *TaskSlice) Delete(taskIndex int) error {
	if taskIndex <= 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	}
	*t = append((*t)[:taskIndex-1], (*t)[taskIndex:]...)
	return nil
}

func main() {
}