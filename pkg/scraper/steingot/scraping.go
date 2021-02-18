package steingot

import (
	"fmt"
	"log"
	"strings"
	"tilescrap/pkg/csvmodel"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var vendor string
var url string

var data []*csvmodel.Model

type SteingotScraper struct{}

func NewScraper(v, URL string) *SteingotScraper {

	vendor = v
	url = URL
	return &SteingotScraper{}
}

func (s *SteingotScraper) Scrap(scrapingData chan<- []*csvmodel.Model) {

	c := colly.NewCollector()
	c.AllowURLRevisit = false

	extensions.RandomUserAgent(c)

	p := c.Clone()
	p.AllowURLRevisit = true

	c.OnRequest(onRequest)

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
		scrapingData <- nil
	})

	c.OnHTML("div.product-catalog__categories", func(e *colly.HTMLElement) {

		e.ForEach(".product-catalog__categories-item ", func(i int, a *colly.HTMLElement) {
			link := a.Attr("href")
			link = a.Request.AbsoluteURL(link)

			p.Visit(link)
		})

	})

	p.OnRequest(onRequest)

	p.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
	})

	p.OnHTML("body", func(b *colly.HTMLElement) {
		b.ForEach("div.product-catalog__product-item", func(i int, prod *colly.HTMLElement) {
			name := prod.ChildText("div.product-catalog__product-item-name")
			price := prod.ChildText("div.product-catalog__product-item-cost:nth-child(1)")
			price = strings.TrimPrefix(price, "Цена: ")
			priceCurrency := prod.ChildText("div.product-catalog__product-item-cost:nth-child(2)")

			data = append(data, &csvmodel.Model{
				Vendor:        vendor,
				Product:       name,
				Price:         price,
				PriceCurrency: priceCurrency,
			})
		})

		nextpage := b.ChildAttr(".product-catalog__page-number-next", "href")
		fmt.Println(nextpage, "\n")
		if nextpage != "" {
			link := b.Request.AbsoluteURL(nextpage)
			p.Visit(link)

		}

	})

	c.OnScraped(func(r *colly.Response) {
		scrapingData <- data
	})

	c.Visit(url)

}

func onRequest(r *colly.Request) {
	log.Println("Переходим на", r.URL.String())
}
