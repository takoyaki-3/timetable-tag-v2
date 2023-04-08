package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	// "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)

	allocator, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocator)
	defer chromedp.Cancel(ctx)

	fmt.Println("hello")

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://gtfs-data.jp/search"),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML(`document.documentElement.innerHTML`, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("hello")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		panic(err)
	}

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			row.Find("td").Each(func(k int, cell *goquery.Selection) {
				fmt.Printf("%v %v\n", cell.Text(), cell.Nodes)
			})
		})
	})
}
