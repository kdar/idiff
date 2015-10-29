idiff
=====

Interface diff is a package that diffs two interface{}s. It can be a basic Go type, struct, map, pointer, array, or slice.

Right now it just provides a formatter to print something you can use in testing.

## Example

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

Output:

    They're not equal:
    map[string]string["last"]: got: "Betterton", expected: "Snow"

## Formatters

For now there is only one test formatter. All the fields of DiffResult are public so feel free to make your own formatter.
