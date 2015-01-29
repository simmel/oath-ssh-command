package main

import (
	"fmt"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error:%s", err)
	}
	fmt.Println(usr.HomeDir)
}
