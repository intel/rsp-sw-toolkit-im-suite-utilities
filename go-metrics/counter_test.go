/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package metrics

import "testing"

func TestNewRegisteredCounter(t *testing.T) {
	counter := NewRegisteredCounter("a test for new counter", nil)
	if counter == nil {
		t.Error("failed to create a new counter")
	}
}
