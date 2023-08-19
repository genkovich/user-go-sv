package token

import (
	"fmt"
	"strconv"
	"user-service/pkg/cache"
)

type GeneratorCached struct {
	Generator Generator
	Cache     cache.Cache
}

func NewGeneratorCached(generator Generator, cache cache.Cache) *GeneratorCached {
	return &GeneratorCached{Generator: generator, Cache: cache}
}

func (g *GeneratorCached) Generate(login string, role string) (*Token, error) {
	cacheKey := fmt.Sprintf("token:%s:%s", login, role)
	cachedToken, err := g.Cache.Get(cacheKey)
	if err == nil && cachedToken != "" {
		return NewToken(cachedToken), nil
	}

	token, err := g.Generator.Generate(login, role)
	if err != nil {
		return nil, err
	}

	durationInSeconds := int((ExpirationDuration).Seconds())
	err = g.Cache.Set(cacheKey, token.GetJwtToken(), strconv.Itoa(durationInSeconds))

	if err != nil {
		return nil, err
	}

	return token, nil
}
