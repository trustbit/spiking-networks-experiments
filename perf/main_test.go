package main

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	n := NewNeuron(10, 10)

	for i := 0; i < 10; i++ {
		n.enqueue(1)
		n.process()

		fmt.Println(fmt.Sprintf("Fired %v, pot %d", n.fired, n.potential))
	}
}
