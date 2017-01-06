package renderer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
)

func TestWrapBasic(t *testing.T) {
	e := echo.New()
	e.Renderer = Wrap(render.New(render.Options{
		Directory: "fixtures",
	}))

	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "gophers")
	})

	status, body := request(echo.GET, "/hello", e)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "<h1>Hello gophers</h1>\n", body)
}

func request(method, path string, e *echo.Echo) (int, string) {
	req, _ := http.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
