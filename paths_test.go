package xdg_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/miquella/xdg"
)

func TestPathGlob(t *testing.T) {
	dir := createDir(t, "abc", "cbc", "xyz")
	defer removeDir(dir)
	expected := []string{filepath.Join(dir, "abc"), filepath.Join(dir, "cbc")}
	sort.Strings(expected)

	p := xdg.Path(dir)
	actual, err := p.Glob("*bc")
	if err != nil {
		t.Fatalf("globbing failed: %v", err)
	}

	sort.Strings(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %#v, got %#v", expected, actual)
	}
}

func TestPathJoin(t *testing.T) {
	path := filepath.FromSlash("/home/test/.local/share")
	expected := filepath.Join(path, "test")

	p := xdg.Path(path)
	actual := p.Join("test")
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestPathsGlob(t *testing.T) {
	dir1 := createDir(t, "abc", "cbc", "xyz")
	defer removeDir(dir1)
	dir2 := createDir(t, "bbc", "tcl", "pcb")
	defer removeDir(dir2)
	expected := []string{filepath.Join(dir1, "abc"), filepath.Join(dir1, "cbc")}
	expected = append(expected, filepath.Join(dir2, "bbc"))
	sort.Strings(expected)

	p := xdg.Paths{xdg.Path(dir1), xdg.Path(dir2)}
	actual, err := p.Glob("*bc")
	if err != nil {
		t.Fatalf("globbing failed: %v", err)
	}

	sort.Strings(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %#v, got %#v", expected, actual)
	}
}

func createDir(t *testing.T, filenames ...string) string {
	dir, err := ioutil.TempDir("", "xdg_path")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	for _, filename := range filenames {
		err = ioutil.WriteFile(filepath.Join(dir, filename), []byte(filename), 0644)
		if err != nil {
			t.Fatalf("failed to write test file %s: %v", filename, err)
		}
	}

	return dir
}

func removeDir(dir string) {
	os.RemoveAll(dir)
}
