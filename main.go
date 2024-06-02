package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	Id        int
	Text      string
	Completed bool
}

var data = map[string][]Todo{
	"todos": {
		Todo{1, "Go Htmx", true},
	},
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, data)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	text := r.PostFormValue("todo")
	templ := template.Must(template.ParseFiles("index.html"))
	todo := Todo{Id: len(data["todos"]) + 1, Text: text, Completed: false}
	data["todos"] = append(data["todos"], todo)

	templ.ExecuteTemplate(w, "todo-list-element", todo)
	fmt.Println("Added todo with text: ", text, todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	
	todos := data["todos"]
	var deletedTodo Todo
	for i, todo := range todos {
		if todo.Id == id {
			deletedTodo = todos[i]
			todos = append(todos[:i], todos[i+1:]...)
			data["todos"] = todos
			break
		}
	}

	templ := template.Must(template.ParseFiles("index.html"))
	templ.ExecuteTemplate(w, "todo-list-element", todos)

	fmt.Println("Deleted todo with id: ", id, deletedTodo)
}

func toggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	todos := data["todos"]
	var updatedTodo Todo
	for i, todo := range todos {
		if todo.Id == id {
			todos[i].Completed = !todos[i].Completed
			updatedTodo = todos[i]
			break
		}
	}

	templ := template.Must(template.ParseFiles("index.html"))
	templ.ExecuteTemplate(w, "todo-list-element", updatedTodo)

	fmt.Println("Toggled todo with id: ", id, updatedTodo)
}

// Only for testing
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {

	http.HandleFunc("/", todoHandler)
	http.HandleFunc("/get-todos", getTodosHandler)
	http.HandleFunc("/add-todo", addTodoHandler)
	http.HandleFunc("/delete-todo", deleteTodoHandler)
	http.HandleFunc("/toggle-todo", toggleTodoHandler)

	fmt.Println("Listening on port 3535")
	err := http.ListenAndServe(":3535", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
