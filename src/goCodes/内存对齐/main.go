package main

import (
	"fmt"
	"unsafe"
)

type Part1 struct{
	a bool
	b int32
	c int8
	d int64
	e byte
}

type Part2 struct{
	a bool
	e byte
	c int8
	b int32
	d int64
}

func main(){
	fmt.Println("bool size:",unsafe.Sizeof(bool(true)))
	fmt.Println("int32 size:",unsafe.Sizeof(int32(0)))
	fmt.Println("int8 size:",unsafe.Sizeof(int8(0)))
	fmt.Println("int64 size:",unsafe.Sizeof(int64(0)))
	fmt.Println("byte size:",unsafe.Sizeof(byte(0)))
	fmt.Println("string size:",unsafe.Sizeof(string("abced")))
	fmt.Println()
	// unsafe.Alignof() 是可以得到对齐系数
	fmt.Println("bool align:",unsafe.Alignof(bool(true)))
	fmt.Println("int32 align:",unsafe.Alignof(int32(0)))
	fmt.Println("int8 align:",unsafe.Alignof(int8(0)))
	fmt.Println("int64 align:",unsafe.Alignof(int64(0)))
	fmt.Println("byte align:",unsafe.Alignof(byte(0)))
	fmt.Println("string align:",unsafe.Alignof(string("abced")))

	part1 := Part1{}
	fmt.Println("Part1 size:",unsafe.Sizeof(part1))
	part2 := Part2{}
	fmt.Println("Part2 size:",unsafe.Sizeof(part2))
}