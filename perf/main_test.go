package main

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	n := NewNeuron(10, 10, 10)

	for i := 0; i < 10; i++ {
		n.enqueue(1)
		n.process()

		fmt.Println(fmt.Sprintf("Fired %v, pot %d", n.fired, n.potential))
	}
}


func BenchmarkNewNeuron(b *testing.B) {

	n1 := NewNeuron(10, 10, 1)
	n2 := NewNeuron(10, 10, 1)
	// ratio 20


	for i := 0; i < 60; i++ {
		c := NewSynapse(n2, 3, 1)
		n1.targets = append(n1.targets, c)
	}

	for n := 0; n < b.N; n++ {
		n1.enqueue(1)

		n1.process()
		n2.process()


	}
}
