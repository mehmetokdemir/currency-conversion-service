package currency

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/patrickmn/go-cache"
	"strings"
)

type Service struct {
	Cache *cache.Cache
}

func NewCurrencyService(goCache *cache.Cache) Service {
	return Service{Cache: goCache}
}

func (s *Service) getCurrenciesFromExternalService() dto.Currencies {
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

func (s *Service) SetLocalCacheToCurrencies() {
	currencyMap := s.getCurrenciesFromExternalService()
	for k, v := range currencyMap {
		s.Cache.Set(strings.ToUpper(k), v, cache.NoExpiration)
	}
}

func (s *Service) CheckIsCurrencyCodeExist(code string) bool {
	_, ok := s.Cache.Get(code)
	return ok
}
