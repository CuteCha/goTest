package faissUsage

/*
#include "cFile.cpp"
*/
import "C"

import "fmt"

func GoSum(a, b int) {
	s := C.sum(C.int(a), C.int(b))
	fmt.Println(s)
}

func Test2() {
	GoSum(4, 5)
}
