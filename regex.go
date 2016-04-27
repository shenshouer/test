package main

import (
	"fmt"
)

import (
	"regexp"
)

func main() {
	r := regexp.MustCompile(`^/Users/sope/Desktop/index1/\d{2}.*$`)
	r2 := regexp.MustCompile(fmt.Sprintf(`^/Users/sope/.*/index$`))

	teststr1 := "/Users/sope/Desktop/index1/img/backhome.gif"
	teststr2 := "/Users/sope/Desktop/index1/24/24.html"
	teststr3 := "/Users/sope/Desktop/index1"
	teststr4 := "/Users/sope/Desktop/index"

	fmt.Println(r.MatchString(teststr1))
	fmt.Println(r.MatchString(teststr2))
	fmt.Println(r2.MatchString(teststr3))
	fmt.Println(r2.MatchString(teststr4))
}
