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
	defer Patch(&get_config_file, func(filename string) (file *mockfile.file) {
		mockfile := mockfile.New(os.Getenv("HOME") + "/.google_authenticator")
		mockfile.Write([]byte("luul"))
		defer mockfile.Close()
		return mockfile
	}).Restore()

	expected := "1337"
	ga_token_file := find_config()
	found := parse_config(ga_token_file)
	if found != expected {
		t.Errorf("%q != %q", found, expected)
	}
}
