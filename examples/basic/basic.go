package main

import (
	"fmt"

	"github.com/kdar/idiff"
)

func main() {
	m1 := map[string]string{
		"name": "John",
		"last": "Snow",
	}
	m2 := map[string]string{
		"name": "John",
		"last": "Betterton",
	}

	diff, equal := idiff.Diff(m1, m2)
	if equal {
		fmt.Println("They're equal!")
	} else {
		fmt.Println("They're not equal:")
		fmt.Println(idiff.FormatTest(diff))
	}
}
