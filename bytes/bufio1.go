package main

import (
	"bytes"
	"fmt"
	"learn/http/go/bufio"
)

func main2() {
	strByte := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'v', 'd', 'm'}
	b1 := bytes.NewBuffer(strByte)

	r1 := bufio.NewReaderSize(b1, 1000)

	fmt.Printf("strByte:%v bufio Size:%v  buffer:%v\n", strByte, r1.Size(), r1.Buffered())

	b2 := make([]byte, 10)
	b2[0] = 2
	b2 = append(b2, 1)
	b3 := []byte{'a', 'b', 'c', 'd', 'e', 'f'}

	n := copy(b2, b3[1:])
	fmt.Printf("n:%v b2:%v b3:%v\n", n, b2, b3)

	d1, err := r1.Peek(2)
	fmt.Printf("d1:%v err:%v\n", d1, err)

	dis1, err := r1.Discard(1)
	fmt.Printf("dis1:%v err:%v\n", dis1, err)

	d2, err := r1.Peek(2)
	fmt.Printf("d2:%v err:%v\n", d2, err)

	a1 := make([]byte, 10)

	n2, err := r1.Read(a1)
	fmt.Printf("n2:%v a1:%v err:%v\n", n2, a1, err)

	r1 = bufio.NewReaderSize(bytes.NewBuffer(strByte), 1000)
	//r1.Peek(4)
	//l, err := r1.ReadSlice('B')
	//fmt.Printf("l:%v err:%v\n", string(l), err)

	l, err := r1.ReadBytes('d')
	fmt.Printf("l:%v err:%v\n", string(l), err)

}
