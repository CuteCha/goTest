package main

import (
	"debugProm/utils"
	"fmt"
)

func test() {
	fmt.Println("hello, let's go")

	messages := make(chan string)
	go sample(messages)
	msg := <-messages
	fmt.Println(msg)

}

func sample(messages chan string) {
	messages <- "ping"
}
func main() {
	test()
	utils.Run()
}
