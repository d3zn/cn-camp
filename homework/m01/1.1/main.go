package main

import "fmt"

func main() {
	s := []string{"I", "am", "stupid", "and", "weak"}
	for i, v := range s {
		switch v {
		case "stupid":
			s[i] = "smart"
		case "weak":
			s[i] = "strong"
		}
	}
	fmt.Println(s)
}
