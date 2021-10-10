package main

import (
	"fmt"
	"greetings"
)


func main() {
	message := greetings.Hello("Felps")
	fmt.Println(message)
}