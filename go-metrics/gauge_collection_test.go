package metrics

import (
	"testing"
	"time"
)

func TestAddAndReadings(t *testing.T) {
	gaugeCollection := NewGaugeCollection()
	expected := []int64{201, 569, 123}

	start := time.Now()
	for _, reading := range expected {
		time.Sleep(100 * time.Millisecond)
		gaugeCollection.Add(reading)
	}

	end := time.Now()

	actual := gaugeCollection.Readings()

	if len(actual) != len(expected) {
		t.Fatalf("Actual collection length %d not as expected %d", len(actual), len(expected))
	}

	for index := range expected {
		if actual[index].Reading != expected[index] {
			t.Errorf("Actual reading at index=%d of %d not as expected %d", index, actual[index].Reading, expected[index])
		}

		if actual[index].Time.UnixNano() < start.UnixNano() || actual[index].Time.UnixNano() > end.UnixNano() {
			t.Errorf("Actual reading's time at index=%d of %v not as expected", index, actual[index].Time)
		}
	}
}

func TestAddWithTagAndReadings(t *testing.T) {
	gaugeCollection := NewGaugeCollection()
	expectedTags := []Tag {{"Tag1", "D123"}, {"Tag2","D789"}, {"Tag3", "D034"}}
	expectedValues := []int64{201, 569, 123}

	start := time.Now()
	for index, _ := range expectedValues {
		time.Sleep(100 * time.Millisecond)
		gaugeCollection.AddWithTag(expectedValues[index], expectedTags[index])
	}

	end := time.Now()

	actual := gaugeCollection.Readings()

	if len(actual) != len(expectedValues) {
		t.Fatalf("Actual collection length %d not as expected %d", len(actual), len(expectedValues))
	}

	for index := range expectedValues {
		if actual[index].Reading != expectedValues[index] {
			t.Errorf("Actual reading at index=%d of %d not as expected %d", index, actual[index].Reading, expectedValues[index])
		}

		if actual[index].Tag != nil {
			if (*actual[index].Tag).Name != expectedTags[index].Name {
				t.Errorf("Actual tag Name at index=%d of %s not as expected %s", index, (*actual[index].Tag).Name, expectedTags[index].Name)
			}
			if (*actual[index].Tag).Value != expectedTags[index].Value {
				t.Errorf("Actual tag Value at index=%d of %s not as expected %s", index, (*actual[index].Tag).Value, expectedTags[index].Value)
			}
		} else {
			t.Errorf("Actual tag at index=%d is nil", index)
		}

		if actual[index].Time.UnixNano() < start.UnixNano() || actual[index].Time.UnixNano() > end.UnixNano() {
			t.Errorf("Actual reading's time at index=%d of %v not as expected", index, actual[index].Time)
		}
	}
}
func TestClearAndIsSet(t *testing.T) {
	gaugeCollection := NewGaugeCollection()

	if gaugeCollection.IsSet() {
		t.Fatal("Actual collection is already IsSet before adding readings")
	}

	gaugeCollection.Add(201)
	gaugeCollection.Add(301)

	actual := gaugeCollection.Readings()

	if len(actual) < 1 {
		t.Fatalf("Actual collection length %d not as expected > 0", len(actual))
	}

	if !gaugeCollection.IsSet() {
		t.Fatal("Actual collection not IsSet as expected")
	}

	gaugeCollection.Clear()
	actual = gaugeCollection.Readings()

	if len(actual) != 0 {
		t.Fatalf("Actual cleared collection length %d not as expected 0", len(actual))
	}

	if gaugeCollection.IsSet() {
		t.Fatal("Actual collection still IsSet")
	}
}

func TestSnapshot(t *testing.T) {
	original := NewGaugeCollection()
	readings := []int64{201, 569, 123, 255, 1001}

	for _, reading := range readings {
		time.Sleep(100 * time.Millisecond)
		original.Add(reading)
	}

	snapshot := original.Snapshot()
	snapShotReadings := snapshot.Readings()
	originalReadings := original.Readings()

	if snapshot.IsSet() != original.IsSet() {
		t.Errorf("Snaphot IsSet=%v not same as original IsSet=%v", snapshot.IsSet(), original.IsSet())

	}
	if len(snapShotReadings) != len(originalReadings) {
		t.Fatalf("Snapshot collection length %d not same as original %d", len(snapShotReadings), len(originalReadings))
	}

	for index := range original.Readings() {
		if snapShotReadings[index].Reading != originalReadings[index].Reading {
			t.Errorf("Actual reading at index=%d of %d not same as original %d", index, snapShotReadings[index].Reading, originalReadings[index].Reading)
		}

		if snapShotReadings[index].Time.Unix() != originalReadings[index].Time.Unix() {
			t.Errorf("Actual reading time at index=%d of %d not as original %d", index, snapShotReadings[index].Time.Unix(), originalReadings[index].Time.Unix())
		}
	}
}

func TestGetOrRegisterGaugeCollection(t *testing.T) {
	var target GaugeCollection
	expected := "Sample Name"

	target = GetOrRegisterGaugeCollection(expected, DefaultRegistry)
	if target == nil {
		t.Fatal("unexpected nil returned")
	}

	actual := DefaultRegistry.Get(expected)
	if actual == nil {
		t.Fatalf("'%s' not register as expected", expected)
	}

	switch actual.(type) {
	case GaugeCollection:
		break
	default:
		t.Fatalf("Returned type is not GaugeCollection as expected")
	}
}
