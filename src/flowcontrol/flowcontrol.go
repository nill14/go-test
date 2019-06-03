package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	basicLoop()
	loop1()

	fmt.Println(sqrt(2), sqrt(-4))

	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)

	fmt.Println(Sqrt(4))
	//	fmt.Println(Sqrt(150.256846897e70))
	fmt.Println(Sqrt(1))

	Switch()

	days()

	timeOfTheDay()

	deferTest()
	deferCounting()
}

func basicLoop() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
}

func loop1() {
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)
}

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func Sqrt(x float64) float64 {
	z := float64(1)
	var diff_1, diff_2 float64
	for i := 0; i < 2 || (i < 1000 && math.Abs(diff_1) < math.Abs(diff_2)); i++ {
		diff_2 = diff_1
		diff_1 = (z*z - x) / (2 * z)
		z -= diff_1
		fmt.Printf("%5d %20.20g %20.20g\n", i, z, diff_1)
	}
	fmt.Printf("Expected sqrt(%g) = %g\n", x, math.Sqrt(x))
	return z
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func Switch() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}

func days() {
	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}
}

func timeOfTheDay() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

func deferTest() {
	defer fmt.Println("world")

	fmt.Println("hello")
}

func deferCounting() {
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("done")
}
