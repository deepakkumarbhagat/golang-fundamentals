package main

import (
	"fmt"
	"time"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	if t.Left != nil {
		Walk(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

func Walking(t *tree.Tree, ch chan int) {
	Walk(t, ch)
	defer close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walking(t1, ch1)
	go Walking(t2, ch2)

	for {
		v1, more1 := <-ch1
		v2, more2 := <-ch2

		if more1 != more2 || v1 != v2 {
			return false
		}

		//return if channel closed
		if !more1 {
			return true
		}
	}
}

func main() {
	t1 := tree.New(2)
	t2 := tree.New(1)

	fmt.Print(Same(t1, t2))

	time.Sleep(100)
}
