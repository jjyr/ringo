package ringo

import "testing"

func TestTemplateManage(t *testing.T) {
	templateManage := newTemplateManage()
	templateManage.SetTemplatePath(".")
	for i, name := range []string{"app.go", "context.go", "binding/form.go", "render/render.go"} {
		templateName := templateManage.FindTemplate(name).Name()
		if name != templateName {
			t.Errorf("test case %d failed, name: %s not match template name: %s", i+1, name, templateName)
		}
	}

	func() {
		defer func() { recover() }()
		templateManage.FindTemplate("not_exists_template")
		t.Errorf("template not found should cause panic")
	}()
}
