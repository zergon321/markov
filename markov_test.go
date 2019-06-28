package markov_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tracer8086/markov"
)

func JSONTest(t *testing.T) {
	chain := markov.CreateNew()

	for i := 0; i < 8; i++ {
		chain.AddEvent("s" + strconv.Itoa(i))
	}

	data, err := chain.ToJSON()
	assert.NoError(t, err)

	_, err = markov.FromJSON(data)
	assert.NoError(t, err)
}
