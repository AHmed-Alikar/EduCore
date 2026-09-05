package main

import "fmt"

type Address struct {
	City    string
	Country string
}

type Student struct {
	ID      int
	Name    string
	Age     int
	Email   string
	Address Address
}

func introduce(student Student) {
	fmt.Println("My name is", student.Name)
}

func (s Student) introduce() {
	fmt.Println("My name is:", s.Name)
}

func (s *Student) updateAge(age int) {
	s.Age = age
}

func main() {
	student := Student{
		ID:    1001,
		Name:  "Ahmed",
		Age:   25,
		Email: "ahmed@gmail.com",
		Address: Address{
			City:    "Mogadishu",
			Country: "Somalia",
		},
	}

	introduce(student)
	student.introduce()

	student.updateAge(26)

	fmt.Println("Student Age:", student.Age)
}