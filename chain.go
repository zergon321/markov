package markov

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Chain represents a Markov's chain with states and their probabilities.
type Chain struct {
	transitions map[string]map[string]int64
	totals      map[string]int64
}

type chainExport struct {
	Transitions map[string]map[string]int64
	Totals      map[string]int64
}

// CreateNew returns a new empty Markov's chain.
func CreateNew() *Chain {
	transitions := make(map[string]map[string]int64, 0)
	totals := make(map[string]int64)

	return &Chain{transitions, totals}
}

// AddState adds a new state in the chain.
func (chain *Chain) AddState(state string) error {
	if _, ok := chain.transitions[state]; ok {
		return fmt.Errorf("State %s already exists in the chain", state)
	}

	chain.transitions[state] = make(map[string]int64, 0)
	chain.totals[state] = 0

	return nil
}

// GetAllStates returns a slice of the all states stored in the chain.
func (chain *Chain) GetAllStates() []string {
	states := make([]string, 0)

	for state := range chain.totals {
		states = append(states, state)
	}

	return states
}

// Transit returns a new state for the system by the current state.
// Returns an empty string if the current state is terminal.
func (chain *Chain) Transit(state string) (string, error) {
	var (
		max  float64
		next string
	)

	if _, ok := chain.transitions[state]; !ok {
		return "", fmt.Errorf("State %s doesn't exist in the chain", state)
	}

	probabilities, err := chain.GetTransitionProbabilities(state)

	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Multiply all the probabilities by random number and set the state with max probability as the next state.
	for key := range probabilities {
		temp := probabilities[key] * rand.Float64()

		if temp > max {
			max = temp
			next = key
		}
	}

	return next, nil
}

// HasState returns 'true' if the specified state exists in the chain and 'false' otherwise.
func (chain *Chain) HasState(state string) bool {
	_, ok := chain.transitions[state]

	return ok
}

// GetTransitionWeights returns a dictionary of transition weights for the state.
func (chain *Chain) GetTransitionWeights(state string) (map[string]int64, error) {
	if _, ok := chain.transitions[state]; !ok {
		return nil, fmt.Errorf("State %s doesn't exist in the chain", state)
	}

	result := make(map[string]int64, 0)

	for key, value := range chain.transitions[state] {
		result[key] = value
	}

	return result, nil
}

// GetTransitionProbabilities returns a dictionary of transition probabilities for the state.
func (chain *Chain) GetTransitionProbabilities(state string) (map[string]float64, error) {
	if _, ok := chain.transitions[state]; !ok {
		return nil, fmt.Errorf("State %s doesn't exist in the chain", state)
	}

	result := make(map[string]float64, 0)

	for key, value := range chain.transitions[state] {
		result[key] = float64(value) / float64(chain.totals[state])
	}

	return result, nil
}

// RemoveState removes all the occurences of the given State from the Markov's chain.
func (chain *Chain) RemoveState(state string) error {
	if _, ok := chain.transitions[state]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", state)
	}

	// Delete a state from other states' transitions.
	for key := range chain.transitions {
		chain.RemoveTransition(key, state)
	}

	// Delete an State.
	delete(chain.transitions, state)
	delete(chain.totals, state)

	return nil
}

// AddTransition adds a new transition from one state to another in the chain.
func (chain *Chain) AddTransition(outgoing string, incoming string, weight int64) error {
	if _, ok := chain.transitions[outgoing]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.transitions[incoming]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.transitions[outgoing][incoming]; ok {
		return fmt.Errorf("The transition from %s to %s already exists in the chain", outgoing, incoming)
	}

	if weight <= 0 {
		return errors.New("Weight should be above zero")
	}

	chain.transitions[outgoing][incoming] = weight
	chain.totals[outgoing] += weight

	return nil
}

// HasTransition returns 'true' if transition exists in the chain and 'false' otherwise.
func (chain *Chain) HasTransition(outgoing string, incoming string) bool {
	_, outgoingExists := chain.transitions[outgoing]
	_, incomingExists := chain.transitions[incoming]

	if outgoingExists && incomingExists {
		_, ok := chain.transitions[outgoing][incoming]

		return ok
	}

	return false
}

// GetTransitionWeight returns a weight of the specified transition.
func (chain *Chain) GetTransitionWeight(outgoing string, incoming string) (int64, error) {
	if _, ok := chain.transitions[outgoing]; !ok {
		return 0, fmt.Errorf("State %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.transitions[incoming]; !ok {
		return 0, fmt.Errorf("State %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.transitions[outgoing][incoming]; !ok {
		return 0, fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	return chain.transitions[outgoing][incoming], nil
}

// GetTransitionProbability returns a probability of the specified transition
func (chain *Chain) GetTransitionProbability(outgoing string, incoming string) (float64, error) {
	if _, ok := chain.transitions[outgoing]; !ok {
		return 0, fmt.Errorf("State %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.transitions[incoming]; !ok {
		return 0, fmt.Errorf("State %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.transitions[outgoing][incoming]; !ok {
		return 0, fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	return float64(chain.transitions[outgoing][incoming]) / float64(chain.totals[outgoing]), nil
}

// UpdateTransition changes the weight value of the existing transition.
func (chain *Chain) UpdateTransition(outgoing string, incoming string, weight int64) error {
	if _, ok := chain.transitions[outgoing]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.transitions[incoming]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.transitions[outgoing][incoming]; !ok {
		return fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	oldWeight := chain.transitions[outgoing][incoming]
	chain.totals[outgoing] -= oldWeight

	chain.transitions[outgoing][incoming] = weight
	chain.totals[outgoing] += weight

	return nil
}

// RemoveTransition removes the transition from the outgoing state to the incoming state.
func (chain *Chain) RemoveTransition(outgoing string, incoming string) error {
	if _, ok := chain.transitions[outgoing]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", outgoing)
	}

	if _, ok := chain.transitions[incoming]; !ok {
		return fmt.Errorf("State %s doesn't exist in the chain", incoming)
	}

	if _, ok := chain.transitions[outgoing][incoming]; !ok {
		return fmt.Errorf("The transition from %s to %s doesn't exist in the chain", outgoing, incoming)
	}

	chain.totals[outgoing] -= chain.transitions[outgoing][incoming]
	delete(chain.transitions[outgoing], incoming)

	return nil
}

// ToJSON serializes the Markov's chain to JSON format.
func (chain *Chain) ToJSON() ([]byte, error) {
	input := chainExport{chain.transitions, chain.totals}
	data, err := json.MarshalIndent(input, "", "    ")

	if err != nil {
		return nil, err
	}

	return data, nil
}

// FromJSON creates a new Markov's chain from its JSON representation.
func FromJSON(data []byte) (*Chain, error) {
	var input chainExport
	err := json.Unmarshal(data, &input)

	if err != nil {
		return nil, err
	}

	return &Chain{input.Transitions, input.Totals}, nil
}
