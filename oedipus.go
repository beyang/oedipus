package oedipus

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var sphinxCmd = "sphinx-build" // TODO: This should be set to the path to the sphinx command

type Doc struct {
	Symbol     string
	Class      string
	SourceFile string
	Start      int
	End        int
	Body       string
}

func GetDocs(docDir string, includeSource bool) ([]Doc, []error) {
	buildDir := filepath.Join(docDir, "_build", "oedipus_html")
	cacheDir := filepath.Join(docDir, "_build", "doctrees")
	if _, err := os.Lstat(buildDir); err != nil {
		err := buildDocs(docDir, buildDir, cacheDir)
		if err != nil {
			return nil, []error{err}
		}
	}

	return extractDocs(buildDir, includeSource)
}

func buildDocs(sourceDir, buildDir, cacheDir string) error {
	cmd := exec.Command(sphinxCmd, "-b", "html", "-d", cacheDir, sourceDir, buildDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error: %v with output:%s\n", err, string(out))
	}
	return nil
}

func extractDocs(buildDir string, includeSource bool) (docs []Doc, errs []error) {
	filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !info.IsDir() && ext == ".html" {
			if b, err := ioutil.ReadFile(path); err == nil {
				log.Printf("Processing %s", path)

				body := string(b)
				parser := new(docParseState)
				theseDocs, theseErrs := parser.extractDocsFromHtml(body, "", includeSource)

				docs = append(docs, theseDocs...)
				errs = append(errs, theseErrs...)
			} else {
				errs = append(errs, err)
			}
		}
		return nil
	})
	return docs, errs
}

type docParseState struct {
	index int

	nextClassStart []int
	nextClassEnd   []int
	nextId         []int

	Stack []Doc
}

var classStartTag = regexp.MustCompile(`\<dl(?: class="([^"]*)")?\>`)
var classEndTag = regexp.MustCompile(`\</dl\>`)
var idTag = regexp.MustCompile(`\<dt id="([^"]*)"\>`)

func (d *docParseState) extractDocsFromHtml(html string, sourceFile string, includeSource bool) (docs []Doc, errs []error) {
	d.Stack = make([]Doc, 0)
	d.index = 0
	d.nextClassStart = classStartTag.FindStringSubmatchIndex(html[d.index:])
	d.nextClassEnd = classEndTag.FindStringIndex(html[d.index:])
	d.nextId = idTag.FindStringSubmatchIndex(html[d.index:])
	for {
		log.Printf("STATE: %+v", d)
		if len(d.nextClassStart) > 0 &&
			(len(d.nextClassEnd) == 0 || d.nextClassStart[0] < d.nextClassEnd[0]) &&
			(len(d.nextId) == 0 || d.nextClassStart[0] < d.nextId[0]) {

			if len(d.nextClassStart) == 4 {
				if d.nextClassStart[3] >= 0 {
					// log.Printf("%v", d.nextClassStart)
					// log.Printf("%s", html[d.nextClassStart[0]:d.nextClassStart[1]])

					className := html[d.nextClassStart[2]:d.nextClassStart[3]]
					d.Stack = append(d.Stack, Doc{Class: className, SourceFile: sourceFile, Start: d.nextClassStart[0]})
					// log.Printf("<dl>: %s", className)
				} else {
					d.Stack = append(d.Stack, Doc{SourceFile: sourceFile, Start: d.nextClassStart[0]})
				}
			}

			d.index = d.nextClassStart[1]
			d.nextClassStart = classStartTag.FindStringSubmatchIndex(html[d.index:])
			for i := 0; i < len(d.nextClassStart); i += 1 {
				d.nextClassStart[i] += d.index
			}
		} else if len(d.nextClassEnd) > 0 &&
			(len(d.nextId) == 0 || d.nextClassEnd[0] < d.nextId[0]) {
			// log.Printf("</dl>")

			doc := d.Stack[len(d.Stack)-1]
			if doc.Class == "function" || doc.Class == "method" || doc.Class == "attribute" || doc.Class == "class" {
				doc.End = d.nextClassEnd[1]
				doc.Body = html[doc.Start:doc.End]
				docs = append(docs, doc)
			}
			d.Stack = d.Stack[:len(d.Stack)-1]
			d.index = d.nextClassEnd[1]
			d.nextClassEnd = classEndTag.FindStringIndex(html[d.index:])
			for i := 0; i < len(d.nextClassEnd); i += 1 {
				d.nextClassEnd[i] += d.index
			}
		} else if len(d.nextId) > 0 {
			// log.Printf("<dt>")
			if len(d.nextId) == 4 {
				d.Stack[len(d.Stack)-1].Symbol = html[d.nextId[2]:d.nextId[3]]
			}

			d.index = d.nextId[1]
			d.nextId = idTag.FindStringSubmatchIndex(html[d.index:])
			for i := 0; i < len(d.nextId); i += 1 {
				d.nextId[i] += d.index
			}
		} else {
			break
		}

	}

	if len(d.Stack) > 0 {
		return nil, []error{fmt.Errorf("Non-empty stack at end of parsing")}
	}

	return docs, errs
}
