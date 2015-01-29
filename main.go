package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

func check_err(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error:%s", err)
	}
	ga_token_file := path.Join(usr.HomeDir, ".google_authenticator")

	file, err := os.Open(ga_token_file)
	check_err(err)
	defer file.Close()

}
