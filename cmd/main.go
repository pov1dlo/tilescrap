package main

import (
	"encoding/csv"
	"os"
	"tilescrap/pkg/csvmodel"
	"tilescrap/pkg/scraper"
	"tilescrap/pkg/scraper/braer"
)

var scrapers []scraper.Scraper
var scrapingData chan []*csvmodel.Model
var writer *csv.Writer

func main() {

	file, _ := os.Create("result.csv")
	defer file.Close()

	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	file.Write(bomUtf8)

	writer = csv.NewWriter(file)
	writer.Comma = '\t'

	writer.Write([]string{"Продавец",
		"Товар",
		"Характеристики",
		"Цена",
		"Валюта",
		"В наличии",
		"Колекция",
		"Описание",
		"Ссылка"})

	writer.Flush()

	scrapers = append(scrapers,
		braer.NewScraper("Braer", "https://тротуарная-плитка-браер.рф/"))

	scrapingData = make(chan []*csvmodel.Model, len(scrapers))

	for _, s := range scrapers {
		go s.Scrap(scrapingData)
	}

	defer close(scrapingData)

	for sd := range scrapingData {
		for _, s := range sd {
			writer.Write([]string{
				s.Vendor,
				s.Product,
				s.Characteristics,
				s.Price,
				s.PriceCurrency,
				s.InStock,
				s.Collection,
				s.Description,
				s.URL})
			writer.Flush()
		}
	}

	file.Close()
}
