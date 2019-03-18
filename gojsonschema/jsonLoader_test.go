/*
 * INTEL CONFIDENTIAL
 * Copyright (2017) Intel Corporation.
 *
 * The source code contained or described herein and all documents related to the source code ("Material")
 * are owned by Intel Corporation or its suppliers or licensors. Title to the Material remains with
 * Intel Corporation or its suppliers and licensors. The Material may contain trade secrets and proprietary
 * and confidential information of Intel Corporation and its suppliers and licensors, and is protected by
 * worldwide copyright and trade secret laws and treaty provisions. No part of the Material may be used,
 * copied, reproduced, modified, published, uploaded, posted, transmitted, distributed, or disclosed in
 * any way without Intel/'s prior express written permission.
 * No license under any patent, copyright, trade secret or other intellectual property right is granted
 * to or conferred upon you by disclosure or delivery of the Materials, either expressly, by implication,
 * inducement, estoppel or otherwise. Any license under such intellectual property rights must be express
 * and approved by Intel in writing.
 * Unless otherwise agreed by Intel in writing, you may not remove or alter this notice or any other
 * notice embedded in Materials by Intel or Intel's suppliers or licensors in any way.
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
