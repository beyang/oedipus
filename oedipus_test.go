package oedipus

import (
	"testing"
)

func TestExtractDocsFromHtml(t *testing.T) {
	testcases := []struct {
		html string
		docs []Doc
	}{
		{
			html: `<body>
<dl class="class">
<dt id="foo">

<dl class="method">
<dt id="foo.bar">
</dl>

<dl class="method">
<dt id="foo.baz">
</dl>

<dl class="data">
<dt>
</dl>

<dl>
<dt>
</dl>

</dl>
</body>
`,
			docs: []Doc{
				{
					Symbol:     "foo.bar",
					Class:      "method",
					SourceFile: "test.html",
					Body: `<dl class="method">
<dt id="foo.bar">
</dl>`,
				},
				{
					Symbol:     "foo.baz",
					Class:      "method",
					SourceFile: "test.html",
					Body: `<dl class="method">
<dt id="foo.baz">
</dl>`,
				},
				{
					Symbol:     "foo",
					Class:      "class",
					SourceFile: "test.html",
					Body: `<dl class="class">
<dt id="foo">

<dl class="method">
<dt id="foo.bar">
</dl>

<dl class="method">
<dt id="foo.baz">
</dl>

<dl class="data">
<dt>
</dl>

<dl>
<dt>
</dl>

</dl>`,
				},
			},
		},
	}

	for _, testcase := range testcases {
		parser := &docParseState{}
		docs, errs := parser.extractDocsFromHtml(testcase.html, "test.html", true)

		if len(errs) != 0 {
			t.Errorf("Unexpected errors: %v", errs)
		} else {
			assertDocEqual(t, testcase.docs, docs)
		}
	}
}

func assertDocEqual(t *testing.T, exp []Doc, act []Doc) {
	if len(exp) != len(act) {
		t.Errorf("Expected %d docs, but found %d: %+v", len(exp), len(act), act)
		return
	}

	for i := 0; i < len(exp); i++ {
		if exp[i].Symbol != act[i].Symbol {
			t.Errorf("Expected doc Symbol %s, but got %s", exp[i].Symbol, act[i].Symbol)
		}
		if exp[i].Class != act[i].Class {
			t.Errorf("Expected doc Class %s, but got %s", exp[i].Class, act[i].Class)
		}
		if exp[i].SourceFile != act[i].SourceFile {
			t.Errorf("Expected doc SourceFile %s, but got %s", exp[i].SourceFile, act[i].SourceFile)
		}
		if exp[i].Body != act[i].Body {
			t.Errorf("Expected doc Body\n%s\nbut got\n%s", exp[i].Body, act[i].Body)
		}
	}
}
