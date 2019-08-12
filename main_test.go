package main_test

import (
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func gpgVerify(args ...string) error {
	cmd := exec.Command("./pgp-verify", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TestAll(t *testing.T) {
	// Setup
	gnupgHome, err := ioutil.TempDir("", "docker-remote")
	ok(t, err)

	t.Run("valid-asc-key", func(t *testing.T) {
		ok(t, gpgVerify("test/channel-rust-beta-date.txt",
			"test/channel-rust-beta-date.txt.asc",
			"test/rust-key.gpg.ascii"))

	})

	t.Run("invalid-asc-key", func(t *testing.T) {
		ok(t, gpgVerify("test/channel-rust-beta-date.txt",
			"test/channel-rust-beta-date.txt.asc",
			"test/rust-key.gpg.ascii"))

	})

	t.Run("accept-expired-key", func(t *testing.T) {
		ok(t, gpgVerify("test/channel-rust-beta-date.txt",
			"test/channel-rust-beta-date.txt.asc",
			"test/rust-key.gpg.ascii"))

	})

	t.Run("reject-expired-key", func(t *testing.T) {
		ok(t, gpgVerify("test/channel-rust-beta-date.txt",
			"test/channel-rust-beta-date.txt.asc",
			"test/rust-key.gpg.ascii"))

	})

	// Teardown
	ok(t, os.RemoveAll(gnupgHome))

}
