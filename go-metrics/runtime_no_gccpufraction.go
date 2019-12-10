// +build !go1.5

package metrics

/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

import "runtime"

func gcCPUFraction(memStats *runtime.MemStats) float64 {
	return 0
}
