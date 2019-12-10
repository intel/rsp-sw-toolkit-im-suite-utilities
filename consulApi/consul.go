/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package consulApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address  string
	WaitTime time.Duration
}

type Client struct {
	config Config
}

func NewClient(config *Config) (*Client, error) {
	return &Client{config: *config}, nil
}

type QueryOptions struct {
	WaitIndex uint64
	WaitTime  time.Duration
}

type KeyValuePairs []*KeyValuePair

type KeyValuePair struct {
	Key         string
	CreateIndex uint64
	ModifyIndex uint64
	LockIndex   uint64
	Flags       uint64
	Value       []byte
	Session     string
}

func (client *Client) GetValue(key string, queryOptions *QueryOptions) (*KeyValuePair, error) {
	//index=1&wait=600000ms

	endpoint := client.buildEndPoint(key)

	httpClient := &http.Client{
		Timeout: time.Second * 1800,
	}

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create GET http.NewRquest to %s: %s", endpoint, err.Error())
	}
	request.Header.Set("content-type", "application/json;charset=utf-8")

	if queryOptions != nil {
		query := request.URL.Query()
		if queryOptions.WaitIndex != 0 {
			query.Add("index", strconv.FormatUint(queryOptions.WaitIndex, 10))
		}

		if queryOptions.WaitTime != 0 {
			query.Add("wait", durToMsec(queryOptions.WaitTime))
		}

		request.URL.RawQuery = query.Encode()
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to get value for %s: %s", key, err.Error())
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if response.StatusCode == http.StatusOK {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body for %s: %s", key, err.Error())
		}

		keyValuePairs := KeyValuePairs{}

		if err := json.Unmarshal(responseData, &keyValuePairs); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response for key %s into a KeyValuePairs struct: %s", key, err.Error())
		}

		var keyValuePair *KeyValuePair = nil

		if len(keyValuePairs) > 0 {
			keyValuePair = keyValuePairs[0]
		}

		return keyValuePair, nil
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return nil, fmt.Errorf("unable to get value for %s: Received %d status", key, response.StatusCode)
}

func (client *Client) PutValue(key string, value string) error {
	endpoint := client.buildEndPoint(key)

	httpClient := &http.Client{
		Timeout: time.Second * 1800,
	}

	request, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(value)))
	if err != nil {
		return fmt.Errorf("unable to create PUT http.NewRquest to %s: %s", endpoint, err.Error())
	}
	request.Header.Set("content-type", "application/json;charset=utf-8")

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("unable to put value for %s: %s", key, err.Error())
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("unable to put value for %s: Received %d status", key, response.StatusCode)
}

func (client *Client) DeleteValue(key string) error {
	return nil
}

func durToMsec(dur time.Duration) string {
	ms := dur / time.Millisecond
	if dur > 0 && ms == 0 {
		ms = 1
	}
	return fmt.Sprintf("%dms", ms)
}

func (client *Client) buildEndPoint(key string) string {
	endpoint := client.config.Address + "/" + strings.TrimPrefix(key, "/")
	return endpoint
}
