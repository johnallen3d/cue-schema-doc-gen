package main

import (
	"testing"
)

func TestMissingSchemaFiles(t *testing.T) {
	_, err := gatherSchema("./test/foo", "./test/bar")

	if err == nil {
		t.Error("Expected error for missing schema got nil")
	}
}

func TestSchemaPresent(t *testing.T) {
	schema, err := gatherSchema("./test/schema", "./test/dist")

	if err != nil {
		t.Errorf("Expected no error when schema present got %v", err)
	}

	if len(schema) == 0 {
		t.Error("Expected more than zero schema to be gathered")
	}
}
