package bank

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type PriceExtractor struct {
	pageReg string
	tag     string
}

func NewPriceExtractor(pageReg, tag string) PriceExtractor {
	return PriceExtractor{pageReg: pageReg, tag: tag}
}

func (ext PriceExtractor) ExtractPrice(body io.Reader) (float32, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return 0.00, fmt.Errorf("cannot parse the body to an HTML document: %w", err)
	}

	re := regexp.MustCompile(ext.pageReg)

	data := doTraverse(doc, ext.tag, re)

	if data == "" {
		return 0.00, fmt.Errorf("the document does not have a price value with the tag: %v", ext.tag)
	}

	return ext.getPrice(data), nil
}

func doTraverse(doc *html.Node, tag string, re *regexp.Regexp) string {
	var traverse func(n *html.Node, tag string, re *regexp.Regexp) string

	traverse = func(n *html.Node, tag string, re *regexp.Regexp) string {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isNodeWithPriceForBuying(c, tag, re) {
				return c.Data
			}

			data := traverse(c, tag, re)

			if data != "" {
				return data
			}
		}

		return ""
	}

	return traverse(doc, tag, re)
}

func isNodeWithPriceForBuying(n *html.Node, tag string, re *regexp.Regexp) bool {
	return n.Type == html.TextNode && n.Parent.Data == tag && re.MatchString(n.Data)
}

func (ext PriceExtractor) getPrice(data string) float32 {
	cleanedData := strings.ReplaceAll(data, "\u00a0", "")
	cleanedData = strings.ReplaceAll(cleanedData, ",", ".")
	price, _ := strconv.ParseFloat(cleanedData, 32)

	return float32(price)
}
