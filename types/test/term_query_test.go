package test

import (
	"engine/types"
	"fmt"
	"testing"
)

func TestTermQuery(t *testing.T) {
	//(A|B|C)&(E|(A&C))
	A := types.Be("A")
	B := types.Be("B")
	C := types.Be("C")
	E := types.Be("E")

	resultStr := A.Or(B, C).And(E.Or(A.And(C)))
	fmt.Println(resultStr.ToString())

}
