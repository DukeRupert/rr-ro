// pkg/template/template.go
package template

import (
    "html/template"
    "io"
    "github.com/labstack/echo/v4"
)

type Template struct {
    Templates *template.Template
}

func NewTemplate() *Template {
    return &Template{
        Templates: template.Must(template.ParseFiles(
            "templates/layout.html",
            "templates/products.html",
        )),
    }
}

// Render implements echo.Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.Templates.ExecuteTemplate(w, name, data)
}