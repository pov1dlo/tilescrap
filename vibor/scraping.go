package vibor

import (
	"log"
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
	c.UserAgent = UserAgent

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visit", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalf("Код ответа: %d\nНе удалось подключиться к сайту по причине %s", r.StatusCode, err)
	})

	c.OnHTML("", func(e *colly.HTMLElement) {

	})

	c.Visit(SiteURL)

}

/* URL = 'https://купиплитку.рф'
GOODS = list()


def scrap():

    headers = {
        'User-agent': user_agent.get_user_agent()
    }

    r = requests.get(url=URL, headers=headers)

    soup = bs(r.text, 'lxml')

    links = soup.select(selector='.title_blok_rith_menu > a')

    for link in links:
        parse_link(link.get('href'), link.get_text())

    with open('{}.json'.format(COMPETITOR), 'w', encoding='utf-8') as outfile:
        json.dump(GOODS, outfile, ensure_ascii=False)


def parse_link(link, category=''):

    if SITE not in link:
        link = SITE + link

    r = requests.get(link)
    soup = bs(r.text, 'lxml')
    goods_list = soup.select(selector='.product-line ')

    for goods in goods_list:

        name = goods.select_one('.products-item__name > span')
        in_stock = goods.select_one('.catalog__status-availability > span')
        collection = goods.select_one('.product-line__collection > .product-line__collection-inner')
        price = goods.select_one('.product-price .product-price__numb > .product-price__c')
        old_price = None
        url_img = goods.select_one('.products-item__image > img')

        args = dict()
        args['date'] = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        args['competitor'] = COMPETITOR
        args['site'] = URL
        args['url'] = link
        args['category'] = category

        if name is not None:
            args['name'] = name.get_text()

        if in_stock is not None:
            args['in_stock'] = in_stock.get_text()

        if collection is not None:
            args['collection'] = collection.get_text()

        if price is not None:
            args['price'] = price.get_text()

        if url_img is not None:
            args['url_img'] = url_img.get('src')

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
