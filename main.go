package main

import (
	"bufio"
	"encoding/base32"
	"fmt"
	"github.com/hgfischer/go-otp"
	"os"
	"os/exec"
	"os/user"
	"path"
	"syscall"
)

func check_err(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func parse_config(filename string) (token string) {
	file, err := os.Open(filename)
	check_err(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	check_err(scanner.Err())
	if len(scanner.Text()) != 16 {
		fmt.Printf("ERROR: Couldn't read exactly 16 bytes from first line of %q. Got this: %q.", filename, scanner.Text())
	}
	token_in_bytes, err := base32.StdEncoding.DecodeString(scanner.Text())
	check_err(err)
	return string(token_in_bytes)
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error:%s", err)
	}
	ga_token_file := path.Join(usr.HomeDir, ".google_authenticator")

	token := parse_config(ga_token_file)

	fmt.Print("Verification code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	check_err(scanner.Err())
	otp_input := scanner.Text()

	totp := &otp.TOTP{Secret: token}
	verified := totp.Verify(otp_input)

	if verified {
		env := os.Environ()

		if os.Getenv("SSH_ORIGINAL_COMMAND") != "" {
			// FIXME Maybe check if $SHELL is set to something?
			shell, err := exec.LookPath(os.Getenv("SHELL"))
			check_err(err)

			// FIXME Maybe check if $USER is set to something?
			args := []string{shell, "-c", os.Getenv("SSH_ORIGINAL_COMMAND")}
			syscall.Exec(shell, args, env)
		} else {
			shell, err := exec.LookPath("login")
			check_err(err)

			// FIXME Maybe check if $USER is set to something?
			args := []string{"login", "-f", os.Getenv("USER")}
			syscall.Exec(shell, args, env)
		}
	} else {
		os.Exit(1)
	}
}
