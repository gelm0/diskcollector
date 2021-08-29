package diskcollector

import (
	"os"
	"testing"

    "github.com/stretchr/testify/assert"
)

func getNotExistingDir(path string) (string) {
    _, err := os.Stat(path)
    if err != nil {
        return path
    } else {
        return getNotExistingDir(path + "-")
    }
}

func TestStatDiskPathNotExists(t *testing.T) {
	// Assumption that noone would ever name their mountpath this (Dangerous!)
    var dirPath = "/-"
	nonExistentPath := getNotExistingDir(dirPath)
	_, err := StatDisk(nonExistentPath)
	if err == nil {
		t.Errorf("Stat %s didin't return an error", nonExistentPath)
	}
}

func TestInitDpInitializesToRoot(t *testing.T) {
	d := &UnixDiskStat{}
	d.InitDp("")
	if d.Path != "/" {
		t.Errorf("Path should be %q but is %q", "/", d.Path)
	}
}

func TestInitDpInitializesToSetPath(t *testing.T) {
	d := &UnixDiskStat{}
	path := "/tmp"
	d.InitDp(path)
	if d.Path != path {
		t.Errorf("Path should be %q but is %q", path, d.Path)
	}
}

func TestInitPdInitializesCorrectly(t *testing.T) {
    // Create temporary directory
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("Failed to create temporary directory")
	}
	defer os.RemoveAll(dir)

    diskStat := InitPd(dir)
    assert.Equal(t, diskStat.Path, dir, "Diskstat and path should be the same")
}
