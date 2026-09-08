package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func students(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		student := Student{
			ID:   1001,
			Name: "Muscab",
			Age:  25,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)

	case http.MethodPost:
		var student Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid JSON")
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Student Created:", student.Name)

	case http.MethodPut:
		var student Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid JSON")
			return
		}

		fmt.Fprintln(w, "Student Updated:", student.Name)

	case http.MethodDelete:
		fmt.Fprintln(w, "Student Deleted")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method Not Allowed")
	}
}

func main() {
	http.HandleFunc("/students", students)

	fmt.Println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}


/*
Another powershell
curl.exe http://localhost:8080/students

then 
$body = '{"name":"Ahmed","age":25}'
Invoke-RestMethod -Uri "http://localhost:8080/students" -Method POST -ContentType "application/json" -Body $body

$body = '{"name":"Muscab Ahmed","age":21}'
Invoke-RestMethod -Uri "http://localhost:8080/students" -Method PUT -ContentType "application/json" -Body $body

curl.exe -X DELETE http://localhost:8080/students

*/