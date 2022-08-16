package testdata

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("test os.Exit(0)")
	os.Exit(0) // want "found os.Exit"
}
