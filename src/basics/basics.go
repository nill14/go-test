package main

import (
	"fmt"
	"math"
	"math/rand"
)

var c, python, java bool
var i, j int = 1, 2

func main() {
	fmt.Printf("hello, world\n")
	fmt.Println("My favorite number is", rand.Intn(10))
	fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))
	fmt.Println(math.Pi)
	fmt.Println(add(42, 13))

	a, b := swap("hello", "world")
	fmt.Println(a, b)

	fmt.Println(split(17))

	var i int
	fmt.Println(i, c, python, java)

	var c, python, java = true, false, "no!"
	k := 3
	fmt.Println(i, j, k, c, python, java)

	defaultValue()

	typeConversions()

	numericConstants()
}

func add(x, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func defaultValue() {
	var i int
	var f float64
	var b bool
	var s string
	fmt.Printf("%v %v %v %q\n", i, f, b, s)
}

func typeConversions() {
	i := 42
	f := float64(i)
	u := uint(f)
	var x, y int = 3, 4
	var g float64 = math.Sqrt(float64(x*x + y*y))
	var t uint = uint(g)
	fmt.Println(x, y, u, f, t)

	v := 0.867 + 0.5i // change me!
	fmt.Printf("v is of type %T\n", v)

	const Pi = 3.14
	const World = "世界"
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")

	const Truth = true
	fmt.Println("Go rules?", Truth)
}

func needInt(x int) int { return x*10 + 1 }
func needFloat(x float64) float64 {
	return x * 0.1
}

func numericConstants() {
	// Create a huge number by shifting a 1 bit left 100 places.
	// In other words, the binary number that is 1 followed by 100 zeroes.
	const Big = 1 << 100
	// Shift it right again 99 places, so we end up with 1<<1, or 2.
	const Small = Big >> 99
	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))
	//	fmt.Println(needInt(Big)) //constant 1267650600228229401496703205376 overflows int
}
