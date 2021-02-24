package consumerProductor

import (
	"fmt"
	"sync"
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

var messages = [][]string{
	{
		"The world itself's",
		"just one big hoax.",
		"Spamming each other with our",
		"running commentary of bullshit,",
	},
	{
		"but with our things, our property, our money.",
		"I'm not saying anything new.",
		"We all know why we do this,",
		"not because Hunger Games",
		"books make us happy,",
	},
	{
		"masquerading as insight, our social media",
		"faking as intimacy.",
		"Or is it that we voted for this?",
		"Not with our rigged elections,",
	},
	{
		"but because we wanna be sedated.",
		"Because it's painful not to pretend,",
		"because we're cowards.",
		"- Elliot Alderson",
		"Mr. Robot",
	},
}

const producerCount int = 4
const consumerCount int = 3

func produce(link chan<- string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, msg := range messages[id] {
		link <- msg
	}
}

func consume(link <-chan string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range link {
		fmt.Printf("Message \"%v\" is consumed by consumer %v\n", msg, id)
	}
}

func Test2() {
	link := make(chan string)
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}

	wp.Add(producerCount)
	wc.Add(consumerCount)

	for i := 0; i < producerCount; i++ {
		go produce(link, i, wp)
	}

	for i := 0; i < consumerCount; i++ {
		go consume(link, i, wc)
	}

	wp.Wait()
	close(link)
	wc.Wait()
}
