package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func buildNixFile(t testing.TB, tempdir string, nixFile string, expectedBuilds int) int {
	store := path.Join(tempdir, "store")
	output := path.Join(tempdir, path.Base(nixFile))
	err := os.MkdirAll(output, 0700)
	ok(t, err)

	buildFlags := fmt.Sprintf("-I nixpkgs=channel:nixos-unstable-small --store '%s' -o '%s'", store, path.Join(output, "result"))
	flags := []string{"-build-flags", buildFlags, "-f", nixFile}
	fmt.Printf("nix-build-uncached %s\n", strings.Join(flags, " "))
	err = realMain(flags)
	ok(t, err)

	files, err := ioutil.ReadDir(output)
	ok(t, err)
	return len(files)
}

func TestAll(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	ok(t, err)
	_, filename, _, _ := runtime.Caller(0)
	asset := path.Join(path.Dir(filename), "test")

	builds := buildNixFile(t, tempdir, path.Join(asset, "test-skip-cached.nix"), 1)
	equals(t, builds, 1)

	builds = buildNixFile(t, tempdir, path.Join(asset, "test-many-drvs.nix"), 1)
	// depends a bit ulimit, but we should be able to build at least 100
	if 100 > builds {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\t: %#v > %#v\033[39m\n\n", filepath.Base(file), line, 100, builds)
		t.FailNow()
	}
}
