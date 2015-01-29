package main

import (
	"bufio"
	"encoding/base32"
	"fmt"
	"github.com/hgfischer/go-otp"
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

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	check_err(scanner.Err())
	if len(scanner.Text()) != 16 {
		fmt.Printf("ERROR: Couldn't read exactly 16 bytes from first line of %q. Got this: %q.", ga_token_file, scanner.Text())
	}
	token_in_bytes, err := base32.StdEncoding.DecodeString(scanner.Text())
	check_err(err)
	token := string(token_in_bytes)

	fmt.Print("Verification code: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	check_err(scanner.Err())
	otp_input := scanner.Text()

	totp := &otp.TOTP{Secret: token}
	verified := totp.Verify(otp_input)

	os.Exit(func(b bool) int {
		if b {
			return 0
		}
		return 1
	}(verified))
}
