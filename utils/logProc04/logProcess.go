package logProc04

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	curTime time.Time
	domain  string
	level   int
	score   float64
}

type Reader interface {
	read(rc chan []byte)
}

type Writer interface {
	write(wc chan *Message)
}

type ReadFromFile struct {
	path string
}

func (r *ReadFromFile) read(rc chan []byte) {
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("os.Open file fail: %s", err.Error()))
	}

	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("rd.ReadBytes file fail: %s", err.Error()))
		}
		rc <- line[:len(line)-1]
	}

}

type WriteToInfluxDB struct {
	influxDBsn string
}

func (w *WriteToInfluxDB) write(wc chan *Message) {
	for data := range wc {
		fmt.Println(data)
		fmt.Println(data.curTime)
		fmt.Printf("%+v\n", data)
	}
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}

func (l *LogProcess) Process() {
	//r := regexp.MustCompile(`([\d\.]+)\s+([a-z]+)\s+([\d\.])`)
	for data := range l.rc {
		v := string(data)
		//ret := r.FindStringSubmatch(v)
		ret := strings.Split(v, "|")
		if len(ret) != 4 {
			log.Printf("r.FindStringSubmatch fail: %s; %d\n", v, len(ret))
			continue
		}

		message := &Message{}

		timeFormat := "2006-01-02 15:04:05 +0000"
		loc, _ := time.LoadLocation("Asia/Shanghai")
		ct, err := time.ParseInLocation(timeFormat, ret[0], loc)
		if err != nil {
			log.Printf("time.ParseInLocation fail: %s; %s\n", ret[0], err.Error())
		}
		message.curTime = ct

		level, err := strconv.Atoi(ret[1])
		if err != nil {
			log.Printf("strconv.Atoi fail: %s; %s\n", ret[1], err.Error())
		}
		message.level = level

		message.domain = ret[2]

		score, err := strconv.ParseFloat(ret[3], 64)
		if err != nil {
			log.Printf("strconv.ParseFloat fail: %s; %s\n", ret[3], err.Error())
		}
		message.score = score

		l.wc <- message
	}
}

func Test() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("start")
	log.Printf("go %s", "process")
	r := &ReadFromFile{
		path: "/Users/cxq/go/src/debugProm/logs/log",
	}

	w := &WriteToInfluxDB{
		influxDBsn: "user_password..",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}

	go lp.read.read(lp.rc)
	go lp.Process()
	go lp.write.write(lp.wc)

	time.Sleep(30 * time.Second)

}
