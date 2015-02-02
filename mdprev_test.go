package mdprev

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewMdPrev(t *testing.T) {
	expectedContent := "Content"
	mdFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(mdFile.Name())

	err = ioutil.WriteFile(mdFile.Name(), []byte(expectedContent), 0644)
	if err != nil {
		panic(err)
	}

	mdPrev := NewMdPrev(mdFile.Name())
	if mdPrev.MdContent != expectedContent {
		t.Fatalf("The file's contents: %s is not eql to the expected %s", mdPrev.MdContent, expectedContent)
	}
}

func TestWatcher(t *testing.T) {
	expectedContent := "Content"
	mdPrev := buildTestMdPrev("")
	defer os.Remove(mdPrev.MdFile)

	mdPrev.Watch()

	var expectedChange bool
	expectedChange = false

	go func() {
		expectedChange = <-mdPrev.MdChanges
	}()

	_ = ioutil.WriteFile(mdPrev.MdFile, []byte(expectedContent), 0644)
	// it sucks, but watcher needs time to notice changes
	time.Sleep(10 * time.Millisecond)

	if expectedChange != true {
		t.Errorf("Expected that watcher will notify about the change, but it didn't")
	}
}

func TestListenAndBraodcastChanges(t *testing.T) {
	content := "A lot of new content"
	mdPrev := buildTestMdPrev(content)
	defer os.Remove(mdPrev.MdFile)
	var b []byte

	go mdPrev.ListenAndBroadcastChanges()
	mdPrev.MdChanges <- true

	b = <-mdPrev.Broadcast
	if string(b) != content {
		t.Errorf("Expecting to broadcast new content but got %s", string(b))
	}
}

func TestMdDirPath(t *testing.T) {
	mdPrev := buildTestMdPrev("")
	if mdPrev.MdDirPath() != filepath.Dir(mdPrev.MdFile) {
		t.Errorf("Expecting an absolute path the directory containg md file, but got %s", mdPrev.MdDirPath())
	}

}

func buildTestMdPrev(content string) *MdPrev {
	mdFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(mdFile.Name(), []byte(content), 0644)
	return NewMdPrev(mdFile.Name())
}
