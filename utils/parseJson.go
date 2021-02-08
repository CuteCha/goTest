package utils

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name    string
	Age     int
	Vip   bool
	Classes []string
	Price   float32
}

func (s *Student) ShowStu() {
	fmt.Println("show Student:")
	fmt.Println("\tName: ", s.Name)
	fmt.Println("\tAge: ", s.Age)
	fmt.Println("\tVip: ", s.Vip)
	fmt.Println("\tPrice: ", s.Price)
	fmt.Printf("\tClasses: ")
	for _, a := range s.Classes {
		fmt.Printf("%s ", a)
	}
	fmt.Println("")
}

func Run() {
	st := &Student{
		"Xiao Ming",
		16,
		true,
		[]string{"Math", "English", "Chinese"},
		9.99,
	}
	fmt.Println("before JSON encoding :")
	st.ShowStu()

	b, err := json.Marshal(st)
	if err != nil {
		fmt.Println("encoding failed")
	} else {
		fmt.Println("encoded data : ")
		fmt.Println(b)
		fmt.Println(string(b))
	}

	ch := make(chan string, 1)
	go func(c chan string, str string) {
		c <- str
	}(ch, string(b))

	strData := <-ch
	fmt.Println("--------------------------------")
	stb := &Student{}
	stb.ShowStu()
	err = json.Unmarshal([]byte(strData), &stb)
	if err != nil {
		fmt.Println("Unmarshal failed")
	} else {
		fmt.Println("Unmarshal success")
		stb.ShowStu()
	}
}
