package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func reverseHellow() {
	fmt.Println(stringutil.Reverse("Hello, OTUS!"))
}

func main() {
	reverseHellow()
}
