package main

import (
	"bufio"
	"fmt"
	"github.com/hgfischer/go-otp"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"syscall"
)

type file interface {
	io.Closer
	io.Writer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
	Truncate(size int64) error
}

var check_err = func(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		fail_out()
	}
}

var find_config = func() (filename string) {
	usr, err := user.Current()
	check_err(err)
	return path.Join(usr.HomeDir, ".google_authenticator")
}

var get_config_file = func(filename string) (config_file file) {
	config_file, err := os.Open(filename)
	check_err(err)
	return config_file
}

var parse_config = func(filename string) (token string) {
	file := get_config_file(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	check_err(scanner.Err())
	if len(scanner.Text()) != 16 {
		check_err(fmt.Errorf("Couldn't read exactly 16 bytes from first line of %q. Got this: %q.", filename, scanner.Text()))
	}
	token_in_bytes := scanner.Text()
	return string(token_in_bytes)
}

var read_otp_input = func() (otp string) {
	if os.Getenv("OTP_TOKEN") != "" {
		return os.Getenv("OTP_TOKEN")
	}

	fmt.Print("Verification code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	check_err(scanner.Err())
	return scanner.Text()
}

var run_appropriately = func() {
	env := os.Environ()

	// FIXME Maybe check if $SHELL is set to something?
	shell, err := exec.LookPath(os.Getenv("SHELL"))
	check_err(err)

	if os.Getenv("SSH_ORIGINAL_COMMAND") != "" {
		args := []string{shell, "-c", os.Getenv("SSH_ORIGINAL_COMMAND")}
		exec_appropriately(shell, args, env)
	} else {
		args := []string{shell, "-i"}
		exec_appropriately(shell, args, env)
	}
}

var exec_appropriately = func(shell string, args []string, env []string) {
	syscall.Exec(shell, args, env)
}

var fail_out = func() {
	os.Exit(1)
}

func main() {
	ga_token_file := find_config()

	token := parse_config(ga_token_file)

	otp_input := read_otp_input()

	totp := &otp.TOTP{Secret: token, IsBase32Secret: true}
	verified := totp.Verify(otp_input)

	if verified {
		run_appropriately()
	} else {
		fail_out()
	}
}
