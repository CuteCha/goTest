package faissUsage
/*
// C 标志io头文件，你也可以使用里面提供的函数
#include <stdio.h>

void pri(){
	printf("hey");
}

int add(int a,int b){
	return a+b;
}
*/
import "C"
import "fmt"

func Test() {
	fmt.Println(C.add(2, 1))
	C.pri()
}
