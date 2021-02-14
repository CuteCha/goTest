package logProc05

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	totalNums int
	tps       float64
	rcLen     int
	wcLen     int
	errNums   int
	runTime   string
}

const (
	TOTAL = 0
	ERROR = 1
)

var monitorCh = make(chan int, 200)

type Monitor struct {
	sTime  time.Time
	data   Info
	tipSeq []int
}

func (m *Monitor) run(lp *LogProcess) {
	go func() {
		for v := range monitorCh {
			switch v {
			case TOTAL:
				m.data.totalNums += 1
			case ERROR:
				m.data.errNums += 1

			}
		}
	}()

	tick := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-tick.C
			m.tipSeq = append(m.tipSeq, m.data.totalNums)
			if len(m.tipSeq) > 2 {
				m.tipSeq = m.tipSeq[1:]
			}

		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.data.runTime = time.Now().Sub(m.sTime).String()
		m.data.rcLen = len(lp.rc)
		m.data.wcLen = len(lp.wc)

		if len(m.tipSeq) >= 2 {
			m.data.tps = float64(m.tipSeq[1]-m.tipSeq[0]) / 5
		}

		ret, _ := json.MarshalIndent(m.data, "", "\t")

		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":9193", nil)
}

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
		monitorCh <- TOTAL
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
			monitorCh <- ERROR
			log.Printf("time.ParseInLocation fail: %s; %s\n", ret[0], err.Error())
		}
		message.curTime = ct

		level, err := strconv.Atoi(ret[1])
		if err != nil {
			monitorCh <- ERROR
			log.Printf("strconv.Atoi fail: %s; %s\n", ret[1], err.Error())
		}
		message.level = level

		message.domain = ret[2]

		score, err := strconv.ParseFloat(ret[3], 64)
		if err != nil {
			monitorCh <- ERROR
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

	monitor := &Monitor{
		sTime: time.Now(),
		data:  Info{},
	}
	monitor.run(lp)

}
