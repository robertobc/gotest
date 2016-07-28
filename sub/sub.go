package sub

import "fmt"

type MySub struct {
	A string
	B int
}

func (m *MySub) Hello() {
	fmt.Println("hello")
}
