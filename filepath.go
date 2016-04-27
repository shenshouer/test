package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println(filepath.Dir("/playlist/zip/zip_20151114232152093.zip/"))
	fmt.Println(filepath.Base("/playlist/zip/zip_20151114232152093.zip"))
	fmt.Println(strings.Split(filepath.Base("/playlist/zip/zip_20151114232152093.zip"), ".")[0])

	fmt.Println(strings.Trim("wqwqwqw.tx1t32", ".txt"))

}
