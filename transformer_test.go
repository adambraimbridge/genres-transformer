package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name  string
		term  term
		genre genre
	}{
		{"Transform term to genre", term{
			CanonicalName: "Market Report",
			RawID:         "R0xfMTY0ODM1-R2VucmVz"},
			genre{
				UUID:      "c6844334-c743-31c7-a4ea-61587e006e9b",
				PrefLabel: "Market Report",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"UjB4Zk1UWTBPRE0xLVIyVnVjbVZ6-R2VucmVz"},
					Uuids: []string{"c6844334-c743-31c7-a4ea-61587e006e9b"},
				},
				PrimaryType:   primaryType,
				TypeHierarchy: genreTypes,
			}},
	}

	for _, test := range tests {
		expectedGenre := transformGenre(test.term, "Genres")
		assert.Equal(test.genre, expectedGenre, fmt.Sprintf("%s: Expected genre incorrect", test.name))
	}

}
