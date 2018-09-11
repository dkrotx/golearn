package main

import (
	"log"
	"net/rpc"
)

type ToDo struct {
	Title, Status string
}

type EditToDo struct {
	Title, NewTitle, NewStatus string
}

type ToDoListReply struct {
    TaskList []ToDo
}

func main() {
    var err error
	var reply ToDo

	// Create a TCP connection to localhost on port 1234
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	finishApp := ToDo{"Finish App", "Started"}
	makeDinner := ToDo{"Make Dinner", "Not Started"}
	walkDog := ToDo{"Walk the dog", "Not Started"}

	client.Call("Task.MakeToDo", finishApp, &reply)
	client.Call("Task.MakeToDo", makeDinner, &reply)
	client.Call("Task.MakeToDo", walkDog, &reply)

	/*client.Call("Task.DeleteToDo", makeDinner, &reply)

	client.Call("Task.MakeToDo", makeDinner, &reply)

	client.Call("Task.GetToDo", "Finish App", &reply)
	log.Println("Finish App: ", reply)

	err = client.Call("Task.EditToDo", EditToDo{"Finish App", "Finish App", "Completed"}, &reply)
	if err != nil {
		log.Fatal("Problem editing ToDo: ", err)
	}*/

    var listReply ToDoListReply
    err = client.Call("Task.ListToDo", 10, &listReply)
	if err != nil {
		log.Fatal("Problem with listing ToDos: ", err)
	}

    log.Printf("Todos: %v\n", listReply)
}
