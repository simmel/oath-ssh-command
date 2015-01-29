package main

import (
	"fmt"
	"os/user"
	"path"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error:%s", err)
	}
	fmt.Println(path.Join(usr.HomeDir, ".google_authenticator"))
}
