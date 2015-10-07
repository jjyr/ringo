package ringo

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
)

type TemplateManage struct {
	templateRoot *template.Template
	templatePath string
}

func newTemplateManage() *TemplateManage {
	templateManage := TemplateManage{}
	templateManage.templateRoot = template.New("")
	return &templateManage
}

// SetTemplate setup the base template
func (tm *TemplateManage) SetTemplate(t *template.Template) {
	tm.templateRoot = t
}

// SetTemplatePath set template path
func (tm *TemplateManage) SetTemplatePath(templatePath string) {
	if tm.templatePath != "" {
		panic(fmt.Errorf("template path already set to %s, can not overwrite with %s", tm.templatePath, templatePath))
	}
	tm.templatePath = templatePath
}

// LoadTemplates load template from templatePath/name
func (tm *TemplateManage) LoadTemplates(names ...string) {
	for _, name := range names {
		content, err := ioutil.ReadFile(path.Join(tm.templatePath, name))
		if err != nil {
			log.Panicf("Load template error: %s", err)
		}
		template.Must(tm.templateRoot.New(name).Parse(string(content)))
	}
}

// FindTemplate find template, not found will cause panic
func (tm *TemplateManage) FindTemplate(name string) *template.Template {
	t := tm.templateRoot.Lookup(name)
	if t == nil {
		if tm.templatePath == "" {
			panic(fmt.Errorf("Can not found template %s", name))
		}
		// lazy load
		tm.LoadTemplates(name)
		return tm.FindTemplate(name)
	}
	return t
}
