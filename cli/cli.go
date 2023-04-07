package cli

import (
	"fmt"
	"os"
)

// func GetScreenplayFile() {}
func GetScreenplayContents(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Read file")
	}
	return string(dat)
}
