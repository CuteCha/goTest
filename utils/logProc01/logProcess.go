package logProc01

import (
	"fmt"
	"strings"
	"time"
)

type LogProcess struct {
	path       string
	influxDBsn string
	rc         chan string
	wc         chan string
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
	lp := &LogProcess{
		rc:         make(chan string),
		wc:         make(chan string),
		path:       "./logs/log",
		influxDBsn: "user_password..",
	}

	go lp.readFromFile()
	go lp.Process()
	go lp.writeToInfluxDB()

	time.Sleep(500 * time.Millisecond)

}
