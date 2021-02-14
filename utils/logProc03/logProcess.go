package logProc03

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	read(rc chan []byte)
}

type Writer interface {
	write(wc chan string)
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
			continue
		}
		rc <- line[:len(line)-1]
	}

}

type WriteToInfluxDB struct {
	influxDBsn string
}

func (w *WriteToInfluxDB) write(wc chan string) {
	for data := range wc {
		fmt.Println(data)
	}
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

func (l *LogProcess) Process() {
	for data := range l.rc {
		l.wc <- strings.ToUpper(string(data))
	}

}

func Test() {
	r := &ReadFromFile{
		path: "/Users/cxq/go/src/debugProm/logs/log",
	}

	w := &WriteToInfluxDB{
		influxDBsn: "user_password..",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.read(lp.rc)
	go lp.Process()
	go lp.write.write(lp.wc)

	time.Sleep(30 * time.Second)

}
