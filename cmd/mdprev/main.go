package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/sparrovv/go-recreation-ground/mdprev"
)

func main() {
	port := flag.String("port", "9900", "port number on which server listens")
	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(1)
	}

	mdFileName := flag.Arg(0)

	if _, err := os.Stat(mdFileName); os.IsNotExist(err) {
		log.Printf("File '%s' doesn't exist\n", mdFileName)
		os.Exit(1)
	}

	mdPrev := mdprev.NewMdPrev(mdFileName)
	mdPrev.Watch()

	go mdPrev.RunServer(*port)
	go mdPrev.ListenAndBroadcastChanges()

	url := "http://localhost:" + *port + "/" + mdPrev.MdFile
	open.Run(url) // Opens in the default browser
	fmt.Println("Server listens on:", url)

	<-mdPrev.Exit
}

func usage() {
	u := "Usage: "
	u += "mdprev [-port] filename.md"
	fmt.Println(u)
}
