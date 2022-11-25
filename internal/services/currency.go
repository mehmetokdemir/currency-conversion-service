package services

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/patrickmn/go-cache"
	"strings"
)

type CurrencyService struct {
	Cache *cache.Cache
}

func NewCurrencyService(goCache *cache.Cache) CurrencyService {
	return CurrencyService{Cache: goCache}
}

func (s *CurrencyService) getCurrenciesFromExternalService() entity.Currencies {
	var currencyMap = make(map[string]string)
	client := resty.New()
	client.SetRetryCount(4)
	rsp, err := client.R().EnableTrace().Get("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json")
	if err != nil {
		return nil
	}

	if err = json.Unmarshal(rsp.Body(), &currencyMap); err != nil {
		return nil
	}

	return currencyMap
}

func (s *CurrencyService) SetLocalCacheToCurrencies() {
	currencyMap := s.getCurrenciesFromExternalService()
	for k, v := range currencyMap {
		s.Cache.Set(strings.ToUpper(k), v, cache.NoExpiration)
	}
}

func (s *CurrencyService) CheckCurrencyCodeIsExist(code string) bool {
	_, ok := s.Cache.Get(code)
	return ok
}
