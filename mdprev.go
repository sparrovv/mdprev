package mdprev

import (
	"io/ioutil"
	"log"
	"path/filepath"

	fsnotify "gopkg.in/fsnotify.v1"
)

type MdPrev struct {
	//File name
	MdFile string
	//File contents
	MdContent string
	//Chan on which changes are pushed
	MdChanges chan bool
	//Broadcast changes to connected clients
	Broadcast chan []byte
	Exit      chan bool
}

func NewMdPrev(mdFile string) *MdPrev {
	ch := make(chan bool)
	b := make(chan []byte)
	exit := make(chan bool)
	mdPrev := &MdPrev{mdFile, "", ch, b, exit}
	mdPrev.loadContent()

	return mdPrev
}

func (m *MdPrev) loadContent() {
	body, _ := ioutil.ReadFile(m.MdFile)
	m.MdContent = string(body)
}

// observe file for changes
func (mdPrev *MdPrev) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					mdPrev.MdChanges <- true
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(mdPrev.MdFile)
	if err != nil {
		panic(err)
	}
}

func (mdPrev *MdPrev) ListenAndBroadcastChanges() {
	for _ = range mdPrev.MdChanges {
		mdPrev.loadContent()
		mdPrev.Broadcast <- []byte(mdPrev.MdContent)
	}
}

func (mdPrev *MdPrev) MdDirPath() string {
	dir, err := filepath.Abs(filepath.Dir(mdPrev.MdFile))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
