package main

import "fmt"

func main() {
	var(
		u8 uint8      = 255
		i8 int8       = 127
		f32 float32   = 0.2
		c64 complex64 = 15 + 12i
		s string      = "Hello, World!"
	)
	fmt.Println(u8, i8, f32, c64, s)
	fmt.Printf("%T %v", u8, u8)
}
