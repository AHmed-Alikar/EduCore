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

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := connectDB()

	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	var user User
	err = conn.QueryRow(
		context.Background(),
		`INSERT INTO users (name, email, password_hash, role)
	 VALUES ($1, $2, $3, $4)
	 RETURNING id, name, email, role`,
		input.Name,
		input.Email,
		hash,
		"student",
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := connectDB()

	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	defer conn.Close(context.Background())

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	var user User

	err = conn.QueryRow(
		context.Background(),
		"SELECT id, name, email, password_hash FROM users WHERE email = $1",
		input.Email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !checkPassword(input.Password, user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := createToken(user)

	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to your profile")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome Admin")
}
func main() {
	http.HandleFunc("/students", getStudentsHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.Handle("/profile", authMiddleware(http.HandlerFunc(profileHandler)))
	fmt.Println("EduCore server running on http://localhost:8080")
	http.Handle(
		"/admin",
		authMiddleware(
			requireRole(
				"admin",
				http.HandlerFunc(adminHandler),
			),
		),
	)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
