/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package gojsonschema

import (
	"strings"
	"testing"
)

func TestDecodeReferenceToken(t *testing.T) {
	token := "abcd012345`~1~0~0~0"
	ret := decodeReferenceToken(token)
	encodedToken := encodeReferenceToken(ret)
	if encodedToken != token {
		t.Errorf("failed to decode reference token  %s  result %s", token, ret)
	}
}

func encodeReferenceToken(token string) string {
	step1 := strings.Replace(token, `~`, `~0`, -1)
	step2 := strings.Replace(step1, `/`, `~1`, -1)
	return step2
}
