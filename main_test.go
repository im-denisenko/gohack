package main

import (
	"os"
	"testing"
)

func TestResolvePath(t *testing.T) {
	cwd, err := os.Getwd()
	Check(err)

	relativePath := "input/a.json"
	absolutePath := cwd + "/input/a.json"

	relativePathResolved, _ := ResolvePath(relativePath)
	absolutePathResolved, _ := ResolvePath(absolutePath)

	if relativePathResolved != cwd+"/"+relativePath {
		t.Fatalf("relative path resolve was incorrect")
	}
	if absolutePathResolved != absolutePath {
		t.Fatalf("absolute path resolve was incorrect")
	}
}
