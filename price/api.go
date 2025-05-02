package price

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	gcache "github.com/patrickmn/go-cache"
	"github.com/usecorn/common-lib/conversions"
	"github.com/usecorn/common-lib/eth"
)

type PriceAPI interface {
	GetAllPrices(ctx context.Context, network string) (map[string]*big.Float, error)
	GetPrice(ctx context.Context, network, token string) (*big.Float, error)
}

type priceAPI struct {
	cacheExpiration time.Duration
	client          *http.Client
	cache           *gcache.Cache
}

func NewPriceAPI(apiTimeout, cacheExpiration time.Duration) PriceAPI {
	return &priceAPI{
		cacheExpiration: cacheExpiration,
		client:          &http.Client{Timeout: apiTimeout},
		cache:           gcache.New(cacheExpiration, 2*cacheExpiration),
	}
}

func (papi *priceAPI) GetAllPrices(ctx context.Context, network string) (map[string]*big.Float, error) {
	// Check if the prices are cached
	if cachedPrices, found := papi.cache.Get(network); found {
		if prices, ok := cachedPrices.(map[string]*big.Float); ok {
			return prices, nil
		}
	}
	var url string
	switch network {
	case eth.EthereumNetwork:
		url = "https://api.usecorn.com/api/v1/price/all"
	case eth.CornMainnet:
		url = "https://api.usecorn.com/api/v1/price/corn/all"
	default:
		return nil, errors.Newf("unsupported network %s", network)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := papi.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	prices := make(map[string]*big.Float, len(result))
	for k, v := range result {
		price, ok := conversions.NewLargeFloat().SetString(v)
		if !ok {
			return nil, err
		}
		prices[k] = price
	}

	// Cache the prices
	papi.cache.Set(network, prices, papi.cacheExpiration)

	return prices, nil
}

func (papi *priceAPI) GetPrice(ctx context.Context, network, token string) (*big.Float, error) {
	prices, err := papi.GetAllPrices(ctx, network)
	if err != nil {
		return nil, err
	}
	price, ok := prices[token]
	if !ok {
		return nil, errors.Newf("price not found for token: '%s'", token)
	}
	return price, nil
}
