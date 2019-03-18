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

package consulApi

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Uncomment if you want to run tests against real Consul service running locally
// var consulUrl = "http://localhost:8500/v1/kv"
var consulUrl string

func TestMain(m *testing.M) {

	var testMockServer *httptest.Server
	if consulUrl == "" {
		mockConsul := NewMockConsul()
		testMockServer = mockConsul.Start()
		consulUrl = testMockServer.URL + "/v1/kv"
	}

	exitCode := m.Run()
	if testMockServer != nil {
		defer testMockServer.Close()
	}
	os.Exit(exitCode)
}

func TestNewClient(t *testing.T) {
	_, err := NewClient(&Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("failed to create NewClient: %s", err.Error())
	}
}

func TestGetValueNotFound(t *testing.T) {
	key := "bogus"

	target, err := NewClient(&Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("failed to create NewClient: %s", err.Error())
	}

	var keyValuePair *KeyValuePair
	keyValuePair, err = target.GetValue(key, nil)
	if err != nil {
		t.Fatalf("unexpected an error: %s", err.Error())
	}

	if keyValuePair != nil {
		t.Fatalf("expected nil keyValuePair for Not Found")
	}
}

func TestPutGetValue(t *testing.T) {
	key := "myKey"
	expected := "This is my value"

	target, err := NewClient(&Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("failed to create NewClient: %s", err.Error())
	}

	if err = target.PutValue(key, expected); err != nil {
		t.Fatalf("failed PutValue of %s to %s key: %s", expected, key, err.Error())
	}

	var keyValuePair *KeyValuePair
	keyValuePair, err = target.GetValue(key, nil)
	if err != nil {
		t.Fatalf("failed GetValue for key %s: %s", key, err.Error())
	}

	if keyValuePair == nil {
		t.Fatalf("unexpected key not found")
	}

	actual := string(keyValuePair.Value)
	if actual != expected {
		t.Fatalf("Actual value received '%s' is not as expected '%s'", actual, expected)
	}
}

func TestBlockingGetValueTimedOut(t *testing.T) {
	key := "SomeKey"
	value := "Some value"
	expectedWait := time.Second * 5
	target, err := NewClient(&Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("failed to create NewClient: %s", err.Error())
	}

	if err = target.PutValue(key, value); err != nil {
		t.Fatalf("failed PutValue of %s to %s key: %s", value, key, err.Error())
	}

	var keyValuePair *KeyValuePair
	keyValuePair, err = target.GetValue(key, nil)
	if err != nil {
		t.Fatalf("failed GetValue for key %s: %s", key, err.Error())
	}

	doneChannel := make(chan bool)
	go func(t *testing.T) {
		_, err := target.GetValue(key, &QueryOptions{WaitIndex: keyValuePair.ModifyIndex, WaitTime: expectedWait})
		if err != nil {
			t.Errorf("failed GetValue with wait for key %s: %s", key, err.Error())
		}

		doneChannel <- true
	}(t)

	startTime := time.Now()
	<-doneChannel

	actualWaited := time.Since(startTime)
	if actualWaited < expectedWait {
		t.Fatalf("Didn't wait as expected. Actual %v, Expected %v", actualWaited, expectedWait)
	}
}

func TestBlockingGetValueNotTimedOut(t *testing.T) {
	key := "AnotherKey"
	value := "Another value"
	expected := "Changed Value"
	waitTime := time.Second * 5

	target, err := NewClient(&Config{Address: consulUrl})
	if err != nil {
		t.Fatalf("failed to create NewClient: %s", err.Error())
	}

	if err = target.PutValue(key, value); err != nil {
		t.Fatalf("failed PutValue of %s to %s key: %s", value, key, err.Error())
	}

	keyValuePair, err := target.GetValue(key, nil)
	if err != nil {
		t.Fatalf("failed GetValue for key %s: %s", key, err.Error())
	}

	doneChannel := make(chan bool)
	go func(t *testing.T) {
		keyValuePair, err = target.GetValue(key, &QueryOptions{WaitIndex: keyValuePair.ModifyIndex, WaitTime: waitTime})
		if err != nil {
			t.Errorf("failed GetValue with wait for key %s: %s", key, err.Error())
		}

		doneChannel <- true
	}(t)

	go func(t *testing.T) {
		time.Sleep(time.Second * 1)

		// Put the new value to trigger GetValue to return before timeout
		if err := target.PutValue(key, expected); err != nil {
			t.Errorf("failed PutValue of %s to %s key: %s", value, key, err.Error())
		}
	}(t)

	startTime := time.Now()
	<-doneChannel

	actualWaited := time.Since(startTime)
	if actualWaited >= waitTime {
		t.Fatalf("Didn't wait as expected. Actual %v, Expected %v", actualWaited, waitTime)
	}

	actual := string(keyValuePair.Value)
	if actual != expected {
		t.Fatalf("Actual value received '%s' is not as expected '%s'", actual, expected)
	}
}
