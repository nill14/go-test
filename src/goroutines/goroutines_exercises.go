package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch) // <- closes the channel when this function returns
	var walk func(t *tree.Tree)
	walk = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walk(t.Left)
		ch <- t.Value
		walk(t.Right)
	}
	walk(t)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if v1 != v2 || ok1 != ok2 {
			return false
		}

		if !ok1 {
			break
		}
	}

	return true
}

func printTree(t *tree.Tree) {
	ch := make(chan int)
	go Walk(t, ch)

	for v := range ch {
		fmt.Printf("%d ", v)
	}
	fmt.Println("")
}

func testTree() {
	t11 := &tree.Tree{nil, 1, nil}
	t12 := &tree.Tree{nil, 2, nil}
	t15 := &tree.Tree{nil, 5, nil}
	t113 := &tree.Tree{nil, 13, nil}

	t31 := &tree.Tree{t11, 1, t12}
	t38 := &tree.Tree{t15, 8, t113}
	t73 := &tree.Tree{t31, 3, t38}

	t53 := &tree.Tree{t31, 3, t15}
	t78 := &tree.Tree{t53, 8, t113}

	printTree(t73)
	printTree(t78)

	printTree(tree.New(1))

	fmt.Println(Same(t73, t73))
	fmt.Println(Same(t73, t78))
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))

	fmt.Println("end")

}
