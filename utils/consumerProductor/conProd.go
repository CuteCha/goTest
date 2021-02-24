package consumerProductor

import (
	"fmt"
	"time"
)

func Producer(ch chan int) {
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Printf("ch: %d, gen: %d, date: %d\n", &ch, i, time.Now().UnixNano())
	}
	close(ch)
}

func Consumer(id int, ch chan int, done chan bool) {
	for {
		value, ok := <-ch
		if ok {
			fmt.Printf("id: %d, recv: %d, date: %d\n", id, value, time.Now().UnixNano())
		} else {
			fmt.Printf("id: %d, closed, date: %d\n", id, time.Now().UnixNano())
			break
		}
	}
	done <- true
}

func Test() {
	ch := make(chan int, 3)

	coNum := 2
	done := make(chan bool, coNum)
	for i := 1; i <= coNum; i++ {
		go Consumer(i, ch, done)
	}

	go Producer(ch)
	for i := 1; i <= coNum; i++ {
		<-done
	}

	//当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))
}
