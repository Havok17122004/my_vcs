package pkg

import "fmt"

// panic the error
func Check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
