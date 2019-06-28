package markov

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Chain represents a Markov's chain with events and their probabilities.
type Chain struct {
	edges  map[string]map[string]int64
	totals map[string]int64
}

// CreateNew returns a new empty Markov's chain.
func CreateNew() *Chain {
	edges := make(map[string]map[string]int64, 0)
	totals := make(map[string]int64)

	return &Chain{edges, totals}
}

// AddEvent adds a new event in the chain.
func (chain *Chain) AddEvent(event string) error {
	if _, ok := chain.edges[event]; ok {
		return fmt.Errorf("Event %s already exists in the chain", event)
	}

	chain.edges[event] = make(map[string]int64, 0)
	chain.totals[event] = 0

	return nil
}

// GetTransitionWeights returns a dictionary of transition weights for the event.
func (chain *Chain) GetTransitionWeights(event string) (map[string]int64, error) {
	if _, ok := chain.edges[event]; !ok {
		return nil, fmt.Errorf("Event %s doesn't exist in the chain", event)
	}

	result := make(map[string]int64, 0)

	for key, value := range chain.edges[event] {
		result[key] = value
	}

	return result, nil
}

// GetTransitionProbabilities returns a dictionary of transition probabilities for the event.
func (chain *Chain) GetTransitionProbabilities(event string) (map[string]float64, error) {
	if _, ok := chain.edges[event]; !ok {
		return nil, fmt.Errorf("Event %s doesn't exist in the chain", event)
	}

	result := make(map[string]float64, 0)

	for key, value := range chain.edges[event] {
		result[key] = float64(value) / float64(chain.totals[event])
	}

	return result, nil
}

// RemoveEvent removes all the occurences of the given event from the Markov's chain.
func (chain *Chain) RemoveEvent(event string) error {
	if _, ok := chain.edges[event]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", event)
	}

	// Delete an event from other events' transitions.
	for key := range chain.edges {
		chain.RemoveTransition(key, event)
	}

	// Delete an event.
	delete(chain.edges, event)
	delete(chain.totals, event)

	return nil
}

// AddTransition adds a new transition from one event to another in the chain.
func (chain *Chain) AddTransition(outgoing string, incoming string, weight int64) error {
	if _, ok := chain.edges[outgoing]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.edges[incoming]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.edges[outgoing][incoming]; ok {
		return fmt.Errorf("The transition from %s to %s already exists in the chain", outgoing, incoming)
	}

	if weight <= 0 {
		return errors.New("Weight should be above zero")
	}

	chain.edges[outgoing][incoming] = weight
	chain.totals[outgoing] += weight

	return nil
}

// GetTransitionWeight returns a weight of the specified transition.
func (chain *Chain) GetTransitionWeight(outgoing string, incoming string) (int64, error) {
	if _, ok := chain.edges[outgoing]; !ok {
		return 0, fmt.Errorf("Event %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.edges[incoming]; !ok {
		return 0, fmt.Errorf("Event %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.edges[outgoing][incoming]; !ok {
		return 0, fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	return chain.edges[outgoing][incoming], nil
}

// GetTransitionProbability returns a probability of the specified transition
func (chain *Chain) GetTransitionProbability(outgoing string, incoming string) (float64, error) {
	if _, ok := chain.edges[outgoing]; !ok {
		return 0, fmt.Errorf("Event %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.edges[incoming]; !ok {
		return 0, fmt.Errorf("Event %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.edges[outgoing][incoming]; !ok {
		return 0, fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	return float64(chain.edges[outgoing][incoming]) / float64(chain.totals[outgoing]), nil
}

// UpdateTransition changes the weight value of the existing transition.
func (chain *Chain) UpdateTransition(outgoing string, incoming string, weight int64) error {
	if _, ok := chain.edges[outgoing]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.edges[incoming]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.edges[outgoing][incoming]; !ok {
		return fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	oldWeight := chain.edges[outgoing][incoming]
	chain.totals[outgoing] -= oldWeight

	chain.edges[outgoing][incoming] = weight
	chain.totals[outgoing] += weight

	return nil
}

// RemoveTransition removes the transition from the outgoing state to the incoming state.
func (chain *Chain) RemoveTransition(outgoing string, incoming string) error {
	if _, ok := chain.edges[outgoing]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.edges[incoming]; !ok {
		return fmt.Errorf("Event %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.edges[outgoing][incoming]; !ok {
		return fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	chain.totals[outgoing] -= chain.edges[outgoing][incoming]
	delete(chain.edges[outgoing], incoming)

	return nil
}

// ToJSON serializes the Markov's chain to JSON format.
func (chain *Chain) ToJSON() ([]byte, error) {
	input := []interface{}{chain.edges, chain.totals}
	data, err := json.MarshalIndent(input, "", "    ")

	if err != nil {
		return nil, err
	}

	return data, nil
}

// FromJSON creates a new Markov's chain from its JSON representation.
func FromJSON(data []byte) (*Chain, error) {
	var input []interface{}
	err := json.Unmarshal(data, &input)

	if err != nil {
		return nil, err
	}

	return &Chain{input[0].(map[string]map[string]int64), input[1].(map[string]int64)}, nil
}
