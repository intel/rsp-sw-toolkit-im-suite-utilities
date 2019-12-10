/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package helper

import "testing"

func TestContains(t *testing.T) {
	expected := "two"
	target := []string{"one", expected, "three"}
	if !Contains(target, expected) {
		t.Errorf("Failed to find '%s' in %v", expected, target)
	}
}

func TestNotContains(t *testing.T) {
	unexpected := "four"
	target := []string{"one", "two", "three"}
	if Contains(target, unexpected) {
		t.Errorf("Unexpectedly found '%s' in %v", unexpected, target)
	}
}

func TestContainsWithInAny(t *testing.T) {
	expected := "two"
	target := []string{"one", "twenty" + expected, "three"}
	if !ContainsWithInAny(target, expected) {
		t.Errorf("Failed to find '%s' in %v", expected, target)
	}
}

func TestNotContainsWithInAny(t *testing.T) {
	expected := "two"
	target := []string{"one", "twenty", "three"}
	if ContainsWithInAny(target, expected) {
		t.Errorf("Unexpectedly found '%s' in %v", expected, target)
	}
}

func TestAreEqualTrue(t *testing.T) {
	source := []string{"one", "twenty", "three"}
	target := []string{"one", "twenty", "three"}

	if !AreEqual(source, target) {
		t.Errorf("Slices are unexpectedly not equal: %v , %v", source, target)
	}
}

func TestAreEqualFalse(t *testing.T) {

	source := []string{"one", "twenty", "three"}
	target := []string{"four", "five", "three"}

	if AreEqual(source, target) {
		t.Errorf("Slices are unexpectedly equal: %v , %v", source, target)
	}

}

func TestAreEqualTrueDiffOrder(t *testing.T) {
	source := []string{"three", "twenty", "one"}
	target := []string{"one", "twenty", "three"}

	if !AreEqual(source, target) {
		t.Errorf("Slices are unexpectedly not equal: %v , %v", source, target)
	}
}
