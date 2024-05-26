// Tests for the utils package
package util_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/Joao-Felisberto/devprivops/util"
)

// Test for the Map function
func TestMap(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	mapped := util.Map(nums, func(n int) int { return 2 * n })

	for i := range expected {
		if expected[i] != mapped[i] {
			t.Errorf("Arrays do not match: expected %d, got %d", expected[i], mapped[i])
		}
	}
}

// Test for the Filter function
func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	expected := []int{2, 4}
	mapped := util.Filter(nums, func(n int) bool { return n%2 == 0 })

	for i := range expected {
		if expected[i] != mapped[i] {
			t.Errorf("Arrays do not match: expected %d, got %d", expected[i], mapped[i])
		}
	}
}

// Test for the MapCast function
func TestMapCast(t *testing.T) {
	m := map[interface{}]interface{}{
		"a": 1,
		"b": 2,
	}

	mMapped := util.MapCast[string, int](m)

	if reflect.TypeOf(mMapped) != reflect.TypeOf(map[string]int{}) {
		t.Errorf("The casting failed! Expected %#v, got %#v", map[string]int{}, mMapped)
	}
}

// Test for the MapToMap function
func TestMapToMap(t *testing.T) {
	stringIntPairs := []struct {
		s string
		i int
	}{
		{"a", 1},
		{"b", 2},
	}

	stringIntMap := util.ArrayToMap(stringIntPairs, func(p struct {
		s string
		i int
	}) (string, int) {
		return p.s, p.i
	})

	if !reflect.DeepEqual(stringIntMap, map[string]int{"a": 1, "b": 2}) {
		t.Errorf("The cast map did not match expectations, got %v\n", stringIntMap)
	}
}

// Test for the Any function
func TestAny(t *testing.T) {
	nums := []int{2, 4, 6, 8}
	found := util.Any(nums, func(n int) bool { return n%2 == 0 })

	if !found {
		t.Errorf("Should have found an even number in the aray.")
	}

	found = util.Any(nums, func(n int) bool { return n%2 != 0 })

	if found {
		t.Errorf("Should not have found an odd number in the aray.")
	}
}

// Test for the CompareSets function
func TestCompareSets(t *testing.T) {
	set1 := []int{1, 2, 3}
	set2 := []int{1, 2, 3}

	if !util.CompareSets(set1, set2) {
		t.Errorf("Sets should be equal: %d %d", set1, set2)
	}

	set1 = []int{1, 2}
	if util.CompareSets(set1, set2) {
		t.Errorf("Sets should be different: set1 has less elements: %d %d", set1, set2)
	}

	set1 = []int{1, 2, 1000}
	if util.CompareSets(set1, set2) {
		t.Errorf("Sets should be different: set1 has an element not present in set2: %d %d", set1, set2)
	}
	set1 = []int{1, 2, 3}

	set2 = []int{1, 2}
	if util.CompareSets(set1, set2) {
		t.Errorf("Sets should be different: set2 has less elements: %d %d", set1, set2)
	}

	set2 = []int{1, 2, 1000}
	if util.CompareSets(set1, set2) {
		t.Errorf("Sets should be different: set2 has an element not present in set2: %d %d", set1, set2)
	}
}

// Test for the CreateFileWithData function
func TestCreateFileWithData(t *testing.T) {
	for _, f := range []string{"file.txt", "./some/file.txt", "./some/dir/and/file.txt"} {
		err := util.CreateFileWithData(f, "data")
		if err != nil {
			t.Fatalf("Failed to create file '%s': %s", f, err)
		}
		defer util.DeleteFileAndParentPath(f)

		rawFileData, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("Failed to read created file '%s': %s", f, err)
		}

		fileData := string(rawFileData)
		if fileData != "data" {
			t.Errorf("The data written is not what was on the file '%s': expected 'data', got '%s'", f, fileData)
		}
	}
}
