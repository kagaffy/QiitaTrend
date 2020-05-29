package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sclevine/agouti"
	"log"
	"net/http"
	"strings"
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

	//driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	//loginQiita(driver)

	e.GET("/trend/weekly", func(c echo.Context) error {
		trend, err := qiitaTrend("/?scope=weekly")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, trend)
	})
	e.GET("/trend/daily", func(c echo.Context) error {
		trend, err := qiitaTrend("")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, trend)
	})
	e.GET("/trend/monthly", func(c echo.Context) error {
		trend, err := qiitaTrend("/?scope=monthly")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, trend)
	})
	e.Logger.Fatal(e.Start(":8081"))
}

func qiitaTrend(endPoint string) (string, error) {
	url := "https://qiita.com" + endPoint

	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			//"--window-size=300,1200",
			"--blink-settings=imagesEnabled=false", // don't load images
			//"--disable-gpu",                        // ref: https://developers.google.com/web/updates/2017/04/headless-chrome#cli
			//"no-sandbox",                           // ref: https://github.com/theintern/intern/issues/878
			//"disable-dev-shm-usage",                // ref: https://qiita.com/yoshi10321/items/8b7e6ed2c2c15c3344c6
		}),
	)

	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}
	// ログインページに遷移
	if err := page.Navigate("https://qiita.com/login"); err != nil {
		log.Fatalf("Failed to navigate:%v", err)
	}
	// ID, Passの要素を取得し、値を設定
	identity := page.FindByID("identity")
	password := page.FindByID("password")
	identity.Fill("YoshikiTsukada")
	password.Fill("monster0323")
	// formをサブミット
	if err := page.FindByClass("loginSessionsForm_submit").Submit(); err != nil {
		log.Fatalf("Failed to login:%v", err)
	}

	err = page.Navigate(url)
	if err != nil {
		log.Println(err)
		fmt.Print("error3")
		return "", err
	}

	content, err := page.HTML()
	if err != nil {
		log.Println(err)
		fmt.Print("error4")
		return "", err
	}

	reader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)

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
