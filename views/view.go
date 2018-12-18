package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
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
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{ 
		Template: t, 
                Layout: layout,
		NotFound: notfound,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer

	w.Header().Set("Content-Type", "text/html")

	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data {
			Yield: data,
		}
	}

//	if v.NotFound {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}

	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
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
	v.Render(w, nil)
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

