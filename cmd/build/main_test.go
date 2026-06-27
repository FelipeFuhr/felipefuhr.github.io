package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_BuildsSiteFromRealContent(t *testing.T) {
	out := filepath.Join(t.TempDir(), "dist")
	err := run([]string{
		"-data=../../data",
		"-templates=../../templates",
		"-assets=../../assets",
		"-static=../../static",
		"-out=" + out,
	}, io.Discard)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "index.html")); err != nil {
		t.Errorf("expected index.html: %v", err)
	}
}

func TestRun_InvalidFlag_ReturnsError(t *testing.T) {
	if err := run([]string{"-nope"}, io.Discard); err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestRun_MissingData_ReturnsError(t *testing.T) {
	if err := run([]string{"-data=/does/not/exist"}, io.Discard); err == nil {
		t.Error("expected error for missing data dir")
	}
}
