package vibor

import (
	"fmt"
	"log"
	"strings"

	uagent "tilescrap/useragent"

	"github.com/gocolly/colly"
)

var SiteURL string
var UserAgent string

func init() {

	SiteURL = "https://купиплитку.рф"
	UserAgent = uagent.GetUserAgent()
}

func Scrap() {

	c := colly.NewCollector()
	c.CheckHead = true
	c.UserAgent = UserAgent

	p := c.Clone()

	c.OnRequest(onRequest)

	c.OnError(onError)

	c.OnHTML("div.title_blok_rith_menu", func(d *colly.HTMLElement) {
		d.ForEach("a", func(i int, e *colly.HTMLElement) {
			link := e.Attr("href")
			link = e.Request.AbsoluteURL(link)
			p.Visit(link)
		})
	})

	p.OnRequest(onRequest)

	p.OnError(onError)

	p.OnHTML("body", func(b *colly.HTMLElement) {
		b.ForEach("div.product-line", func(i int, prod *colly.HTMLElement) {

			name := prod.ChildText(".products-item__name > span")
			inStock := prod.ChildText(".catalog__status-availability > span")
			collection := prod.ChildText(".product-line__collection > .product-line__collection-inner")
			price := prod.ChildText(".product-price .product-price__numb > .product-price__c")
			//old_price := 0
			//urlImg := prod.ChildAttr(".products-item__image > img", "src")

			//fmt.Println(name, inStock, collection, price, urlImg)

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

			fmt.Println(char)

		})
	})

	c.Visit(SiteURL)

}

func onRequest(r *colly.Request) {
	log.Println("Переходим на", r.URL)
}

func onError(r *colly.Response, err error) {
	log.Fatalf("Не удалось подключиться к сайту по причине %s\nКод ответа: %d\n", err, r.StatusCode)
}

/*
   char = ''

   product_char = goods.select('.product-line-char')
   for pch in product_char:
       product_char_key = pch.select_one('.product-line-char__title')
       product_char_value = pch.select_one('.product-property__value')

       if product_char_key is not None and product_char_value is not None:
           char += "{}: {}".format(product_char_key.get_text(), product_char_value.get_text()) \
               .replace('\t', '') \
               .replace('\n', '') \
               .strip()

           char += '\n'

   product_char_color = goods.select_one('.product-line-char-color')
   if product_char_color is not None:
       product_char_key = product_char_color.select_one('.product-line-char__title')
       product_char_value = product_char_color.select_one('.product-line__collection-inner > img')

       if product_char_key is not None and product_char_value is not None:
           char += "{}: {}".format(product_char_key.get_text(), product_char_value.get_text())\
               .replace('\t', '')\
               .replace('\n', '')\
               .strip()

           char += '\n'

   args['char'] = char

   GOODS.append(args)
*/
