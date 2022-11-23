package cache

import (
	"github.com/go-resty/resty/v2"
)

type Currency struct {
}

func GetCurrenciesAndSetLocalCache() {
	client := resty.New()
	client.SetRetryCount(4)
	rsp, _ := client.R().EnableTrace().Get("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json")
	if rsp != nil {

	}
}
