package main

import "fmt"

type Student struct {
	ID   int
	Name string
	Age  int
}

func main() {
	ch := make(chan Student)

	go func() {
		ch <- Student{
			ID:   1001,
			Name: "Muscab",
			Age:  20,
		}
	}()

	go func() {
		ch <- Student{
			ID:   2002,
			Name: "Ali",
			Age:  43,
		}
	}()

	students1 := <-ch
	students2 := <-ch

	fmt.Println("Name 1:", students1.Name)
	fmt.Println("Name 2:", students2.Name)
}