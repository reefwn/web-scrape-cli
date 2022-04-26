package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"reefwn/web-scrape-cli/utils"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)


func HandleImageCmd(imgCmd *flag.FlagSet, imgUrl, imgAttr, imgFolder *string) {
	imgCmd.Parse(os.Args[2:])

	if *imgUrl == "" {
		log.Fatal("url is required")
	}

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(*imgUrl),
		// TODO: scroll page
		// chromedp.ScrollIntoView("footer"),
		// chromedp.WaitVisible("footer > div"),
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
		src, _ := s.Attr(*imgAttr)

		fmt.Println("downloading ... ", src)

		err := utils.DownloadFile(src, *imgFolder, strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}
	})
}