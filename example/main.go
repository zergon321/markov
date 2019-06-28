package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/tracer8086/markov"
)

func main() {
	chain := markov.CreateNew()

	// Add states.
	for i := 0; i < 8; i++ {
		chain.AddState("s" + strconv.Itoa(i))
	}

	// Add transitions.
	chain.AddTransition("s0", "s1", 10)
	chain.AddTransition("s1", "s2", 2)
	chain.AddTransition("s1", "s3", 8)
	chain.AddTransition("s2", "s4", 10)
	chain.AddTransition("s3", "s4", 1)
	chain.AddTransition("s3", "s5", 4)
	chain.AddTransition("s3", "s6", 5)
	chain.AddTransition("s5", "s7", 3)
	chain.AddTransition("s6", "s7", 7)
	chain.AddTransition("s5", "s6", 4)
	chain.AddTransition("s5", "s5", 7)

	// Output the chain.
	data, err := chain.ToJSON()

	if err != nil {
		log.Fatalln("Couldn't serialize chain to JSON:", err)
	}

	fmt.Println(string(data))

	restoredChain, err := markov.FromJSON(data)

	if err != nil {
		log.Fatalln("Couldn't restore the chain from JSON:", err)
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			probability, err := restoredChain.GetTransitionProbability("s"+strconv.Itoa(i), "s"+strconv.Itoa(j))

			if err == nil {
				fmt.Printf("s%d - s%d: %0.2f\n", i, j, probability)
			}
		}
	}

	fmt.Println(restoredChain.HasState("s0"))
	fmt.Println(restoredChain.HasState("s8"))
	fmt.Println(restoredChain.HasTransition("s3", "s6"))
	fmt.Println(restoredChain.HasTransition("s4", "s5"))
}
