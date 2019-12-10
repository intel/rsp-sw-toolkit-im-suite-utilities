/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package metrics

import "testing"

func TestGetOrRegisterEmptyFunc(t *testing.T) {
	r := NewRegistry()
	empty := ""
	if metric := r.GetOrRegister("a test for register", empty); metric == nil {
		t.Error("failed to create a new metric for empty func")
	}
}

func TestUnregisterAll(t *testing.T) {
	r := NewRegistry()
	metricName := "a test for register"
	empty := ""
	r.GetOrRegister(metricName, empty)

	r.UnregisterAll()

	if metric := r.Get(metricName); metric != nil {
		t.Error("failed to unregister all metrics")
	}
}
