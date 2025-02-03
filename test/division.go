package test

import (
	"engine/utils"
	"fmt"
	"testing"
)

func TestDivision(t *testing.T) {
	docs := []*utils.Doc{
		&utils.Doc{1, []string{"a", "b", "c"}},
		&utils.Doc{2, []string{"b", "c"}},
		&utils.Doc{3, []string{"c"}}}

	index := utils.Division(docs)

	for k, value := range index {
		fmt.Println(k, value)
	}
}

// Run the test:
// go test -v ./test -run=^TestDivision$ -count=1
