package braer

import (
	"fmt"
	"log"
	uagent "tilescrap/useragent"

	"github.com/gocolly/colly"
)

var SiteURL string
var UserAgent string

func init() {

	SiteURL = "https://тротуарная-плитка-браер.рф/"
	UserAgent = uagent.GetUserAgent()

}

func Scrap() {

	go scrapOfficial1()
	go scrapOfficial2()
}

func scrapOfficial1() {

	//SiteURL := "https://braer.ru/catalog"

	/*
		links = soup.select(selector='.category > a')


		Нужно искользовать Selenium
	*/

}

func scrapOfficial2() {

	c := colly.NewCollector()
	c.AllowURLRevisit = false
	c.UserAgent = UserAgent

	productColly := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Переходим на", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s", r.StatusCode, err)
	})

	c.OnHTML("body", func(b *colly.HTMLElement) {

		SiteURL = b.Request.AbsoluteURL(SiteURL)

		b.ForEach(".ul-goods-view-item a.ul-goods-view-link", func(i int, e *colly.HTMLElement) {

			link := e.Attr("href")

			link = e.Request.AbsoluteURL(link)
			if link != "" && link != SiteURL {
				productColly.Visit(link)
			}
		})
	})

	productColly.OnHTML("body", func(b *colly.HTMLElement) {

		b.ForEach("div.ul-goods-view-item", func(i int, e *colly.HTMLElement) {
			title := e.ChildText("div.ul-goods-view-details > div.ul-goods-view-title > div.js-goods-contenteditable")
			price := e.ChildAttr("div.ul-goods-view-price meta[itemprop=\"price\"]", "content")
			priceCurrency := e.ChildAttr("div.ul-goods-view-price meta[itemprop=\"priceCurrency\"]", "content")
			description := e.ChildText("div.ul-goods-view-descr > div.note > span")

			fmt.Println(title, ":", price, ":", priceCurrency, ":", description)

		})
	})

	productColly.OnRequest(func(r *colly.Request) {
		log.Println("Переходим на ", r.URL)
	})

	c.Visit(SiteURL)
}
