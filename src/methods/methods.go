package main

import (
	"fmt"
	"image"
	"io"
	"math"
	"strings"
	"time"
)

type Vertex struct {
	X, Y float64
}

type MyFloat float64

type Abser interface {
	Abs() float64
}

func main() {
	receivers()
	receiversNonStruct()
	pointerReceivers()
	interfaces()
	interfaceValues()
	typeAssertions()
	typeSwitches()
	stringers()
	errors()
	exerciseErrors()
	readers()
	exerciseRotReader()
	images()

}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func receivers() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
	fmt.Println(Abs(v))
}

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func receiversNonStruct() {
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func Scale(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func pointerReceivers() {
	v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())

	v = Vertex{3, 4}
	Scale(&v, 10)
	fmt.Println(v.Abs())
}

func interfaces() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f  // a MyFloat implements Abser
	a = &v // a *Vertex implements Abser

	// In the following line, v is a Vertex (not *Vertex)
	// and does NOT implement Abser.
	a = v

	fmt.Println(a.Abs())
}

type I interface {
	Println()
}

type Text struct {
	S string
}

func (t *Text) Println() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

type Float float64

func (f Float) Println() {
	fmt.Println(f)
}

func interfaceValues() {
	var i I
	describe(i)
	//	i.Println() runtime error

	var j interface{} = "Hello"
	describe(j)

	var t *Text
	i = t
	describe(i)
	i.Println()

	i = &Text{"Hello"}
	describe(i)
	i.Println()

	i = Float(math.Pi)
	describe(i)
	i.Println()
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func typeAssertions() {
	var i interface{} = "hello"

	s := i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	f, ok := i.(float64)
	fmt.Println(f, ok)

	//	f = i.(float64) // panic
	//	fmt.Println(f)
}

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

func typeSwitches() {
	do(21)
	do("hello")
	do(true)
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

func stringers() {
	a := Person{"Arthur Dent", 42}
	z := Person{"Zaphod Beeblebrox", 9001}
	fmt.Println(a, z)
}

type IPAddr [4]byte

// Add a "String() string" method to IPAddr.
func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func stringers2() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}

func errors() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

type ErrNegativeSqrt float64

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

	z := float64(1)
	var diff_1, diff_2 float64
	for i := 0; i < 2 || (i < 1000 && math.Abs(diff_1) < math.Abs(diff_2)); i++ {
		diff_2 = diff_1
		diff_1 = (z*z - x) / (2 * z)
		z -= diff_1
	}
	return z, nil
}

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %g", float64(e))
}

func exerciseErrors() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2)) //returns "cannot Sqrt negative number: -2".
}

func readers() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}

func images() {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(m.Bounds())
	fmt.Println(m.At(0, 0).RGBA())
}
