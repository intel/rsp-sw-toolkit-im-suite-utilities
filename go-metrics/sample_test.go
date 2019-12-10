/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package metrics

import "testing"

func TestExpDecaySampleHeapDown(t *testing.T) {
	expSample := newExpDecaySampleHeap(5)

	expSample.Push(expDecaySample{})
	expSample.down(expSample.Size())
	if expSample.Size() != 1 {
		t.Errorf("size wrong after down: %d", expSample.Size())
	}
}
