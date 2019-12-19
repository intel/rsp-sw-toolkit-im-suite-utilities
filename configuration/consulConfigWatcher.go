/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package configuration

import (
	"fmt"
	"log"
	"time"

	"github.com/intel/rsp-sw-toolkit-im-suite-utilities/consulApi"
)

const (
	watchTimeout = time.Minute * 10
)

type Watcher struct {
	watchKey string
	consul   *consulApi.Client
}

func NewWatcher(consul *consulApi.Client, key string) (*Watcher, error) {
	if consul == nil {
		return nil, fmt.Errorf("consul can not be nil")
	}

	if key == "" {
		return nil, fmt.Errorf("key can not be empty")
	}

	watcher := Watcher{
		consul:   consul,
		watchKey: key,
	}

	return &watcher, nil
}

func (watcher *Watcher) Start(changeCallback func([]byte)) error {
	keyValuePair, err := watcher.consul.GetValue(watcher.watchKey, nil)
	if err != nil {
		return fmt.Errorf("unable to GET watch key (%s) data: %s", watcher.watchKey, err.Error())
	}

	if keyValuePair == nil {
		return fmt.Errorf("unable to GET watch key (%s) data: key is not found", watcher.watchKey)
	}

	go func(targetIndex uint64) {

		queryOptions := consulApi.QueryOptions{
			WaitIndex: targetIndex,
			WaitTime:  watchTimeout,
		}

		for {
			keyValuePair, err = watcher.consul.GetValue(watcher.watchKey, &queryOptions)
			if err != nil {
				log.Printf("Error watching %s key: %s", watcher.watchKey, err.Error())
				time.Sleep(time.Second * 10) // Assume consul is restarting, so want long wait.
				continue
			}

			if keyValuePair.ModifyIndex == targetIndex {
				// No change , so must have timed out. Try again
				continue
			}

			changeCallback(keyValuePair.Value)

			// This is required so we block waiting for the next change
			targetIndex = keyValuePair.ModifyIndex
		}
	}(keyValuePair.ModifyIndex)

	return nil
}
