package vibor

import (
	"fmt"
	"log"
	"strings"
	"tilescrap/pkg/csvmodel"
	"tilescrap/pkg/useragent"

	"github.com/gocolly/colly"
)

var vendor string
var url string
var userAgent string

var data []*csvmodel.Model

type VibortScraper struct{}

func NewScraper(v, URL string) *VibortScraper {

	vendor = v
	url = URL
	userAgent = useragent.GetUserAgent()
	return &VibortScraper{}
}

func (s *VibortScraper) Scrap(scrapingData chan<- []*csvmodel.Model) {

	c := colly.NewCollector()
	c.AllowURLRevisit = false
	//c.CheckHead = true
	c.UserAgent = userAgent

	p := c.Clone()

	c.OnRequest(onRequest)

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
		scrapingData <- nil
	})

	c.OnHTML("div.title_blok_rith_menu", func(d *colly.HTMLElement) {
		d.ForEach("a", func(i int, e *colly.HTMLElement) {
			link := e.Attr("href")
			link = e.Request.AbsoluteURL(link)
			p.Visit(link)
		})
	})

	/* c.OnHTML("body", func(b *colly.HTMLElement) {
		b.ForEach("div.title_blok_rith_menu a", func(i int, a *colly.HTMLElement) {
			link := a.Attr("href")
			link = a.Request.AbsoluteURL(link)
			p.Visit(link)
		})
	}) */

	p.OnRequest(onRequest)

	p.OnError(func(r *colly.Response, err error) {
		log.Printf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s\n", r.StatusCode, err)
	})

	p.OnHTML("body", func(b *colly.HTMLElement) {
		b.ForEach("div.product-line", func(i int, prod *colly.HTMLElement) {

			name := prod.ChildText(".products-item__name > span")
			inStock := prod.ChildText(".catalog__status-availability > span")
			collection := prod.ChildText(".product-line__collection > .product-line__collection-inner")
			price := prod.ChildText(".product-price .product-price__numb > .product-price__c")

			var char string
			prod.ForEach(".product-line-char", func(i int, pch *colly.HTMLElement) {
				key := pch.ChildText(".product-line-char__title")
				val := pch.ChildText(".product-property__value")
				key = strings.ReplaceAll(key, "\t", "")
				key = strings.ReplaceAll(key, "\n", "")
				val = strings.ReplaceAll(val, "\t", "")
				val = strings.ReplaceAll(val, "\n", "")

				if key != "" && val != "" {
					char += fmt.Sprintf("%s: %s\n", key, val)
				}
			})

			data = append(data, &csvmodel.Model{
				Vendor:          vendor,
				Product:         name,
				Characteristics: char,
				Price:           price,
				InStock:         inStock,
				Collection:      collection,
				URL:             b.Request.URL.String(),
			})
		})
	})

	c.OnScraped(func(r *colly.Response) {
		scrapingData <- data
	})

	c.Visit(url)

}

func onRequest(r *colly.Request) {
	log.Println("Переходим на", r.URL.String())
}
