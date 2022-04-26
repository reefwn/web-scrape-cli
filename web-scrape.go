package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func main () {
	imgCmd := flag.NewFlagSet("image", flag.ExitOnError)
	imgUrl := imgCmd.String("u", "url", "url to scrape")

	switch os.Args[1] {
		case "image": 
			handleImageScrape(imgCmd, imgUrl)
	}
}

func handleImageScrape(imgCmd *flag.FlagSet, imgUrl *string) {
	imgCmd.Parse(os.Args[2:])

	start := time.Now()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var res string
	var html string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(*imgUrl),
		// TODO: scroll page
		// chromedp.ScrollIntoView("footer"),
		// chromedp.WaitVisible("footer > div"),
		chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			dom, er := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			html = dom
			return er
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Selection.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		fmt.Println("src", i, src)
	})

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}