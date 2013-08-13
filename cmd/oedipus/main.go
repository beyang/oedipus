package main

import (
	"flag"
	ed "github.com/beyang/oedipus"
	"log"
)

var docDir = flag.String("docDir", "", "docs directory, e.g., 'django/docs'")

func main() {
	flag.Parse()

	if *docDir == "" {
		log.Fatal("Must specify value for docDir")
	}

	docs, errs := ed.GetDocs(*docDir, true)
	for _, doc := range docs {
		log.Printf("%s:\n%s\n", doc.Symbol, doc.Body)
	}
	for _, err := range errs {
		log.Printf("Error: %v", err)
	}
}
