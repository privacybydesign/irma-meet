package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/privacybydesign/irma-meet/internal/i18n"
)

//go:embed templates
var templateFS embed.FS

const (
	docsURL = "https://privacybydesign.foundation/uitgifte/"

	pageHome         = "home"
	pageAuthenticate = "authenticate"
	pageCreateFree   = "create_free"
	pageCreateExam   = "create_exam"
	pageCreated      = "created"
	pageJoin         = "join"
	pageError        = "error"
)

var pageNames = []string{pageHome, pageAuthenticate, pageCreateFree, pageCreateExam, pageCreated, pageJoin, pageError}

// view is the root value every template renders. Translations are methods on
// it so templates can write {{.T "key"}} without passing a locale around.
type view struct {
	tr   *i18n.Translator
	Lang string
	Data any
}

// T returns the translation of key; html/template escapes the result.
func (v view) T(key string) string {
	return v.tr.T(v.Lang, key)
}

// HTML returns a translation that is itself markup. It is only ever called
// from templates with a literal key, so the text is part of the source tree.
func (v view) HTML(key string) template.HTML {
	return template.HTML(v.tr.T(v.Lang, key)) //nolint:gosec // static translation string
}

// HTMLf is HTML for translations with %s placeholders. The arguments are
// escaped, so user input can be put into a trusted sentence.
func (v view) HTMLf(key string, args ...string) template.HTML {
	escaped := make([]any, len(args))
	for i, arg := range args {
		escaped[i] = template.HTMLEscapeString(arg)
	}

	return template.HTML(fmt.Sprintf(v.tr.T(v.Lang, key), escaped...)) //nolint:gosec // arguments are escaped above
}

// YiviMessages is the feedback shown when a Yivi session fails, keyed by
// the status custom.js reacts to.
func (v view) YiviMessages() map[string]string {
	return map[string]string{
		"cancelled": v.T("The meeting could not be started because the Yivi session was cancelled. This usually means your Yivi app does not yet contain the credentials required for this meeting, or you declined the disclosure request. Please make sure the requested data is loaded in your Yivi app and try again."),
		"timeout":   v.T("The meeting could not be started because the Yivi session timed out. Please try again and approve the disclosure request in your Yivi app."),
		"error":     v.T("The meeting could not be started because the Yivi server could not be reached. Please try again in a moment."),
		"fallback":  v.T("The meeting could not be started because the Yivi session did not complete. Please make sure your Yivi app contains the required credentials and try again."),
		"docsLabel": v.T("Read which credentials you need and how to load them"),
		"docsUrl":   docsURL,
	}
}

type renderer struct {
	tr    *i18n.Translator
	pages map[string]*template.Template
}

func newRenderer(tr *i18n.Translator) (*renderer, error) {
	sub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		return nil, err
	}

	r := &renderer{tr: tr, pages: map[string]*template.Template{}}

	for _, name := range pageNames {
		tmpl, err := template.ParseFS(sub, "layout.html", "partials.html", name+".html")
		if err != nil {
			return nil, fmt.Errorf("page %s: %w", name, err)
		}

		r.pages[name] = tmpl
	}

	return r, nil
}

// render writes the page with the given status. The body is built first so a
// template error does not leave a half-written response.
func (r *renderer) render(w http.ResponseWriter, page, lang string, status int, data any) error {
	var body bytes.Buffer

	v := view{tr: r.tr, Lang: lang, Data: data}
	if err := r.pages[page].ExecuteTemplate(&body, "layout.html", v); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, err := w.Write(body.Bytes())

	return err
}
