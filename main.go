package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

type Crawl struct {
	Name    string
	Price   string
	Rating  string
	Reviews string
}

var CrawledItems []Crawl

func RunCrwler() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.Async(true),
	)

	c.OnHTML(" div[class^=\"puis-card-container\"]", func(h *colly.HTMLElement) {

		sr := strings.Split(h.ChildText("span.a-icon-alt"), "von")

		CrawledItems = append(CrawledItems, Crawl{
			Name:    h.ChildText("span.a-size-base-plus"),
			Price:   h.ChildText("span.a-offscreen"),
			Rating:  sr[0],
			Reviews: h.ChildText("span.s-underline-text"),
		})

	})

	c.OnRequest(func(h *colly.Request) {
		log.Printf("Connecting to %v", h.URL)
	})

	for i := 0; i < 25; i++ {

		url := fmt.Sprintf(`https://www.vbin.in/`, i)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Connot connect to site, Error: %v ", err)
		}
	}
	c.Wait()
}

func main() {
	RunCrwler()

	//json, err := json2.Marshal(CrawledItems)

	//if err != nil {
	//	log.Printf("Connot Marshal json, Error: %v ", err)
	//}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(CrawledItems)
	})

	log.Fatal(app.Listen(":3000"))
}
