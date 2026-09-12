package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func getStudent(conn *pgx.Conn) {

	var id int
	var name string
	var age int

	err := conn.QueryRow(
		context.Background(),
		"SELECT id, name, age FROM students WHERE id = $1",
		1001,
	).Scan(&id, &name, &age)

	if err != nil {
		fmt.Println("Failed to get student:", err)
		return
	}

	fmt.Println("ID:", id)
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
}

func createStudent(conn *pgx.Conn) {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO students (id, name, age) VALUES ($1, $2, $3)",
		1002,
		"Ali",
		43,
	)

	if err != nil {
		fmt.Println("Failed to create student:", err)
		return
	}

	fmt.Println("Student created successfully")
}

func getStudents(conn *pgx.Conn) {
	rows, err := conn.Query(
		context.Background(),
		"SELECT id, name, age FROM students",
	)

	if err != nil {
		fmt.Println("Failed to get students:", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int

		err := rows.Scan(&id, &name, &age)

		if err != nil {
			fmt.Println("Failed to scan student:", err)
			return
		}

		fmt.Println("ID:", id)
		fmt.Println("Name:", name)
		fmt.Println("Age:", age)
		fmt.Println("---")
	}
}

func updateStudent(conn *pgx.Conn) {
	_, err := conn.Exec(
		context.Background(),
		"UPDATE students SET age = $1 WHERE id = $2",
		27,
		1001,
	)

	if err != nil {
		fmt.Println("Failed to update student:", err)
		return
	}

	fmt.Println("Student updated successfully")
}

func deleteStudent(conn *pgx.Conn) {
	_, err := conn.Exec(
		context.Background(),
		"DELETE FROM students WHERE id = $1",
		1002,
	)

	if err != nil {
		fmt.Print("Failed to delete Database: ", err)
		return
	}

	fmt.Println("Student deleted successfully")
}

func getStudentsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := connectDB()

	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(
		context.Background(),
		"SELECT id, name, age FROM students",
	)

	if err != nil {
		http.Error(w, "Failed to get students", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var students []Student

	for rows.Next() {
		var student Student

		err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Age,
		)

		if err != nil {
			http.Error(w, "Failed to scan student", http.StatusInternalServerError)
			return
		}

		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func main() {
	http.HandleFunc("/students", getStudentsHandler)

	fmt.Println("EduCore server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
