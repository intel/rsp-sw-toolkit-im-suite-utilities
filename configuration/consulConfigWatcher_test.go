/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/intel/rsp-sw-toolkit-im-suite-utilities/consulApi"
)

const DefaultConfigurationFile = "configuration.json"

var consul *consulApi.Client

// Uncomment if you want to run tests against real Consul service running locally
//var consulUrl = "http://localhost:8500/v1/kv"
var consulUrl string
var originalBytes []byte

func TestMain(m *testing.M) {
	var readErr error
	var testMockServer *httptest.Server
	if consulUrl == "" {
		mockConsul := consulApi.NewMockConsul()
		testMockServer = mockConsul.Start()
		consulUrl = testMockServer.URL + "/v1/kv"
	}

	os.Setenv("consulUrl", consulUrl)

	originalBytes, readErr = ioutil.ReadFile(DefaultConfigurationFile)
	if readErr != nil {
		fmt.Printf("Unable to read local default %s: %s", DefaultConfigurationFile, readErr.Error())
		os.Exit(1)
	}

	var clientErr error
	consul, clientErr = consulApi.NewClient(&consulApi.Config{Address: consulUrl})

	if clientErr != nil {
		fmt.Printf("Unable to create new Consul Client: %s", clientErr.Error())
		os.Exit(1)
	}

	exitCode := m.Run()
	restoreDefaultConfiguration()

	if testMockServer != nil {
		defer testMockServer.Close()
	}

	os.Exit(exitCode)
}

func restoreDefaultConfiguration() {
	log.Print("Restoring default config contents")
	os.Setenv("consulUrl", consulUrl)
	ioutil.WriteFile(DefaultConfigurationFile, originalBytes, 0666)
}

func TestNewWatcher(t *testing.T) {
	target, err := NewWatcher(consul, "config/unit-test-config")

	if target == nil {
		t.Error("Watcher not created")
	}

	if err != nil {
		t.Errorf("NewWatcher returned error %s", err.Error())
	}
}

func TestNewWatcherInputError(t *testing.T) {
	_, err := NewWatcher(nil, "config/unit-test-config")

	if err == nil {
		t.Errorf("Expecting error for consul nil")
	}

	_, err = NewWatcher(consul, "")
	if err == nil {
		t.Errorf("Expecting error for key empty")
	}
}

func TestStart(t *testing.T) {
	appConfigKey := "config/unit-test-app"
	appConfigValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"8585\", \"testing\" : {\"val1\" : 1}}"
	expectedValue := "{\"name\" : \"Default Config Unit Test\",	\"port\": \"1212\", \"testing\" : {\"val1\" : 1}}"

	target, err := NewWatcher(consul, appConfigKey)

	ensureConfigInConsul(consulUrl, appConfigKey, appConfigValue, t)

	doneChannel := make(chan bool)

	if err != nil {
		t.Errorf("NewWatcher returned error %s", err.Error())
	}

	if target == nil {
		t.Error("Watcher not created")
	}

	if err != nil {
		t.Errorf("NewWatcher returned error %s", err.Error())
	}

	if err := target.Start(func(actualValue []byte) {
		if string(actualValue) != expectedValue {
			t.Fatalf("actual (%s) value not expected (%s)", actualValue, expectedValue)
		}
		doneChannel <- true
	}); err != nil {
		t.Errorf("Watcher not started: %s", err.Error())
	}

	ensureConfigInConsul(consulUrl, appConfigKey, expectedValue, t)

	// Give long time for change notification to timeout
	go func() {
		time.Sleep(time.Second * 10)
		doneChannel <- true
	}()

	// block until done or timed out
	<-doneChannel
}

func TestStartKeyNotFound(t *testing.T) {
	appConfigKey := "config/Bogus"

	target, _ := NewWatcher(consul, appConfigKey)
	err := target.Start(nil)
	if err == nil {
		t.Errorf("expecting error for %s key not found", appConfigKey)
	}
}

func TestStartNoConsulRunning(t *testing.T) {
	appConfigKey := "config/unit-test-app"
	badConsulUrl := "localhost:8080"

	badConsul, _ := consulApi.NewClient(&consulApi.Config{Address: badConsulUrl})

	target, _ := NewWatcher(badConsul, appConfigKey)
	err := target.Start(nil)
	if err == nil {
		t.Error("expecting error for consul not running")
	}
}
