//go:build !linux
// +build !linux

package testtools

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func LoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func TempFile(t *testing.T, relativePath, goldenName string) func() {
	ex, _ := os.Executable()
	tmpDir := filepath.Join(filepath.Dir(ex), filepath.Dir(relativePath))
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	filename := filepath.Base(relativePath)
	tmpFile, err := os.Create(filepath.Join(tmpDir, filename))
	if err != nil {
		t.Fatal(err)
	}

	golden := filepath.Join("testdata", goldenName)
	data, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Write(data) //nolint

	return func() {
		os.Remove(filename)
		os.RemoveAll(tmpDir)
	}
}
