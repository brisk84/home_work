package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	reversedString := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(reversedString)
}
