/*
 * INTEL CONFIDENTIAL
 * Copyright (2018) Intel Corporation.
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

package configuration

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.impcloud.net/RSP-Inventory-Suite/utilities/consulApi"

	"time"

	"github.impcloud.net/RSP-Inventory-Suite/utilities/helper"
)

const ConsulTime = time.Second * 20

func TestNewConfiguration(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	if target == nil {
		t.Error("NewConfiguration returned nil")
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	if target.parsedJson == nil {
		t.Error("Parsed JSON is missing")
	}
}

func TestLoadSimple(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	err = target.Load("./testData/simpleConfig.json")

	if err != nil {
		t.Fatalf("Config file not loaded: %s", err.Error())
	}

	if target.parsedJson == nil {
		t.Error("Parsed JSON is missing")
	}
}

func TestLoadComplex(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	err = target.Load("./testData/complexConfig.json")

	if err != nil {
		t.Fatalf("Config file not loaded: %s", err.Error())
	}

	if target.parsedJson == nil {
		t.Error("Parsed JSON is missing")
	}
}

func TestLoadNestedConfig(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	err = target.Load("./testData/nestedConfig.json")

	if err != nil {
		t.Fatalf("Config file not loaded: %s", err.Error())
	}

	if target.parsedJson == nil {
		t.Error("Parsed JSON is missing")
	}
}

func TestLoadBadFile(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	err = target.Load("./testData/badConfig.json")

	if err == nil {
		t.Error("Expected an error, but didn't get it")
	}
}

func TestLoadMissingFile(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	err = target.Load("noconfig.json")

	if err == nil {
		t.Error("Expected error, but didn't get it")
	}
}

func TestGetParsedJson(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	err = target.Load("./testData/complexConfig.json")
	if err != nil {
		t.Fatalf("Config file not loaded: %s", err.Error())
	}

	parsedJson := target.GetParsedJson()
	if parsedJson == nil {
		t.Fatal("Parsed JSON is missing")
	}

	if !target.pathExistsInConfigFile("list") {
		t.Fatal("Parsed JSON missing expected 'list' field")
	}
}

func TestGetStringSimple(t *testing.T) {
	expected := "RRP"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	actual, err := target.GetString("name")

	if err != nil {
		t.Errorf("GetString returned error for 'name': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'name' is incorrect. Expected='%s', Actual='%s'", expected, actual)
	}
}

func TestGetStringComplex(t *testing.T) {
	expected := "Arizona"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	actual, err := target.GetString("complex.location")

	if err != nil {
		t.Errorf("GetString returned error for 'complex.location': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'complex.location' is incorrect. Expected='%s', Actual='%s'", expected, actual)
	}
}

func TestGetNestedJSON(t *testing.T) {
	expectedMap := map[string]interface{}{
		"face-smile":      map[string]interface{}{"image": "face-smile:image", "fps": float64(5), "floatValue": 26.53},
		"people-counting": map[string]interface{}{"image": "people-counting:image", "fps": float64(7), "isTrue": true}}
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	target.Load("./testData/nestedConfig.json")

	mapOfInterface, err := target.GetNestedJSON("algo")

	if err != nil {
		t.Errorf("GetNestedJSON returned error for 'complex.algo': %s", err.Error())
	}

	if len(mapOfInterface) != 2 {
		t.Errorf("Expected number of values in map  is %d, but found %d", 2, len(mapOfInterface))
	}

	equals := reflect.DeepEqual(expectedMap, mapOfInterface)
	if !equals {
		t.Errorf("expectedMap %s and converted mapOfInterface %s are not equal", expectedMap, mapOfInterface)
	}
}

func TestGetNestedMapOfMapString(t *testing.T) {
	expectedMap := map[string]map[string]string{
		"face-smile":      {"image": "face-smile:image", "fps": "5", "floatValue": "26.53"},
		"people-counting": {"image": "people-counting:image", "fps": "7", "isTrue": "true"}}
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	target.Load("./testData/nestedConfig.json")

	mapOfMapOfString, err := target.GetNestedMapOfMapString("algo")

	if err != nil {
		t.Errorf("GetNestedMapOfMapString returned error for 'algo': %s", err.Error())
	}

	if len(mapOfMapOfString) != 2 {
		t.Errorf("Expected number of values in map  is %d, but found %d", 2, len(mapOfMapOfString))
	}

	equals := reflect.DeepEqual(expectedMap, mapOfMapOfString)
	if !equals {
		t.Errorf("expectedMap %v and converted mapOfMapOfString %v are not equal", expectedMap, mapOfMapOfString)
	}
}

func TestGetStringEnv(t *testing.T) {
	expected := "Testing"
	key := "UNIT_TEST_STRING"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, expected)
	actual, err := target.GetString(key)

	if err != nil {
		t.Errorf("GetString returned error for '%s': %s", key, err.Error())
	}

	if actual != expected {
		t.Errorf("Environment variable '%s' incorrect: Expected='%s', Actual='%s'", key, expected, actual)
	}
}

func TestGetStringBadSimplePath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetString("bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetStringBadComplexPath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetString("complex.bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetStringWrongType(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetString("id")

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetIntSimple(t *testing.T) {
	expected := 999
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	actual, err := target.GetInt("id")

	if err != nil {
		t.Errorf("GetInt returned error for 'id': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'id' is incorrect. Expected='%d', Actual='%d'", expected, actual)
	}
}

func TestGetIntComplex(t *testing.T) {
	expected := 1209
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	actual, err := target.GetInt("complex.streetNumber")

	if err != nil {
		t.Errorf("GetInt returned error for 'complex.streetNumber': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'complex.streetNumber' is incorrect. Expected='%d', Actual='%d'", expected, actual)
	}
}

func TestGetIntEnv(t *testing.T) {
	expected := 777
	key := "UNIT_TEST_INT"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, strconv.Itoa(expected))
	actual, err := target.GetInt(key)

	if err != nil {
		t.Errorf("Get Int returned error for '%s' environment variable: %s", key, err.Error())
	}

	if actual != expected {
		t.Errorf("Environment variable '%s' incorrect: Expected='%d', Actual='%d'", key, expected, actual)
	}
}

func TestGetIntBadSimplePath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetInt("bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetIntBadComplexPath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetInt("complex.bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetIntWrongType(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetInt("name")

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetIntEnvWrongType(t *testing.T) {
	key := "UNIT_TEST_INT"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, "bogus")
	_, err = target.GetInt(key)

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetFloatSimple(t *testing.T) {
	expected := 568.99
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	actual, err := target.GetFloat("price")

	if err != nil {
		t.Errorf("GetInt returned error for 'id': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'id' is incorrect. Expected='%f', Actual='%f'", expected, actual)
	}
}

func TestGetFloatComplex(t *testing.T) {
	expected := 568.99
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	actual, err := target.GetFloat("complex.price")

	if err != nil {
		t.Errorf("GetInt returned error for 'complex.price': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'complex.price' is incorrect. Expected='%f', Actual='%f'", expected, actual)
	}
}

func TestGetFloatEnv(t *testing.T) {
	expected := 77.70
	key := "UNIT_TEST_FLOAT"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, "77.70")
	actual, err := target.GetFloat(key)

	if err != nil {
		t.Errorf("Get Int returned error for '%s' environment variable: %s", key, err.Error())
	}

	if actual != expected {
		t.Errorf("Environment variable '%s' incorrect: Expected='%f', Actual='%f'", key, expected, actual)
	}
}

func TestGetFloatBadSimplePath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetFloat("bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetFloatBadComplexPath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetFloat("complex.bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetFloatWrongType(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetFloat("name")

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetFloatEnvWrongType(t *testing.T) {
	key := "UNIT_TEST_FLOAT"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, "bogus")
	_, err = target.GetFloat(key)

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetBoolSimple(t *testing.T) {
	expected := true
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	actual, err := target.GetBool("bool")

	if err != nil {
		t.Errorf("GetInt returned error for 'id': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'id' is incorrect. Expected='%v', Actual='%v'", expected, actual)
	}
}

func TestGetBoolComplex(t *testing.T) {
	expected := true
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	actual, err := target.GetBool("complex.isHome")

	if err != nil {
		t.Errorf("GetInt returned error for 'complex.isHome': %s", err.Error())
	}

	if actual != expected {
		t.Errorf("Value for 'complex.isHome' is incorrect. Expected='%v', Actual='%v'", expected, actual)
	}
}

func TestGetBoolEnv(t *testing.T) {
	expected := true
	key := "UNIT_TEST_BOOL"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, "true")
	actual, err := target.GetBool(key)

	if err != nil {
		t.Errorf("Get Int returned error for '%s' environment variable: %s", key, err.Error())
	}

	if actual != expected {
		t.Errorf("Environment variable '%s' incorrect: Expected='%v', Actual='%v'", key, expected, actual)
	}
}

func TestGetBoolBadSimplePath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetBool("bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetBoolBadComplexPath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetBool("complex.bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestGetBoolWrongType(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/simpleConfig.json")

	_, err = target.GetBool("id")

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetBoolEnvWrongType(t *testing.T) {
	key := "UNIT_TEST_BOOL"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	os.Setenv(key, strconv.Itoa(123))
	_, err = target.GetBool(key)

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Error("Expected error not returned")
	}
}

func TestGetStringSlice(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	expected := []string{"one", "two", "three"}
	actual, err := target.GetStringSlice("list")

	if err != nil {
		t.Fatalf("GetStringSlice returned error for 'list': %s", err.Error())
	}

	if actual == nil {
		t.Fatal("Slice not found")
	}

	if len(expected) != len(actual) {
		t.Fatalf("Slice is wrong length. Expected='%v', Actual='%v'", expected, actual)
	}

	for _, item := range expected {
		if !helper.Contains(actual, item) {
			t.Errorf("Acutal slice is missing %s. Expected='%v', Actual='%v'", item, expected, actual)
		}
	}
}

func TestGetStringSliceEnv(t *testing.T) {
	key := "UNIT_TEST_LIST"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	os.Setenv(key, "one,two,three")

	expected := []string{"one", "two", "three"}
	actual, err := target.GetStringSlice(key)

	if err != nil {
		t.Errorf("GetStringSlice returned error for 'list': %s", err.Error())
	}

	if actual == nil {
		t.Fatal("Slice not found")
	}

	if len(expected) != len(actual) {
		t.Fatalf("Slice is wrong length. Expected='%v', Actual='%v'", expected, actual)
	}

	for _, item := range expected {
		if !helper.Contains(actual, item) {
			t.Errorf("Acutal slice is missing %s. Expected='%v', Actual='%v'", item, expected, actual)
		}
	}
}

func TestGetStringSliceEnvNeedsTrim(t *testing.T) {
	key := "UNIT_TEST_LIST"
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}

	os.Setenv(key, " one  , two  , three ")

	expected := []string{"one", "two", "three"}
	actual, err := target.GetStringSlice(key)

	if err != nil {
		t.Errorf("GetStringSlice returned error for 'list': %s", err.Error())
	}

	if actual == nil {
		t.Fatal("Slice not found")
	}

	if len(expected) != len(actual) {
		t.Fatalf("Slice is wrong length. Expected='%v', Actual='%v'", expected, actual)
	}

	for _, item := range expected {
		if !helper.Contains(actual, item) {
			t.Errorf("Acutal slice is missing %s. Expected='%v', Actual='%v'", item, expected, actual)
		}
	}
}

func TestGetStringSliceBadType(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetStringSlice("numbers")

	if err == nil || !strings.Contains(err.Error(), "unable to convert") {
		t.Fatal("Didn't get expected error")
	}
}

func TestGetStringSliceBadPath(t *testing.T) {
	target, err := NewConfiguration()

	if err != nil {
		t.Errorf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/complexConfig.json")

	_, err = target.GetStringSlice("bogus")

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestSectionedConfigServiceSpecific(t *testing.T) {
	configValue := "serviceName"
	expected := "RRP Inventory Service"

	target, err := NewSectionedConfiguration("inventory-service")

	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/sectionedConfig.json")

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}

	configValue = "serverReadTimeOutSeconds"
	expectedInt := 900
	actualInt, err := target.GetInt(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actualInt != expectedInt {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestSectionedConfigNotSpecified(t *testing.T) {
	configValue := "sampleString"
	expected := "value"

	target, err := NewConfiguration()

	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/sectionedConfig.json")

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}

	configValue = "sampleInt"
	expectedInt := 123
	actualInt, err := target.GetInt(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actualInt != expectedInt {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestSectionedConfigServiceGlobal(t *testing.T) {
	configValue := "port"
	expected := "8080"

	target, err := NewSectionedConfiguration("notification-service")

	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/sectionedConfig.json")

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}

	configValue = "responseLimit"
	expectedInt := 10000
	actualInt, err := target.GetInt(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actualInt != expectedInt {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}

}

func TestSectionedConfigServiceOverrideGlobal(t *testing.T) {
	configValue := "port"
	expected := "8085"

	target, err := NewSectionedConfiguration("rules-service")

	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/sectionedConfig.json")

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestSectionedConfigFailFind(t *testing.T) {
	configValue := "bogus"

	target, err := NewSectionedConfiguration("inventory-service")

	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.Load("./testData/sectionedConfig.json")

	_, err = target.GetString(configValue)

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}

	configValue = "bogus-service.serviceName"

	_, err = target.GetString(configValue)

	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("Expected error not returned")
	}
}

func TestConfiguration_LoadConsulNoRunning(t *testing.T) {
	restoreDefaultConfiguration()
	configKey := "config/unit-test-app"
	os.Setenv("consulConfigKey", configKey)
	os.Setenv("consulUrl", "http://bogus:8500")

	_, err := NewConfiguration()
	if err == nil {
		t.Fatalf("Expected 404 error")
	}
}

func TestConfiguration_LoadConsulDefaultConfig(t *testing.T) {
	restoreDefaultConfiguration()

	configValue := "port"
	expected := "8080"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	// Ensure that Consul doesn't have the current config.
	removeConfigFromConsul(consulUrl, appConfigKey, t)

	// Loads consul with the default configuration.json file
	target, err := NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	verifyConfigInConsul(consulUrl, appConfigKey, t, nil)

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestConfiguration_LoadConsulConfigExists(t *testing.T) {
	restoreDefaultConfiguration()

	configValue := "port"
	expected := "8585"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8585\"}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	target, err := NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	verifyConfigInConsul(consulUrl, appConfigKey, t, nil)

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestConfiguration_ConsulConfigChangeApplied(t *testing.T) {
	restoreDefaultConfiguration()

	configValue := "port"
	expected := "1212"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"testing\" : {\"val1\" : 1}}"
	doneChannel := make(chan bool)
	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	target, err := NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}
	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"1212\", \"testing\" : {\"val1\" : 1}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {

		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	<-doneChannel

	actual, err := target.GetString(configValue)

	if err != nil {
		t.Fatalf("%s not found as expected", configValue)
	}

	if actual != expected {
		t.Errorf("%s actual value='%s' not as expected='%s'", configValue, actual, expected)
	}
}

func TestConfiguration_ConsulConfigGlobalChangedNotified(t *testing.T) {
	restoreDefaultConfiguration()

	section := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(section)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "port" {
			t.Fatalf("expected only port to change and got %s changed", changes[0].Name)
		}

		if changes[0].Value != "1212" {
			t.Fatalf("expected port value to change to 1212 and got %v", changes[0].Value)
		}

		if changes[0].Operation != Updated {
			t.Fatalf("expected Update change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"1212\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigGlobalAddNotified(t *testing.T) {
	restoreDefaultConfiguration()

	section := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(section)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "url" {
			t.Fatalf("expected only url to change and got %s changed", changes[0].Name)
		}

		if changes[0].Value != "localhost" {
			t.Fatalf("expected url value to change to be localhost and got %v.", changes[0].Value)
		}

		if changes[0].Operation != Added {
			t.Fatalf("expected Update change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"url\": \"localhost\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigGlobalDeleteNotified(t *testing.T) {
	restoreDefaultConfiguration()

	section := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\", \"port\": \"8080\", \"url\": \"localhost\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(section)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "url" {
			t.Fatalf("expected only url to change and got %s.", changes[0].Name)
		}

		if changes[0].Operation != Deleted {
			t.Fatalf("expected Deleted change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigTargetSectionChangedNotified(t *testing.T) {
	restoreDefaultConfiguration()

	targetSection := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\", \"num\": 100,	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 1}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(targetSection)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "unit-test.val1" {
			t.Fatalf("expected only unit-test.val1 to change and got %s changed", changes[0].Name)
		}

		if changes[0].Value != float64(9) {
			t.Fatalf("expected unit-test.val1 value to change to be 9 and got %d", changes[0].Value)
		}

		if changes[0].Operation != Updated {
			t.Fatalf("expected Update change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\", \"num\": 100,	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 9}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigTargetSectionAddNotified(t *testing.T) {
	restoreDefaultConfiguration()

	targetSection := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 1}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(targetSection)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "unit-test.url" {
			t.Fatalf("expected only unit-test.url to change and got %s changed", changes[0].Name)
		}

		if changes[0].Value != "localhost" {
			t.Fatalf("expected url value to be localhost and got %v", changes[0].Value)
		}

		if changes[0].Operation != Added {
			t.Fatalf("expected Added change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 1, \"url\": \"localhost\" }}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigTargetSectionDeleteNotified(t *testing.T) {
	restoreDefaultConfiguration()

	targetSection := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 1, \"url\": \"localhost\" }}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(targetSection)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 1 {
			t.Fatalf("expected only 1 change and got %d changes", len(changes))
		}

		if changes[0].Name != "unit-test.url" {
			t.Fatalf("expected only unit-test.url to change and got %s changed", changes[0].Name)
		}

		if changes[0].Operation != Deleted {
			t.Fatalf("expected Deleted change operation and got %v", changes[0].Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\", \"" + targetSection + "\" : {\"val1\" : 1 }}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigMultipleChangesNotified(t *testing.T) {
	restoreDefaultConfiguration()

	section := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\", \"port\": \"8080\", \"" + section + "\" : {\"val1\" : 1, \"val2\" : 2}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(section)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		if len(changes) != 4 {
			t.Fatalf("expected only 4 change and got %d changes", len(changes))
		}

		changeMap := make(map[string]ChangeDetails)

		for _, change := range changes {
			changeMap[change.Name] = change
		}

		item, exists := changeMap["port"]
		if !exists {
			t.Fatalf("expected port to change, but didn't")
		}

		if item.Value != "1212" {
			t.Fatalf("expected port value to be 1212 and got %v", changes[0].Value)
		}

		if item.Operation != Updated {
			t.Fatalf("expected port Update change operation and got %v", item.Operation)
		}

		item, exists = changeMap["unit-test.val1"]
		if !exists {
			t.Fatalf("expected  unit-test.val1 to change , but didn't")
		}

		if item.Value != float64(9) {
			t.Fatalf("expected unit-test.val1 value to change to be 9 and got %d", changes[0].Value)
		}

		if item.Operation != Updated {
			t.Fatalf("expected Update change operation and got %v", item.Operation)
		}

		item, exists = changeMap["url"]
		if !exists {
			t.Fatalf("expected url to change, but didn't")
		}

		if item.Value != "localhost" {
			t.Fatalf("expected url value to be localhost and got %v", changes[0].Value)
		}

		if item.Operation != Added {
			t.Fatalf("expected Added change operation and got %v", item.Operation)
		}

		item, exists = changeMap["unit-test.val2"]
		if !exists {
			t.Fatalf("expected unit-test.val2 to to change, but didn't")
		}

		if item.Operation != Deleted {
			t.Fatalf("expected Deleted change operation and got %v", item.Operation)
		}

		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"1212\", \"url\": \"localhost\", \"" + section + "\" : {\"val1\" : 9}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(ConsulTime)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if !callbackCalled {
		t.Errorf("never received expected config changed callback")
	}
}

func TestConfiguration_ConsulConfigNoChangedNotNotified(t *testing.T) {
	restoreDefaultConfiguration()

	targetSection := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\",  \"testing\" : {\"val1\" : 1}, \"" + targetSection + "\" : {\"val1\" : 1}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(targetSection)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	// Write new config that has no changes so it should not trigger the config changed callback.
	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\",  \"testing\" : {\"val1\" : 1}, \"" + targetSection + "\" : {\"val1\" : 1}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(time.Second * 2)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if callbackCalled {
		t.Errorf("unexpected config changed callback when No Change")
	}
}

func TestConfiguration_ConsulConfigNonTargetSectionChangedNotNotified(t *testing.T) {
	restoreDefaultConfiguration()

	targetSection := "unit-test"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\",  \"not-target\" : {\"val1\" : 1}, \"" + targetSection + "\" : {\"val1\" : 1}}"

	appConfigKey := "config/unit-test-app-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	os.Setenv("consulConfigKey", appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	callbackCalled := false
	doneChannel := make(chan bool)

	target, err := NewSectionedConfiguration(targetSection)
	if err != nil {
		t.Fatalf("NewConfiguration returned error %s", err.Error())
	}

	target.SetConfigChangeCallback(func(changes []ChangeDetails) {
		callbackCalled = true
		doneChannel <- true
	})

	// Give watcher time to spin up
	time.Sleep(time.Second * 1)

	// Write new config that has a change in non target section so it should not trigger the config changed callback.
	appConfigValue = "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8080\",  \"not-target\" : {\"val1\" : 99}, \"" + targetSection + "\" : {\"val1\" : 1}}"
	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(time.Second * 2)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel

	if callbackCalled {
		t.Errorf("unexpected config changed callback when No Change to target targetSection or globals")
	}
}

func ensureConfigInConsul(consulUrl string, configKey string, config string, t *testing.T) {
	consul, err := consulApi.NewClient(&consulApi.Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("not able to communicate with Consul service at %s: %s", consulUrl, err.Error())
	}

	err = consul.PutValue(configKey, config)
	if err != nil {
		t.Fatalf("unable to ensure config %s is in Consul service at %s: %s", configKey, consulUrl, err.Error())
	}
}

func removeConfigFromConsul(consulUrl string, configKey string, t *testing.T) {

	consul, err := consulApi.NewClient(&consulApi.Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("not able to communicate with Consul service at %s: %s", consulUrl, err.Error())
	}

	err = consul.DeleteValue(configKey)
	if err != nil {
		t.Fatalf("unable to delete config %s from Consul service at %s: %s", configKey, consulUrl, err.Error())

	}
}

func verifyConfigInConsul(consulUrl string, configKey string, t *testing.T, expected *string) {
	consul, err := consulApi.NewClient(&consulApi.Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("not able to communicate with Consul service at %s: %s", consulUrl, err.Error())
	}

	keyValuePair, err := consul.GetValue(configKey, nil)
	if err != nil {
		t.Fatalf("unable to get config %s from Consul service at %s: %s", configKey, consulUrl, err.Error())
	}

	if keyValuePair == nil {
		t.Fatalf("config %s not found in Consul service at %s", configKey, consulUrl)
	}

	if expected != nil && string(keyValuePair.Value) != *expected {
		t.Fatalf("Configuration contents in Consul not as expected. Actual %s. Expected %s", string(keyValuePair.Value), *expected)
	}
}
