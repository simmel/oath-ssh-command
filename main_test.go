package main

import (
	"github.com/seletskiy/go-mock-file"
	"os"
	"reflect"
	"testing"
)

/*
https://groups.google.com/d/msg/golang-nuts/6AN1E2CJOxI/kjP-VhXVleoJ
Usage:
	var find_config = func() (filename string) {
		return "something"
	}

	func TestConfig(t *testing.T) {
		defer Patch(&find_config, func() (filename string) {
			return "something else"
		}).Restore()

		expected := "something else"
		config := find_config()
		if config != expected {
			t.Errorf("%q == %q", expected, config)
		}
	}
*/

// Restorer holds a function that can be used
// to restore some previous state.
type Restorer func()

// Restore restores some previous state.
func (r Restorer) Restore() {
	r()
}

// Patch sets the value pointed to by the given destination to the given
// value, and returns a function to restore it to its original value.  The
// value must be assignable to the element type of the destination.
func Patch(dest, value interface{}) Restorer {
	destv := reflect.ValueOf(dest).Elem()
	oldv := reflect.New(destv.Type()).Elem()
	oldv.Set(destv)
	valuev := reflect.ValueOf(value)
	if !valuev.IsValid() {
		// This isn't quite right when the destination type is not
		// nilable, but it's better than the complex alternative.
		valuev = reflect.Zero(destv.Type())
	}
	destv.Set(valuev)
	return func() {
		destv.Set(oldv)
	}
}

// My tests

func TestConfig(t *testing.T) {
	expected := os.Getenv("HOME") + "/.google_authenticator"
	config := find_config()
	if config != expected {
		t.Errorf("%q == %q", expected, config)
	}
}

func TestConfigParsing(t *testing.T) {
	defer Patch(&get_config_file, func(filename string) (config_file file) {
		config_file = mockfile.New(os.Getenv("HOME") + "/.google_authenticator")
		config_file.Write([]byte("KJGFSRCFIFCE4MCC"))
		defer config_file.Close()
		return config_file
	}).Restore()

	expected := "KJGFSRCFIFCE4MCC"
	ga_token_file := find_config()
	found := parse_config(ga_token_file)
	if found != expected {
		t.Errorf("%q != %q", found, expected)
	}
}

func TestExec(t *testing.T) {
	expected := []string{"login", "-f", os.Getenv("USER")}
	defer Patch(&exec_appropriately, func(shell string, args []string, env []string) {
		if len(args) != len(expected) {
			t.Errorf("args length %d differs from expected length %d", len(args), len(expected))
		}
		for i := range args {
			if args[i] != expected[i] {
				t.Errorf("args[%d] %q != expected[%d] %q", i, args[i], i, expected[i])
			}
		}
	}).Restore()

	run_appropriately()
}

func TestExecEnv(t *testing.T) {
	expected := "/bin/true"
	defer Patch(&exec_appropriately, func(shell string, args []string, env []string) {
		found := args[len(args)-1]
		if found != expected {
			t.Errorf("%q != %q", found, expected)
		}
	}).Restore()

	check_err(os.Setenv("SSH_ORIGINAL_COMMAND", expected))
	run_appropriately()
	check_err(os.Unsetenv("SSH_ORIGINAL_COMMAND"))
}

func TestOtpTokenEnv(t *testing.T) {
	expected := "313373"
	check_err(os.Setenv("OTP_TOKEN", expected))
	found := read_otp_input()
	if found != expected {
		t.Errorf("%q != %q", found, expected)
	}
	check_err(os.Unsetenv("OTP_TOKEN"))
}
