/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package gojsonschema

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestSchemaValidationFromFile(t *testing.T) {
	filename, err := os.Getwd()
	if err != nil {
		t.Error("failed to find filename path")
	}
	schemaFilePath := path.Join(filename, "schema.json")
	jsonFilePath := path.Join(filename, "test.json")
	schemaLoader := NewReferenceLoader("file:///" + schemaFilePath)
	documentLoader := NewReferenceLoader("file:///" + jsonFilePath)

	result, err := Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatal(err)
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		t.Error("Invalid document")
	}
}

func TestLoadFromFile(t *testing.T) {
	filename, err := os.Getwd()
	if err != nil {
		t.Error("failed to find filename path")
	}
	jsonFilePath := path.Join(filename, "test.json")
	documentLoader := NewReferenceLoader("file:///" + jsonFilePath)

	contents, loadErr := documentLoader.loadFromFile(jsonFilePath)
	if loadErr != nil {
		t.Error(loadErr)
	}

	if contents == nil {
		t.Error("Empty contents found!")
	}
}
