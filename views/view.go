package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"lenslockedbr.com/context"

	"github.com/gorilla/csrf"
)

var (
	LayoutDir string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
	NotFound bool
}

func NewView(layout string, notfound bool, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap {
		"csrfField": func() (template.HTML, error) {
			return " ", errors.New("csrfField is not implemented")
		},
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{ 
		Template: t, 
                Layout: layout,
		NotFound: notfound,
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, 
                      data interface{}) {
	var buf bytes.Buffer
	var vd Data

	w.Header().Set("Content-Type", "text/html")

	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data {
			Yield: data,
		}	
	}

	// Lookup the alert and assign it if one is persisted
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}

	vd.User = context.User(r.Context())

	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap {
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem "+
                              "persists, please email " + 
                              "support@lenslockedbr.com",
                           http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

/////////////////////////////////////////////////////////////////////
//
// HELPER FUNCTIONS
//
/////////////////////////////////////////////////////////////////////

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

//
// addTemplate takes in a slice of strings representing file paths
// for templates, and it prepends the TemplateDir directory to each
// string in the slice
//
// Eg the input {"home"} would result in the output {"views/home"}
// if TemplateDir == "views/"
//
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

//
// addTemplateExt takes in a slice of strings representing file paths
// for templates, and it appends the TemplateExt extension to each
// string in the slice
//
// Eg the input {"home"} would result in the output {"home.gohtml"}
// if TemplateExt == ".gohtml"
//
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

