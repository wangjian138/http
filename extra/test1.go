package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test1Slice := []byte{'a', 'b', 'c'}
	fmt.Println(test1Slice, string(test1Slice))
	rand.Seed(time.Now().Unix())
	v, err := rand.Read(test1Slice)
	fmt.Printf("v:%v test1Slice:%v err:%v", v, test1Slice, err)

	permNum := rand.Perm(7)
	fmt.Printf("permNum:%v\n", permNum)
}
