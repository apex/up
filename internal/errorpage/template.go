package errorpage

import "html/template"

// defaultPage is the default error page.
var defaultPage = template.Must(template.New("errorpage").Parse(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{.StatusText}} – {{.StatusCode}}</title>
    <style>
      html, body {
        margin: 0;
        width: 100%;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
        font-size: 15px;
      }

      .Error {
        font-size: 1.35em;
        {{with .Variables.color}}
          color: {{.}};
        {{else}}
          color: #6061BE;
        {{end}}
      }

      .Error .message {
        font-weight: 200;
        letter-spacing: 0.095em;
      }

      .Error .message a {
        text-decoration: none;
        color: inherit;
        {{with .Variables.color}}
          border-bottom: 1px dotted {{.}};
        {{else}}
          border-bottom: 1px dotted #6061BE;
        {{end}}
      }

      .Error .status {
        font-weight: 700;
      }

      .Error .code {
        display: none;
      }

      .Error .dot {
        font-weight: 100;
      }

      @media screen and (max-width: 800px) {
        body {
          font-size: 10px;
        }

        .Error {
          display: flex;
          flex-direction: column;
        }

        .Error .status {
          font-size: 1.1em;
        }

        .Error .dot {
          display: none
        }
      }
    </style>
  </head>
  <body>
    <div class="Error">
      <span class="status">{{.StatusText}}</span>
      <span class="code">{{.StatusCode}}</span>
      <span class="dot">•</span>
      {{with .Variables.support_email}}
        <span class="message">Please try your request again or <a href="mailto:{{.}}">contact support</a>.</span>
      {{else}}
        <span class="message">Please try your request again or contact support.</span>
      {{end}}
    </div>
  </body>
</html>`))
