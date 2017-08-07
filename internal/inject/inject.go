// Package inject provides script and style injection utilities.
package inject

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"strings"

	"github.com/apex/log"
	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

// TODO: template support
// TODO: move config to "config" pkg

// locations valid.
var locations = []string{
	"head",
	"body",
}

// types valid.
var types = []string{
	"literal",
	"comment",
	"style",
	"script",
	"inline style",
	"inline script",
	"google analytics",
	"segment",
}

// Rules is a set of rules mapped by location.
type Rules map[string][]*Rule

// Default rules.
func (r Rules) Default() error {
	for pos, rules := range r {
		for i, rule := range rules {
			if err := rule.Default(); err != nil {
				return errors.Wrapf(err, "%s rule #%d", pos, i+1)
			}
		}
	}

	return nil
}

// Validate rules.
func (r Rules) Validate() error {
	for pos, rules := range r {
		if err := validate.List(pos, locations); err != nil {
			return errors.Wrap(err, "invalid location")
		}

		for i, rule := range rules {
			if err := rule.Validate(); err != nil {
				return errors.Wrapf(err, "%s rule #%d", pos, i+1)
			}
		}
	}

	return nil
}

// Apply rules to html.
func (r Rules) Apply(html string) string {
	for pos, rules := range r {
		log.Debugf("injecting %s rules", pos)
		for _, rule := range rules {
			log.Debugf("  inject %s %q", rule.Type, rule.Value)
			switch pos {
			case "head":
				html = Head(html, rule.Apply(html))
			case "body":
				html = Body(html, rule.Apply(html))
			}
		}
	}

	return html
}

// Rule is an injection rule.
type Rule struct {
	// Type of injection, defaults to "literal" unless File is used,
	// or Value contains .js or .css extensions.
	Type string `json:"type"`

	// Value is the literal, inline string, or src/href of the injected tag.
	Value string `json:"value"`

	// File is used to load source from disk instead of providing Value. Note
	// that if Type is not explicitly provided, then it will default to
	// "inline script" or "inline style" for .js and .css files respectively.
	File string `json:"file"`
}

// Apply rule to html.
func (r *Rule) Apply(html string) string {
	switch r.Type {
	case "literal":
		return r.Value
	case "script":
		return Script(r.Value)
	case "style":
		return Style(r.Value)
	case "inline script":
		return ScriptInline(r.Value)
	case "inline style":
		return StyleInline(r.Value)
	case "comment":
		return Comment(r.Value)
	case "segment":
		return Segment(r.Value)
	case "google analytics":
		return GoogleAnalytics(r.Value)
	default:
		return ""
	}
}

// Default applies defaults.
func (r *Rule) Default() error {
	if r.Type == "" {
		r.Type = "literal"
	}

	if r.File != "" {
		if err := r.defaultFile(r.File); err != nil {
			return err
		}
	}

	return nil
}

// Validate returns an error if incorrect.
func (r *Rule) Validate() error {
	if err := validate.List(r.Type, types); err != nil {
		return errors.Wrap(err, "invalid .type")
	}

	if strings.TrimSpace(r.Value) == "" {
		return errors.Errorf(`.value is required`)
	}

	return nil
}

// defaultFile defaults value" from the given path.
func (r *Rule) defaultFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	r.Value = string(b)
	return nil
}

// Head injects a string before the closing head tag.
func Head(html, s string) string {
	return strings.Replace(html, "</head>", "  "+s+"\n  </head>", 1)
}

// Body injects a string before the closing body tag.
func Body(html, s string) string {
	return strings.Replace(html, "</body>", "  "+s+"\n  </body>", 1)
}

// Script returns an script.
func Script(src string) string {
	return `<script src="` + html.EscapeString(src) + `"></script>`
}

// ScriptInline returns an inline script.
func ScriptInline(s string) string {
	return `<script>` + s + `</script>`
}

// Style returns an style.
func Style(href string) string {
	return `<link rel="stylesheet" href="` + html.EscapeString(href) + `">`
}

// StyleInline returns an inline style.
func StyleInline(s string) string {
	return `<style>` + s + `</style>`
}

// Comment returns an html comment.
func Comment(s string) string {
	return "<!-- " + html.EscapeString(s) + " -->"
}

// Segment inline script with key.
func Segment(key string) string {
	return ScriptInline(`
  !function(){var analytics=window.analytics=window.analytics||[];if(!analytics.initialize)if(analytics.invoked)window.console&&console.error&&console.error("Segment snippet included twice.");else{analytics.invoked=!0;analytics.methods=["trackSubmit","trackClick","trackLink","trackForm","pageview","identify","reset","group","track","ready","alias","debug","page","once","off","on"];analytics.factory=function(t){return function(){var e=Array.prototype.slice.call(arguments);e.unshift(t);analytics.push(e);return analytics}};for(var t=0;t<analytics.methods.length;t++){var e=analytics.methods[t];analytics[e]=analytics.factory(e)}analytics.load=function(t){var e=document.createElement("script");e.type="text/javascript";e.async=!0;e.src=("https:"===document.location.protocol?"https://":"http://")+"cdn.segment.com/analytics.js/v1/"+t+"/analytics.min.js";var n=document.getElementsByTagName("script")[0];n.parentNode.insertBefore(e,n)};analytics.SNIPPET_VERSION="4.0.0";
  analytics.load("` + key + `");
  analytics.page();
  }}();
`)
}

// GoogleAnalytics inline script with tracking key.
func GoogleAnalytics(trackingID string) string {
	return ScriptInline(`
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', '` + trackingID + `', 'auto');
  ga('send', 'pageview');
`)
}

// Var injection.
func Var(kind, name string, v interface{}) string {
	b, _ := json.Marshal(v)
	return ScriptInline(kind + ` ` + name + ` = ` + string(b))
}
