package main

import "fmt"

type Student struct {
	ID    int
	Name  string
	Age   int
	Email string
}

func updateStudent(student *Student, name string, age int) {
	if student == nil {
		fmt.Println("Student not found")
		return
	}

	student.Name = name
	student.Age = age

	fmt.Println("Student updated successfully")
}

func main() {
	student := Student{
		ID:    1001,
		Name:  "Ahmed",
		Age:   25,
		Email: "ahmed@gmail.com",
	}

	fmt.Println("Before:")
	fmt.Println("Name:", student.Name)
	fmt.Println("Age:", student.Age)

	updateStudent(&student, "Mohamed", 22)

	fmt.Println()
	fmt.Println("After:")
	fmt.Println("Name:", student.Name)
	fmt.Println("Age:", student.Age)

	fmt.Println()

	updateStudent(nil, "Ali", 20)
}