package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// run_tests.go is a test runner for the Libros project
// It runs all unit and integration tests and provides a summary of results
func main() {
	fmt.Println("=== Libros Test Suite ===")
	fmt.Println()

	// Get the project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("Error finding project root: %v\n", err)
		os.Exit(1)
	}

	// Change to project root for running tests
	err = os.Chdir(projectRoot)
	if err != nil {
		fmt.Printf("Error changing to project root: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running tests from: %s\n", projectRoot)
	fmt.Println()

	// Run unit tests
	fmt.Println("🧪 Running Unit Tests...")
	unitResult := runTests("./tests/unit/...")
	
	fmt.Println()
	
	// Run integration tests
	fmt.Println("🔗 Running Integration Tests...")
	integrationResult := runTests("./tests/integration/...")
	
	fmt.Println()
	
	// Run any other tests in the project
	fmt.Println("📦 Running Main Package Tests...")
	mainResult := runTests("./...")
	
	fmt.Println()
	fmt.Println("=== Test Summary ===")
	
	if unitResult && integrationResult && mainResult {
		fmt.Println("✅ All tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some tests failed!")
		os.Exit(1)
	}
}

// getProjectRoot finds the root directory of the project by looking for go.mod
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for go.mod file going up the directory tree
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}
	
	return "", fmt.Errorf("could not find go.mod file")
}

// runTests runs go test on the specified package and returns success status
func runTests(pkg string) bool {
	fmt.Printf("Testing package: %s\n", pkg)
	
	// Run tests with verbose output and coverage
	cmd := exec.Command("go", "test", "-v", "-cover", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	
	if err != nil {
		fmt.Printf("❌ Tests failed for package %s\n", pkg)
		return false
	}
	
	fmt.Printf("✅ Tests passed for package %s\n", pkg)
	return true
}