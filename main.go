package main

import (
	"debugProm/utils/logProc03"
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
	logProc03.Test()
}
