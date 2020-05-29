package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"log"
	"strings"
	"time"
)

//func main() {
//	e := echo.New()
//
//	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
//		AllowOrigins: []string{"*"},
//	}))
//	e.HidePort = true
//	e.HideBanner = true
//	e.Use(middleware.Recover())
//	e.Use(middleware.Logger())
//
//	e.GET("/trend/weekly", func(c echo.Context) error {
//		trend, err := qiitaTrend("/?scope=weekly")
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, err)
//		}
//		return c.String(http.StatusOK, trend)
//	})
//	e.GET("/trend/daily", func(c echo.Context) error {
//		trend, err := qiitaTrend("")
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, err)
//		}
//		return c.String(http.StatusOK, trend)
//	})
//	e.GET("/trend/monthly", func(c echo.Context) error {
//		trend, err := qiitaTrend("/?scope=monthly")
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, err)
//		}
//		return c.String(http.StatusOK, trend)
//	})
//	e.Logger.Fatal(e.Start(":8081"))
//}

func qiitaTrend(endPoint string) (string, error) {
	//url := "https://qiita.com" + endPoint
	url := "https://qiita.com"
	//url := "https://qiita.com/?scope=monthly"

	driver := agouti.ChromeDriver()
	err := driver.Start()
	if err != nil {
		log.Println(err)
		fmt.Print("error1")
		return "", err
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Println(err)
		fmt.Print("error2")
		return "", err
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

	//doc, err := goquery.NewDocument(url)
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

func main() {
	driver := agouti.ChromeDriver()

	// headlessにする場合はこっちを使う
	// driver := agouti.ChromeDriver(
	//  agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
	// )
	if err := driver.Start(); err != nil {
		log.Fatalf("driverの起動に失敗しました : %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("セッション作成に失敗しました : %v", err)
	}

	// ブラウザはChromeを指定して起動
	//driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			//"--window-size=300,1200",
			//"--blink-settings=imagesEnabled=false", // don't load images
			"--disable-gpu",                        // ref: https://developers.google.com/web/updates/2017/04/headless-chrome#cli
			"no-sandbox",                           // ref: https://github.com/theintern/intern/issues/878
			//"disable-dev-shm-usage",                // ref: https://qiita.com/yoshi10321/items/8b7e6ed2c2c15c3344c6
		}),
	)
	if err := driver.Start(); err != nil {
		fmt.Print("error1")
		fmt.Print(err)
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		fmt.Print("error2")
		fmt.Print(err)
		log.Fatalf("Failed to open page:%v", err)
	}
	// ログインページに遷移
	if err := page.Navigate("https://qiita.com/login"); err != nil {
	//if err := page.Navigate(""); err != nil {
		fmt.Print("error3")
		//fmt.Print(err)
		log.Fatalf("%v", err)
	}
	// ID, Passの要素を取得し、値を設定
	identity := page.FindByID("identity")
	password := page.FindByID("password")
	identity.Fill("YoshikiTsuikada")
	password.Fill("monster0323")
	// formをサブミット
	if err := page.FindByClass("loginSessionsForm_submit").Submit(); err != nil {
		fmt.Print("error4")
		fmt.Print(err)
		//log.Fatalf("Failed to login:%v", err)
	}
	// 処理完了後、3秒間ブラウザを表示しておく
	time.Sleep(3 * time.Second)
}
