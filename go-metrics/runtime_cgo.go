// +build cgo
// +build !appengine

package metrics

/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}
