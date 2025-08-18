package main

import (
	"fmt"
	"todolist/internal"
)

func main() {
	var todolist internal.Todolist
	// todolist.Create("swimming")
	// todolist.Create("walk the dog")
	// todolist.Create("go shopping")
	// todolist.Complete(1)
	todolist.ReadFromFile("monday.tdl.json")

	for i, task := range todolist {
		fmt.Println(i, task.Name, task.Completed)
	}
	// todolist.Edit(2, "get groceries")
	// fmt.Println("-----------------")
	// todolist.Create("shop")
	// for i, task := range todolist {
	// 	fmt.Println(i, task.Name, task.Completed)
	// }
	// fmt.Println("-----------------")
	// todolist.Delete(2)
	// for i, task := range todolist {
	// 	fmt.Println(i, task.Name, task.Completed)
	// }
	// todolist.Store("monday.tdl.json")
}