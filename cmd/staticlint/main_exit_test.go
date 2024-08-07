package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGetExitAnalyzer(t *testing.T) {
	analyzer := GetExitAnalyzer()
	if analyzer == nil {
		t.Fatalf("GetExitAnalyzer() returned nil")
	}
	expectedName := "os_exit_checker"
	if analyzer.Name != expectedName {
		t.Errorf("GetExitAnalyzer().Name returned %s where %s was expected", analyzer.Name, expectedName)
	}
}

func TestExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), GetExitAnalyzer(), "./...")
}
