package student

type Student struct {
	ID   int
	Name string
	Age  int
}

func GetStudent() Student {
	return Student{
		ID:   1001,
		Name: "Ahmed",
		Age:  25,
	}
}