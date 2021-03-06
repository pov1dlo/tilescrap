package braer

import (
	"log"
	"tilescrap/pkg/csvmodel"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var vendor string
var url string

var data []*csvmodel.Model

type BraerScraper struct{}

func NewScraper(v, URL string) *BraerScraper {

	vendor = v
	url = URL
	return &BraerScraper{}
}

func (s *BraerScraper) Scrap(scrapingData chan<- []*csvmodel.Model) {

	c := colly.NewCollector()
	c.AllowURLRevisit = false

	extensions.RandomUserAgent(c)

	p := c.Clone()

	c.OnRequest(onRequest)

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
		scrapingData <- nil
	})

	c.OnHTML("body", func(b *colly.HTMLElement) {

		url = b.Request.AbsoluteURL(url)

		b.ForEach(".ul-goods-view-item a.ul-goods-view-link", func(i int, e *colly.HTMLElement) {

			link := e.Attr("href")

			link = e.Request.AbsoluteURL(link)
			if link != "" && link != url {
				p.Visit(link)
			}
		})
	})

	p.OnHTML("body", func(b *colly.HTMLElement) {

		b.ForEach("div.ul-goods-view-item", func(i int, e *colly.HTMLElement) {
			title := e.ChildText("div.ul-goods-view-details > div.ul-goods-view-title > div.js-goods-contenteditable")
			price := e.ChildAttr("div.ul-goods-view-price meta[itemprop=\"price\"]", "content")
			priceCurrency := e.ChildAttr("div.ul-goods-view-price meta[itemprop=\"priceCurrency\"]", "content")
			description := e.ChildText("div.ul-goods-view-descr > div.note > span")

			data = append(data, &csvmodel.Model{
				Vendor:        vendor,
				Product:       title,
				Price:         price,
				PriceCurrency: priceCurrency,
				Description:   description,
				URL:           b.Request.URL.String(),
			})
		})
	})

	p.OnRequest(func(r *colly.Request) {
		log.Println("Переходим на ", r.URL.String())
	})

	p.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
		scrapingData <- nil
	})

	c.OnScraped(func(r *colly.Response) {
		scrapingData <- data
	})

	c.Visit(url)
}

func onRequest(r *colly.Request) {
	log.Println("Переходим на", r.URL.String())
}
