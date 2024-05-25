package bank

import (
	pageEnt "PriceWatcher/internal/entities/page"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

type BankRequester struct{}

func (r BankRequester) RequestPage(url string) (pageEnt.Response, error) {
	browser := rod.New().Timeout(time.Minute).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.sberbank.ru/ru/quotes/metalbeznal")
	time.Sleep(20 * time.Second)
	page.HTML()

	html, err := page.HTML()
	if err != nil {
		return pageEnt.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}

	respReader := strings.NewReader(html)

	return pageEnt.Response{Body: respReader}, nil
}
