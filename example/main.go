package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/syntaqx/renderer"
	"github.com/unrolled/render"
)

type Person struct {
	Name string
}

func main() {
	e := echo.New()

	// Enable debug logging
	e.Debug = true

	// Keeps the DefaultFuncs provided by renderer
	funcs := []template.FuncMap{renderer.DefaultFuncs}

	// Create an instance of unrolled/render with app-specific configurations.
	r := render.New(render.Options{
		Layout:        "layout",
		Directory:     "templates",
		Extensions:    []string{".html"},
		IsDevelopment: e.Debug,
		Funcs:         funcs,
	})

	// Wrap the render instance with the Renderer compliant interface.
	e.Renderer = renderer.Wrap(r)

	// Attach middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Let's give it a go!
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", Person{
			Name: "syntaqx",
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
