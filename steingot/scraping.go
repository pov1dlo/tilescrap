package steingot

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func Scrap() {

	SiteURL := "https://steingot.ru/katalog/"
	UserAgent := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36"

	c := colly.NewCollector()
	c.AllowURLRevisit = false
	c.UserAgent = UserAgent

	c.OnHTML("dic.product-catalog__categories", func(e *colly.HTMLElement) {

		//e.ForEach("div.product-catalog__categories-item > a.product-catalog__categories-item", func(i int, a *colly.HTMLElement) {
		e.ForEach(".product-catalog__categories-item ", func(i int, a *colly.HTMLElement) {
			link := a.Attr("href")
			link = a.Request.AbsoluteURL(link)
			category := a.Text

			fmt.Println(link, ":", category)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(SiteURL)
}
