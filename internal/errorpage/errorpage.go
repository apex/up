// Package errorpage provides error page loading utilities.
package errorpage

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Page is a single .html file matching
// one or more status codes.
type Page struct {
	Name     string
	Code     int
	Range    bool
	Template *template.Template
}

// Match returns true if the page matches code.
func (p *Page) Match(code int) bool {
	switch {
	case p.Code == code:
		return true
	case p.Range && p.Code == code/100:
		return true
	case p.Name == "error" && code >= 400:
		return true
	case p.Name == "default" && code >= 400:
		return true
	default:
		return false
	}
}

// Specificity returns the specificity, where higher is more precise.
func (p *Page) Specificity() int {
	switch {
	case p.Name == "default":
		return 4
	case p.Name == "error":
		return 3
	case p.Range:
		return 2
	default:
		return 1
	}
}

// Render the page.
func (p *Page) Render(data interface{}) (string, error) {
	var buf bytes.Buffer

	if err := p.Template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Pages is a group of .html files
// matching one or more status codes.
type Pages []Page

// Match returns the matching page.
func (p Pages) Match(code int) *Page {
	for _, page := range p {
		if page.Match(code) {
			return &page
		}
	}

	return nil
}

// Load pages in dir.
func Load(dir string) (pages Pages, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "reading dir")
	}

	for _, file := range files {
		if !isErrorPage(file.Name()) {
			continue
		}

		path := filepath.Join(dir, file.Name())

		t, err := template.New(file.Name()).ParseFiles(path)
		if err != nil {
			return nil, errors.Wrap(err, "parsing template")
		}

		name := stripExt(file.Name())
		code, _ := strconv.Atoi(name)

		if isRange(name) {
			code = int(name[0] - '0')
		}

		page := Page{
			Name:     name,
			Code:     code,
			Range:    isRange(name),
			Template: t,
		}

		pages = append(pages, page)
	}

	pages = append(pages, Page{
		Name:     "default",
		Template: defaultPage,
	})

	Sort(pages)
	return
}

// Sort pages by specificity.
func Sort(pages Pages) {
	sort.Slice(pages, func(i int, j int) bool {
		a := pages[i]
		b := pages[j]
		return a.Specificity() < b.Specificity()
	})
}

// isErrorPage returns true if it looks like an error page.
func isErrorPage(path string) bool {
	if filepath.Ext(path) != ".html" {
		return false
	}

	name := stripExt(path)

	if name == "error" {
		return true
	}

	if isRange(name) {
		return true
	}

	_, err := strconv.Atoi(name)
	return err == nil
}

// isRange returns true if the name matches xx.s
func isRange(name string) bool {
	return strings.HasSuffix(name, "xx")
}

// stripExt returns path without extname.
func stripExt(path string) string {
	return strings.Replace(path, filepath.Ext(path), "", 1)
}
