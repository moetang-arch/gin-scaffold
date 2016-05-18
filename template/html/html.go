package html

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type htmlTemplateRender struct {
	RenderMap       map[string]map[string]*template.Template
	DefaultTemplate *template.Template
}

func NewHtmlTemplateRender(defaultTemplate *template.Template) *htmlTemplateRender {
	r := htmlTemplateRender{
		RenderMap:       make(map[string]map[string]*template.Template),
		DefaultTemplate: defaultTemplate,
	}
	return r
}

func (this *htmlTemplateRender) AddTemplate(templateName, typeName string, template *template.Template) {
	m, ok := this.RenderMap[templateName]
	if ok {
		m[typeName] = templateName
	} else {
		m = make(map[string]*template.Template)
		m[typeName] = templateName
	}
}

func (this *htmlTemplateRender) CustomHTML(c *gin.Context, code int, templateName, typeName, name string, obj interface{}) {
	m, ok := this.RenderMap[templateName]
	c.Status(code)
	var t *template.Template = this.DefaultTemplate
	if ok {
		t1, ok := m[typeName]
		if ok {
			t = t1
		}
	}
	renderByTemplate(c, t, name, obj)
}

func (this *htmlTemplateRender) HTML(c *gin.Context, code int, name string, obj interface{}) {
	c.Status(code)
	renderByTemplate(c, this.DefaultTemplate, name, obj)
}

func renderByTemplate(c *gin.Context, t *template.Template, name string, obj interface{}) {
	if err := render(c.Writer, t, name, obj); err != nil {
		panic(err)
	}
}

var htmlContentType = []string{"text/html; charset=utf-8"}

func render(w http.ResponseWriter, t *template.Template, name string, obj interface{}) error {
	writeContentType(w, htmlContentType)
	if len(name) == 0 {
		return t.Execute(w, obj)
	}
	return t.ExecuteTemplate(w, name, obj)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
