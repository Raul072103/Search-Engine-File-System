package cache

import (
	"MyFileExporer/backend/internal/repo/database"
	"MyFileExporer/common/models"
	"slices"
	"time"
)

// DefaultCacheSize is the default size for the slice containing files
const (
	DefaultCacheSize = 2048
	DefaultTTL       = time.Minute * 10
)

type Cache struct {
	cache *cache
}

type cache struct {
	// requestMap is the in-memory cache, keeps a
	requestsMap map[string]*cacheEntry
	// janitor handles the clearing of cache
	janitor *janitor
	// janitorMap is a map that gets updated each time a new request is cached
	janitorMap map[uint64]map[string]struct{}
}

type cacheEntry struct {
	TTL   uint64
	files []models.File
}

type janitor struct {
}

func newCache() *cache {
	requestMap := make(map[string]*cacheEntry, 2048)
	janitorMap := make(map[uint64]map[string]struct{}, 24)

	return &cache{
		requestsMap: requestMap,
		janitor:     nil,
		janitorMap:  janitorMap,
	}
}

func (cache *cache) Find(fs *database.FileSearchRequest) []models.File {
	key := buildKeyFromRequest(fs)

	entry, cacheHit := cache.requestsMap[key]
	if cacheHit {
		return entry.files
	}

	return nil
}

func (cache *cache) Add(fs *database.FileSearchRequest, files []models.File) error {
	return nil
}

func buildKeyFromRequest(fs *database.FileSearchRequest) string {
	result := ""

	if fs.Words != nil {
		slices.Sort(*fs.Words)

		for _, word := range *fs.Words {
			result += word
		}
	}

	if fs.Name != nil {
		result += *fs.Name
	}

	if fs.Extension != nil {
		slices.Sort(*fs.Extension)

		for _, extension := range *fs.Extension {
			result += extension
		}
	}

	return result
}
