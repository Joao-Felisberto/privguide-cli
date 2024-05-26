// Tests for the fs package
// The tests are made for the internal functions that are independent of the local and global directories,
// not the exported ones, since those are easier to test.
// They have been exported in the `export_test.go` file.
package fs_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/util"
)

const (
	localRoot       = ".devprivops"  // local directory for the tests
	globalRoot      = "g_devprivops" // global directory for the tests
	descriptionRoot = "descriptions" // descriptions directory within each root
	regulationRoot  = "regulations"  // regulations directory within each root
)

// Tests for the getFile function
func TestGetFile(t *testing.T) {
	err := os.Mkdir(globalRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(globalRoot)

	err = os.Mkdir(localRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localRoot)

	etcFile := fmt.Sprintf("%s/f1.txt", globalRoot)
	err = os.WriteFile(etcFile, []byte("etc"), 0666)
	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.ExGetFile("f1.txt", localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if f != etcFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", etcFile, f)
	}

	localFile := fmt.Sprintf("%s/f1.txt", localRoot)
	err = os.WriteFile(localFile, []byte("loc"), 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localFile)

	f, err = fs.ExGetFile("f1.txt", localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if f != localFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", localFile, f)
	}

	os.Remove(etcFile)

	f, err = fs.ExGetFile("f1.txt", localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if f != localFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", localFile, f)
	}
}

// Test for the getDescriptions function.
func TestGetDescriptions(t *testing.T) {
	// Somehow this test seems not to cover some lines of the original function that are sure tested here, please revise
	err := os.Mkdir(globalRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(globalRoot)

	err = os.Mkdir(localRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localRoot)

	globalDescriptions := fmt.Sprintf("%s/%s", globalRoot, descriptionRoot)
	localDescriptions := fmt.Sprintf("%s/%s", localRoot, descriptionRoot)
	err = os.Mkdir(globalDescriptions, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(globalDescriptions)

	err = os.Mkdir(localDescriptions, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localDescriptions)

	descs, err := fs.ExGetDescriptions(descriptionRoot, localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(descs) != 0 {
		t.Errorf("There should be no descriptions, found: %s", descs)
	}

	expectedLocal := []string{"a", "b", "c"}
	for _, v := range expectedLocal {
		d := fmt.Sprintf("%s/%s", localDescriptions, v)
		err = os.Mkdir(d, 0766)
		if err != nil {
			t.Fatal(err)
		}
		// defer os.Remove(d)
	}

	descs, err = fs.ExGetDescriptions(descriptionRoot, localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	expectedLocal = util.Map(expectedLocal, func(d string) string { return fmt.Sprintf("%s/%s", descriptionRoot, d) })
	if !reflect.DeepEqual(descs, expectedLocal) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expectedLocal, descs)
	}

	expectedGlobal := []string{"d", "e", "f"}
	for _, v := range expectedGlobal {
		d := fmt.Sprintf("%s/%s", localDescriptions, v)
		err = os.Mkdir(d, 0766)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(d)
	}

	descs, err = fs.ExGetDescriptions(descriptionRoot, localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	expected := expectedGlobal
	expected = util.Map(expectedGlobal, func(d string) string { return fmt.Sprintf("%s/%s", descriptionRoot, d) })
	expected = append(expected, expectedLocal...)
	if !util.CompareSets(descs, expected) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expected, descs)
	}

	for _, v := range expectedLocal {
		d := fmt.Sprintf("%s/%s", localRoot, v)
		os.Remove(d)
	}

	descs, err = fs.ExGetDescriptions(descriptionRoot, localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	expectedGlobal = util.Map(expectedGlobal, func(d string) string { return fmt.Sprintf("%s/%s", descriptionRoot, d) })
	if !reflect.DeepEqual(descs, expectedGlobal) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expectedGlobal, descs)
	}
}

// Tests for the getRegulations function
func TestGetRegulations(t *testing.T) {
	err := os.Mkdir(globalRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(globalRoot)

	err = os.Mkdir(localRoot, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localRoot)

	globalRegulations := fmt.Sprintf("%s/%s", globalRoot, regulationRoot)
	localRegulations := fmt.Sprintf("%s/%s", localRoot, regulationRoot)
	err = os.Mkdir(globalRegulations, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(globalRegulations)

	err = os.Mkdir(localRegulations, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localRegulations)

	descs, err := fs.ExGetRegulations(localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(descs) != 0 {
		t.Errorf("There should be no regulations, found: %s", descs)
	}

	expectedLocal := []string{"a", "b", "c"}
	for _, v := range expectedLocal {
		d := fmt.Sprintf("%s/%s", localRegulations, v)
		err = os.Mkdir(d, 0766)
		if err != nil {
			t.Fatal(err)
		}
		// defer os.Remove(d)
	}

	descs, err = fs.ExGetRegulations(localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(descs, expectedLocal) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expectedLocal, descs)
	}

	expectedGlobal := []string{"d", "e", "f"}
	for _, v := range expectedGlobal {
		d := fmt.Sprintf("%s/%s", localRegulations, v)
		err = os.Mkdir(d, 0766)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(d)
	}

	descs, err = fs.ExGetRegulations(localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	expected := expectedGlobal
	expected = append(expected, expectedLocal...)
	if !util.CompareSets(descs, expected) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expected, descs)
	}

	for _, v := range expectedLocal {
		d := fmt.Sprintf("%s/%s", localRegulations, v)
		os.Remove(d)
	}

	descs, err = fs.ExGetRegulations(localRoot, globalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(descs, expectedGlobal) {
		t.Errorf("Results did not match expectations: expected %s, got %s", expectedGlobal, descs)
	}
}
