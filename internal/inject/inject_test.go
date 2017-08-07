package inject_test

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
	"github.com/apex/up/internal/inject"
)

var html = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Example</title>
  </head>
  <body>
    <p>Hello World</p>
  </body>
</html>
`

func ExampleStyle() {
	fmt.Printf("%s\n", inject.Style(`/sloth.css`))
	// Output:
	// <link rel="stylesheet" href="/sloth.css">
}

func ExampleStyleInline() {
	fmt.Printf("%s\n", inject.StyleInline(`body { display: none }`))
	// Output:
	// <style>body { display: none }</style>
}

func ExampleScript() {
	fmt.Printf("%s\n", inject.Script(`/sloth.js`))
	// Output:
	// <script src="/sloth.js"></script>
}

func ExampleScriptInline() {
	fmt.Printf("%s\n", inject.ScriptInline(`const user = { "name": "Tobi" }`))
	// Output:
	// <script>const user = { "name": "Tobi" }</script>
}

func ExampleComment() {
	fmt.Printf("%s\n", inject.Comment(`Hello World`))
	// Output:
	// <!-- Hello World -->
}

func ExampleHead() {
	s := inject.Head(html, `<link rel="stylesheet" href="/style.css">`)
	fmt.Printf("%s\n", s)
	// Output:
	// <!doctype html>
	// <html>
	//   <head>
	//     <meta charset="utf-8">
	//     <title>Example</title>
	//     <link rel="stylesheet" href="/style.css">
	//   </head>
	//   <body>
	//     <p>Hello World</p>
	//   </body>
	// </html>
}

func ExampleBody() {
	s := inject.Body(html, inject.Comment("Version 1.0.3"))
	fmt.Printf("%s\n", s)
	// Output:
	// <!doctype html>
	// <html>
	//   <head>
	//     <meta charset="utf-8">
	//     <title>Example</title>
	//   </head>
	//   <body>
	//     <p>Hello World</p>
	//     <!-- Version 1.0.3 -->
	//   </body>
	// </html>
}

func ExampleSegment() {
	fmt.Printf("%s\n", inject.Segment(`KEY HERE`))
	// Output:
	// <script>
	//   !function(){var analytics=window.analytics=window.analytics||[];if(!analytics.initialize)if(analytics.invoked)window.console&&console.error&&console.error("Segment snippet included twice.");else{analytics.invoked=!0;analytics.methods=["trackSubmit","trackClick","trackLink","trackForm","pageview","identify","reset","group","track","ready","alias","debug","page","once","off","on"];analytics.factory=function(t){return function(){var e=Array.prototype.slice.call(arguments);e.unshift(t);analytics.push(e);return analytics}};for(var t=0;t<analytics.methods.length;t++){var e=analytics.methods[t];analytics[e]=analytics.factory(e)}analytics.load=function(t){var e=document.createElement("script");e.type="text/javascript";e.async=!0;e.src=("https:"===document.location.protocol?"https://":"http://")+"cdn.segment.com/analytics.js/v1/"+t+"/analytics.min.js";var n=document.getElementsByTagName("script")[0];n.parentNode.insertBefore(e,n)};analytics.SNIPPET_VERSION="4.0.0";
	//   analytics.load("KEY HERE");
	//   analytics.page();
	//   }}();
	// </script>
}

func ExampleGoogleAnalytics() {
	fmt.Printf("%s\n", inject.GoogleAnalytics(`KEY HERE`))
	// Output:
	// <script>
	//   (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
	//   (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
	//   m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
	//   })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');
	//
	//   ga('create', 'KEY HERE', 'auto');
	//   ga('send', 'pageview');
	// </script>
}

func ExampleVar() {
	user := map[string]string{
		"name": "Tobi",
	}

	fmt.Printf("%s\n", inject.Var("const", "user", user))
	// Output:
	// <script>const user = {"name":"Tobi"}</script>
}

func TestRule_Default(t *testing.T) {
	r := inject.Rule{Value: `<script></script>`}
	assert.NoError(t, r.Default(), "default")
	assert.NoError(t, r.Validate(), "validate")
	assert.Equal(t, "literal", r.Type)
}

func TestRule_Validate(t *testing.T) {
	r := inject.Rule{Type: "whatever"}
	assert.NoError(t, r.Default(), "default")
	assert.EqualError(t, r.Validate(), `invalid .type: "whatever" is invalid, must be one of:

  • literal
  • comment
  • style
  • script
  • inline style
  • inline script
  • google analytics
  • segment`)
}

func TestRules_Default(t *testing.T) {
	t.Run("type literal", func(t *testing.T) {
		rules := inject.Rules{
			"head": []*inject.Rule{
				{
					Value: `<script>var user = {}</script>`,
				},
			},
		}

		assert.NoError(t, rules.Default(), "default")
		assert.NoError(t, rules.Validate(), "validate")
	})
}

func TestRules_Validate(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		rules := inject.Rules{
			"head": []*inject.Rule{
				{
					Type: "inline script",
					// Value: "var user = {}",
				},
			},
		}

		assert.NoError(t, rules.Default(), "default")
		assert.EqualError(t, rules.Validate(), `head rule #1: .value is required`)
	})
}
