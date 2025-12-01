package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed static/*
var staticFiles embed.FS

// Note: No need to embed index.html anymore, it's compiled into the code!

type PageData struct {
	Match    bool
	ErrorMsg string
	Checked  bool
}

// Helper to bridge Templ components with Echo
func render(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	assetFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	e.StaticFS("/static", assetFS)

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/favicon.svg")
	})

	e.GET("/", func(c echo.Context) error {
		return render(c, Home())
	})

	e.POST("/check", checkRegexHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func checkRegexHandler(c echo.Context) error {
	pattern := c.FormValue("regex")
	testString := c.FormValue("testString")

	if pattern == "" && testString == "" {
		return render(c, Result(nil))
	}

	data := &PageData{Checked: true}

	re, err := regexp.Compile(pattern)
	if err != nil {
		data.ErrorMsg = err.Error()
	} else {
		data.Match = re.MatchString(testString)
	}

	return render(c, Result(data))
}
