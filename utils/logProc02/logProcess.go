package logProc02

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	read(rc chan string)
}

type Writer interface {
	write(wc chan string)
}

type ReadFromFile struct {
	path string
}

func (r *ReadFromFile) read(rc chan string) {
	line := "message"
	rc <- line
}

type WriteToInfluxDB struct {
	influxDBsn string
}

func (w *WriteToInfluxDB) write(wc chan string) {
	fmt.Println(<-wc)
}

type LogProcess struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}

func (l *LogProcess) readFromFile() {
	line := "message"
	l.rc <- line
}

func (l *LogProcess) Process() {
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (l *LogProcess) writeToInfluxDB() {
	fmt.Println(<-l.wc)
}

func Test() {
	r := &ReadFromFile{
		path: "./logs/log",
	}

	w := &WriteToInfluxDB{
		influxDBsn: "user_password..",
	}

	lp := &LogProcess{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.read(lp.rc)
	go lp.Process()
	go lp.write.write(lp.wc)

	time.Sleep(500 * time.Millisecond)

}
