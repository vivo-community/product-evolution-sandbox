package widgets_import

import (
	"errors"
	"fmt"
	"github.com/oxtoacart/bpool"
	"log"
	"path/filepath"
	"text/template"
)

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

func Preview(tmpl string) {
	log.Println(tmpl)
}

// NOTE: stole code from here:
//https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6
func LoadTemplates(conf Config) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(conf.Templates.Layout + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob(conf.Templates.Include + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mainTemplate := template.New("main")
	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	log.Println("elastic mapping templates loading successful")

	bufpool = bpool.NewBufferPool(64)
	//log.Println("buffer allocation successful")
}

func RenderTemplate(name string) (string, error) {
	tmpl, ok := templates[name]

	if !ok {
		msg := fmt.Sprintf("could not find template %s\n", name)
		return "", errors.New(msg) 
	}
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// TODO: bogus 'data' = "hello" - should do something else
	err := tmpl.Execute(buf, "hello")
	if err != nil {
		msg := fmt.Sprintf("error executing template %s\n", err)
		return "", errors.New(msg)
	}
	return buf.String(), nil
}
