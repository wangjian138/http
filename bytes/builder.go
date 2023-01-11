package main

import (
	"fmt"
	"strings"
)

func main() {
	builder := strings.Builder{}

	builder.WriteByte('a')
	builder.WriteByte('b')

	str := builder.String()
	fmt.Printf("str:%v\n", str)

	a1 := 0b000
	a2 := a1 ^ 0
	fmt.Printf("a2:%v\n", a2)
}
