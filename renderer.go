package renderer

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/unrolled/render"
)

// DefaultFuncs is a map of highly reusable template functions that can
// optionally be appended to a `render.Render` FuncMap
var DefaultFuncs = template.FuncMap{
	"dateFormat":   dateFormat,
	"htmlEscape":   htmlEscape,
	"htmlUnescape": htmlUnescape,
	"safeHTML":     safeHTML,
	"safeURL":      safeURL,
	"dict":         dictionary,
	"querify":      querify,
}

// RenderWrapper ...
type RenderWrapper struct {
	render *render.Render
}

// Render ...
func (r *RenderWrapper) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := r.render.HTML(w, 0, name, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return err
}

// Wrap ...
func Wrap(r *render.Render) *RenderWrapper {
	return &RenderWrapper{r}
}
