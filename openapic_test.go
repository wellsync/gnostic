package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func test_compiler(t *testing.T, input_file string, reference_file string, expect_errors bool) {
	// output_file is the name of the text-format proto file that is generated by the compiler
	var suffix string
	if expect_errors {
		suffix = ".errors"
	} else {
		suffix = ".text"
	}
	output_file := strings.Replace(filepath.Base(input_file), filepath.Ext(input_file), suffix, 1)
	// remove any preexisting output_file
	var err error
	err = exec.Command("rm", "-f", output_file).Run()
	if err != nil {
		t.Logf("Failed to remove file %s", output_file)
		t.FailNow()
	}
	// run the compiler
	err = exec.Command("openapic", "-text_out", output_file, "-errors", input_file).Run()
	if err != nil && !expect_errors {
		t.Logf("Compile failed: %+v", err)
		t.FailNow()
	}
	// verify the output_file against the reference_file
	err = exec.Command("diff", output_file, reference_file).Run()
	if err != nil {
		t.Logf("Diff failed: %+v", err)
		t.FailNow()
	}
}

func test_normal(t *testing.T, input_file string, reference_file string) {
	test_compiler(t, input_file, reference_file, false)
}

func test_errors(t *testing.T, input_file string, reference_file string) {
	test_compiler(t, input_file, reference_file, true)
}

func TestPetstoreJSON(t *testing.T) {
	test_normal(t,
		"examples/petstore.json",
		"test/petstore.text")
}

func TestPetstoreYAML(t *testing.T) {
	test_normal(t,
		"examples/petstore.yaml",
		"test/petstore.text")
}

func TestSeparateYAML(t *testing.T) {
	test_normal(t,
		"examples/v2.0/yaml/petstore-separate/spec/swagger.yaml",
		"test/v2.0/yaml/petstore-separate/spec/swagger.text")
}

func TestSeparateJSON(t *testing.T) {
	test_normal(t,
		"examples/v2.0/json/petstore-separate/spec/swagger.json",
		"test/v2.0/yaml/petstore-separate/spec/swagger.text") // yaml and json results should be identical
}

func TestRemotePetstoreJSON(t *testing.T) {
	test_normal(t,
		"https://raw.githubusercontent.com/googleapis/openapi-compiler/master/examples/petstore.json",
		"test/petstore.text")
}

func TestRemotePetstoreYAML(t *testing.T) {
	test_normal(t,
		"https://raw.githubusercontent.com/googleapis/openapi-compiler/master/examples/petstore.yaml",
		"test/petstore.text")
}

func TestRemoteSeparateYAML(t *testing.T) {
	test_normal(t,
		"https://raw.githubusercontent.com/googleapis/openapi-compiler/master/examples/v2.0/yaml/petstore-separate/spec/swagger.yaml",
		"test/v2.0/yaml/petstore-separate/spec/swagger.text")
}

func TestRemoteSeparateJSON(t *testing.T) {
	test_normal(t,
		"https://raw.githubusercontent.com/googleapis/openapi-compiler/master/examples/v2.0/json/petstore-separate/spec/swagger.json",
		"test/v2.0/yaml/petstore-separate/spec/swagger.text")
}

func TestErrorBadProperties(t *testing.T) {
	test_errors(t,
		"examples/errors/petstore-badproperties.yaml",
		"test/errors/petstore-badproperties.errors")
}

func TestErrorUnresolvedRefs(t *testing.T) {
	test_errors(t,
		"examples/errors/petstore-unresolvedrefs.yaml",
		"test/errors/petstore-unresolvedrefs.errors")
}
