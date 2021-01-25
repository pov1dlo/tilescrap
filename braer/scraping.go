package braer

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

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

	SiteURL := "https://тротуарная-плитка-браер.рф/"
	UserAgent := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36"

	c := colly.NewCollector()
	c.AllowURLRevisit = false
	c.UserAgent = UserAgent

	productColly := c.Clone()

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
		log.Println("Visiting ", r.URL)
	})

	c.Visit(SiteURL)
}
