package jsonParse

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name    string
	Age     int
	Vip     bool
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

type Stu struct {
	Name  string `json:"name"`
	Age   int
	HIgh  bool
	sex   string
	Class *Class `json:"class"`
}

type Class struct {
	Name  string
	Grade int
}

func Struct2Json() {
	//实例化一个数据结构，用于生成json字符串
	stu := Stu{
		Name: "张三",
		Age:  18,
		HIgh: true,
		sex:  "男",
	}

	//指针变量
	cla := new(Class)
	cla.Name = "1班"
	cla.Grade = 3
	stu.Class = cla

	//Marshal失败时err!=nil
	jsonStu, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}

	//jsonStu是[]byte类型，转化成string类型便于查看
	fmt.Println(string(jsonStu))
}

func Json2Struct() {
	//json字符中的"引号，需用\进行转义，否则编译出错
	//json字符串沿用上面的结果，但对key进行了大小的修改，并添加了sex数据
	//data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	data := `{"name":"张三","age":18,"high":true,"class":{"name":"1班","grade":3}}`
	str := []byte(data)

	//1.Unmarshal的第一个参数是json字符串，第二个参数是接受json解析的数据结构。
	//第二个参数必须是指针，否则无法接收解析的数据，如stu仍为空对象StuRead{}
	//2.可以直接stu:=new(StuRead),此时的stu自身就是指针
	stu := Stu{}
	err := json.Unmarshal(str, &stu)

	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(stu)
	fmt.Println(*stu.Class)
}

type User struct {
	//name string  //fmt.Println(string(userJSON)) 数据无法解析
	//age  int
	Name string
	Age  int
}

func ULTest() {
	user := User{"Tom", 18}
	userJSON, err := json.Marshal(user)

	if err == nil {
		fmt.Println(string(userJSON))
	}

	user2 := new(User)
	err = json.Unmarshal(userJSON, user2)

	if err == nil {
		fmt.Println(user2)
		fmt.Println(*user2)
	}
}
