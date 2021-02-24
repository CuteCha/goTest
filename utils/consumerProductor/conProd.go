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

var messages2 = []string{
	"The world itself's",
	"just one big hoax.",
	"Spamming each other with our",
	"running commentary of bullshit,",
	"masquerading as insight, our social media",
	"faking as intimacy.",
	"Or is it that we voted for this?",
	"Not with our rigged elections,",
	"but with our things, our property, our money.",
	"I'm not saying anything new.",
	"We all know why we do this,",
	"not because Hunger Games",
	"books make us happy,",
	"but because we wanna be sedated.",
	"Because it's painful not to pretend,",
	"because we're cowards.",
	"- Elliot Alderson",
	"Mr. Robot",
}

const producerCount2 int = 3
const consumerCount2 int = 3

var workers []*producers

type producers struct {
	myQ  chan string
	quit chan bool
	id   int
}

func execute(jobQ chan<- string, workerPool chan *producers, allDone chan<- bool) {
	for _, j := range messages2 {
		jobQ <- j
	}
	close(jobQ)
	for _, w := range workers {
		w.quit <- true
	}
	close(workerPool)
	allDone <- true
}

func produce2(jobQ <-chan string, p *producers, workerPool chan *producers) {
	for {
		select {
		case msg := <-jobQ:
			{
				workerPool <- p
				if len(msg) > 0 {
					fmt.Printf("Job \"%v\" produced by worker %v\n", msg, p.id)
				}
				p.myQ <- msg
			}
		case <-p.quit:
			return
		}
	}
}

func consume2(cIdx int, workerPool <-chan *producers) {
	for {
		worker := <-workerPool
		if msg, ok := <-worker.myQ; ok {
			if len(msg) > 0 {
				fmt.Printf("Message \"%v\" is consumed by consumer %v from worker %v\n", msg, cIdx, worker.id)
			}
		}
	}
}

func Test3() {
	jobQ := make(chan string)
	allDone := make(chan bool)
	workerPool := make(chan *producers)

	for i := 0; i < producerCount2; i++ {
		workers = append(
			workers,
			&producers{
				myQ:  make(chan string),
				quit: make(chan bool),
				id:   i,
			})
		go produce2(jobQ, workers[i], workerPool)
	}

	go execute(jobQ, workerPool, allDone)

	for i := 0; i < consumerCount2; i++ {
		go consume2(i, workerPool)
	}
	<-allDone
}
