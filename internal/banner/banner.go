package banner

import (
	_ "embed"
	"fmt"
)

//go:embed ascii.txt
var ascii string

func Print() {
	fmt.Print("\033[H\033[2J")
	fmt.Println(ascii)
	fmt.Println("Project: Spotis")
	fmt.Println("Author: Mayusha256")
	fmt.Println()
}
