package test

import (
	"fmt"
	"testing"

	"github.com/zlican/engine/utils"
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

func TestInterset(t *testing.T) {
	min := 10
	var arr1 = []int{10, 15, 16, 19, 22}
	var arr2 = []int{11, 15, 16, 17, 20, 22}
	arr1BitMap := utils.CreateBitMap(min, arr1)
	arr2BitMap := utils.CreateBitMap(min, arr2)
	res := utils.IntersectionOfBitMap(arr1BitMap, arr2BitMap, min)
	fmt.Println(res)
}
