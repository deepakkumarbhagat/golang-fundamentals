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

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	go func() {
		for x := range ch1 {
			ch3 <- x
		}
		close(ch3)
	}()

	for {
		select {
		case <-ch3:
			return true
		case y := <-ch2:
			x := <-ch3
			if x != y {
				return false
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}

	return true
}

func main() {
	t1 := tree.New(2)
	t2 := tree.New(2)

	fmt.Print(Same(t1, t2))

	time.Sleep(100)
}
