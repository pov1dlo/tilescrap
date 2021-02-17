package scraper

import "tilescrap/pkg/csvmodel"

//Scraper ...
type Scraper interface {
	Scrap(chan<- []*csvmodel.Model)
}
