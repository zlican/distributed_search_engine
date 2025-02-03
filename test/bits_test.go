package test

import (
	"engine/utils"
	"fmt"
	"testing"
)

func TestBits(t *testing.T) {
	var n uint64
	n = utils.SetBit1(n, 11)
	n = utils.SetBit1(n, 28)

	fmt.Println(utils.IsBit1(n, 11)) // true
	fmt.Println(utils.IsBit1(n, 28)) // true
	fmt.Println(utils.IsBit1(n, 29)) // false

	fmt.Println(utils.CountBit1(n))
	fmt.Printf("%064b\n", n)
}
