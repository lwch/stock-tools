package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lwch/runtime"
)

var urls = []string{
	"http://quote.cfi.cn/stockList.aspx?t=11", // 上交所
	"http://quote.cfi.cn/stockList.aspx?t=13", // 深交所
}

func main() {
	for _, u := range urls {
		run(u)
	}
}

func run(u string) {
	resp, err := http.Get(u)
	runtime.Assert(err)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	runtime.Assert(err)
	doc.Find("#divcontent table td").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		text = text[strings.Index(text, "(")+1 : len(text)-1]
		fmt.Println(text)
	})
}
