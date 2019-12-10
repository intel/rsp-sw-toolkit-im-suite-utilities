/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package metrics

import (
	"testing"
)

func TestNewRegisteredMeter(t *testing.T) {
	meter := NewRegisteredMeter("a test for new meter", nil)
	if meter == nil {
		t.Error("failed to create a new meter")
	}
}
