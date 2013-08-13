package main

import (
	"flag"
	"fmt"
	ed "github.com/beyang/oedipus"
	"log"
)

var docDir = flag.String("docDir", "", "docs directory, e.g., 'django/docs'")
var hideBody = flag.Bool("hideBody", false, "hides documentation body if true")
var sphinxCmd = flag.String("sphinxCmd", "", "path to sphinx-build command")
var recompile = flag.Bool("recompile", false, "if true, recompile docs from source")

func main() {
	flag.Parse()

	if *docDir == "" {
		log.Fatal("Must specify value for docDir")
	}

	if *sphinxCmd == "" {
		log.Fatal("Should specify sphinxCmd, e.g., 'sphinx-build'")
	}

	docs, errs := ed.GetDocs(*sphinxCmd, *docDir, true, *recompile)
	for _, doc := range docs {
		if *hideBody {
			fmt.Printf("%s(%s)\t%s:%d:%d\n", doc.Symbol, doc.Class, doc.SourceFile, doc.Start, doc.End)
		} else {
			fmt.Printf("======== %s(%s):%s:%d:%d ========\n%s\n\n", doc.Symbol, doc.Class, doc.SourceFile, doc.Start, doc.End, doc.Body)
		}
	}
	for _, err := range errs {
		log.Printf("Error: %v", err)
	}
}
