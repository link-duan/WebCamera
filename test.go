package main

type T struct {
	name string
}

func main() {
	name := "hello"
	str := `{"name":"`+name+`"}`
	println(str)
}