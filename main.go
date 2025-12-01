package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed index.html
var indexHTML embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type PageData struct {
	Match    bool
	ErrorMsg string
	Checked  bool
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseFS(indexHTML, "index.html")),
	}
	e.Renderer = t

	assetFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	e.StaticFS("/static", assetFS)

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.POST("/check", checkRegexHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func checkRegexHandler(c echo.Context) error {
	pattern := c.FormValue("regex")
	testString := c.FormValue("testString")

	data := PageData{Checked: true}

	if pattern == "" && testString == "" {
		return c.NoContent(http.StatusOK)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		data.ErrorMsg = err.Error()
	} else {
		data.Match = re.MatchString(testString)
	}

	return c.Render(http.StatusOK, "result", data)
}
