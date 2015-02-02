package mdprev

import (
	"strings"
	"testing"
)

func TestToHTML(t *testing.T) {
	markdown := "#h1"
	html := ToHTML(markdown)

	if strings.Contains(html.String(), markdown) != true {
		t.Errorf("There's not markdown in the html template")
	}

	if strings.Contains(html.String(), "* marked - a markdown parser") != true {
		t.Errorf("Marked.js (markdown parser) is missing")
	}

	if strings.Contains(html.String(), `/*github-markdown.css*/`) != true {
		t.Errorf("github-markdown.css is missing")
	}
}
