package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.HidePort = true
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/trend", func(c echo.Context) error {
		trend, err := qiitaTrend()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, trend)
	})
	e.Logger.Fatal(e.Start(":8081"))
}

func qiitaTrend() (string, error) {
	url := "https://qiita.com/"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	selection := doc.Find("div[data-hyperapp-app='Trend']")
	trend, ok := selection.Attr("data-hyperapp-props")
	if !ok {
		return "", errors.New("Internal Server Error")
	}
	return trend, nil
}
