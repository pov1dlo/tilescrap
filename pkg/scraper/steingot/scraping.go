package steingot

import (
	"fmt"
	"log"
	"tilescrap/pkg/csvmodel"
	"tilescrap/pkg/useragent"

	"github.com/gocolly/colly"
)

var vendor string
var url string
var userAgent string

var data []*csvmodel.Model

type SteingotScraper struct{}

func NewScraper(v, URL string) *SteingotScraper {

	vendor = v
	url = URL
	userAgent = useragent.GetUserAgent()
	return &SteingotScraper{}
}

func (s *SteingotScraper) Scrap(scrapingData chan<- []*csvmodel.Model) {

	scrapingData <- nil
	return

	// парсинг еще не готов
	c := colly.NewCollector()
	c.AllowURLRevisit = false
	c.UserAgent = userAgent

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

	c.Visit(url)

}
