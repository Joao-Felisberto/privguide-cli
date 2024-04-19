package fs_test

/*
import (
	"fmt"
	"os"
	"testing"

	"github.com/Joao-Felisberto/devprivops/fs"
)

const appName = "devprivops"

func TestGetFile(t *testing.T) {
	etcDir := fmt.Sprintf("/etc/%s", appName)
	err := os.Mkdir(etcDir, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(etcDir)

	etcFile := fmt.Sprintf("%s/%s/f1.txt", etcDir, appName)
	err = os.WriteFile(etcFile, []byte("etc"), 0666)
	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.GetFile("f1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if f != etcFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", etcFile, f)
	}

	localDir := fmt.Sprintf("../.%s", appName)
	localFile := fmt.Sprintf("%s/%s/f1.txt", localDir, appName)
	err = os.WriteFile(localFile, []byte("loc"), 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(localFile)

	f, err = fs.GetFile("f1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if f != localFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", localFile, f)
	}

	os.Remove(etcFile)

	f, err = fs.GetFile("f1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if f != localFile {
		t.Errorf("Found file does not match expectation: expected %s, got %s", localFile, f)
	}
}
*/
