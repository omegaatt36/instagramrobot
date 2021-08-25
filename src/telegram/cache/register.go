package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

var Cache *goCache.Cache

func Register() {
	// TODO: load from env
	Cache = goCache.New(3*time.Hour, 30*time.Minute)
}

func Set(key string, value Media) {
	Cache.Set(key, value, goCache.DefaultExpiration)
}

func Get(key string) (Media, bool) {
	if media, found := Cache.Get(key); found {
		return media.(Media), true
	}
	return Media{}, false
}
