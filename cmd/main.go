package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"tilescrap/pkg/csvmodel"
	"tilescrap/pkg/scraper"
	"tilescrap/pkg/scraper/braer"
	"tilescrap/pkg/scraper/steingot"
	"tilescrap/pkg/scraper/vibor"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var scrapers []scraper.Scraper
var scrapingData chan []*csvmodel.Model
var writer *csv.Writer
var useTBot bool

func main() {

	tokenTBot := flag.String("token", "", "Telegram token")
	chatID := flag.Int64("chat", 0, "Telegram channel")

	flag.Parse()

	if *tokenTBot != "" && *chatID != 0 {

		useTBot = true
	}

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
		"Коллекция",
		"Описание",
		"Ссылка"})

	writer.Flush()

	scrapers = append(scrapers,
		braer.NewScraper("Braer", "https://тротуарная-плитка-браер.рф/"),
		steingot.NewScraper("Steingot", "https://steingot.ru/katalog/"),
		vibor.NewScraper("Выбор", "https://купиплитку.рф"))

	scrapingData = make(chan []*csvmodel.Model, len(scrapers))

	for _, s := range scrapers {
		go s.Scrap(scrapingData)
	}

	for i := 0; i < len(scrapers); i++ {
		data := <-scrapingData

		if data == nil {
			continue
		}

		for _, s := range data {
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

	if useTBot == true {
		bot, err := tb.NewBot(
			tb.Settings{
				Token: *tokenTBot,
			})

		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 1)

		group := tb.ChatID(*chatID)

		currentDate := time.Now().Format("02.01.2006")
		document := tb.Document{

			File:     tb.FromDisk("result.csv"),
			FileName: fmt.Sprintf("Цены конкурентов от %s.csv", currentDate),
		}

		_, err = bot.Send(group, &document)
		if err != nil {
			log.Println(err)
		}

	}
}
