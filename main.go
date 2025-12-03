package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/egreerdp/regex-checker/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed static/*
var staticFiles embed.FS

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
		return views.Render(c, views.Home())
	})

	e.POST("/check", checkRegexHandler)

	listenAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(listenAddr))
}

func checkRegexHandler(c echo.Context) error {
	pattern := c.FormValue("regex")
	testString := c.FormValue("testString")

	if pattern == "" && testString == "" {
		return views.Render(c, views.Result(nil))
	}

	data := &views.PageData{Checked: true}

	re, err := regexp.Compile(pattern)
	if err != nil {
		data.ErrorMsg = err.Error()
	} else {
		data.Match = re.MatchString(testString)
	}

	return views.Render(c, views.Result(data))
}
