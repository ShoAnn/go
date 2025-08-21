package internal

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Task struct {
	Name		string
	Completed	bool
	CreatedAt	time.Time
}

type Todolist []Task

func (t *Todolist) Store(filePath string) error {
	//encode todolist to json -> write to file
	encodedList, err := json.Marshal(t) 
	if err != nil {
		return errors.New("json encoding failed")
	}
	writeError := os.WriteFile(filePath, encodedList, 0644)
	if writeError != nil {
		return errors.New("file write failed")
	}
	return nil
}

func (t *Todolist) Create(taskName string) {
	newTask := Task{
		Name: taskName,
		Completed: false,
		CreatedAt: time.Now(),
	}
	*t = append(*t, newTask)
}

func (t *Todolist) ReadFromFile(filePath string) error {
	fileContent, err := os.ReadFile(filePath)
	// read file
	if err != nil {
		return errors.New("import failed")
	}
	return  json.Unmarshal(fileContent, t)
}

func (t *Todolist) Complete(taskIndex int) error {
	if taskIndex < 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	}
	(*t)[taskIndex].Completed = true
	return nil
}

func (t *Todolist) Edit(taskIndex int, newName string) error {
	if taskIndex <= 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	} else if newName == "" {
		return errors.New("invalid task name")
	}
	(*t)[taskIndex].Name = newName
	return nil
}


func (t *Todolist) Delete(taskIndex int) error {

	if taskIndex <= 0 || taskIndex > len(*t) {
		return errors.New("invalid task id")
	}
	*t = append((*t)[:taskIndex], (*t)[taskIndex+1:]...)
	return nil

}