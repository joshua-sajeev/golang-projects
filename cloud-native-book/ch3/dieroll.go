package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type dieRollFunc func(int) int

func fakeDieRoll(size int) int {
	return 42
}

func getDieRolls() []dieRollFunc {
	return []dieRollFunc{
		dieRoll, fakeDieRoll,
	}
}
func dieRoll(size int) int {
	return rand.Intn(size) + 1
}

func rollTwo(size1, size2 int) (int, int) {
	return dieRoll(size1), dieRoll(size2)
}

func returnsNames(input1 string, input2 int) (theResult string, err error) {
	theResult = "modified " + input1 + "," + strconv.Itoa(input2)
	return theResult, err
}
func main() {
	fmt.Printf("Rolled a die of size %d, result: %d\n", 6, dieRoll(6))

	res1, res2 := rollTwo(6, 10)
	fmt.Printf("Rolled a pair of dice (%d,%d), results: %d, %d\n", 6, 10, res1, res2)

	named, err := returnsNames("globule", 42)
	fmt.Printf("Named params returned: '%s',%v\n", named, err)
}
