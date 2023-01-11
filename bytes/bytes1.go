package main

import (
	"bytes"
	"fmt"
)

func main1() {
	strByte := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'v', 'd', 'b'}
	b1 := bytes.NewBuffer(strByte)
	fmt.Printf("b1:%v cap:%v\n", b1.Bytes(), cap(strByte))

	s, err := b1.ReadString(byte('c'))
	fmt.Printf("s:%v err:%v\n", s, err)

	//s1, n := b1.ReadString(byte('b'))
	//fmt.Printf("s:%v n:%v\n", s1, n)

	err = b1.UnreadByte()
	fmt.Printf("err:%v\n", err)

	str, size, err := b1.ReadRune()
	fmt.Printf("str:%v size:%v err:%v\n", str, size, err)

	b5 := b1.Next(2)
	fmt.Printf("b5:%v\n", b5)

	s2 := make([]byte, 3)
	n, err := b1.Read(s2)
	fmt.Printf("n:%v s2:%v err:%v cap:%v\n", n, s2, err, cap(s2))

	n, err = b1.WriteRune(100)
	fmt.Printf("n:%v err:%v byte:%v\n", n, err, b1.Bytes())

	err = b1.WriteByte('g')
	fmt.Printf("n:%v err:%v byte:%v\n", n, err, b1.Bytes())

	b2 := bytes.NewBuffer([]byte("a"))
	n1, err := b1.WriteTo(b2)
	fmt.Printf("n1:%v err:%v byte:%v b1:%v\n", n1, err, b2.Bytes(), b1.Bytes())

	n3, err := b1.ReadFrom(b2)
	fmt.Printf("n3:%v err:%v byte:%v b1:%v\n", n3, err, b2.Bytes(), b1.Bytes())

}
